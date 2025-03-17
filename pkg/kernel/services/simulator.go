package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/utils"
	"github.com/samber/lo"
	"github.com/schollz/progressbar/v3"
)

type SimulationStatus string

const (
	StatusPending   SimulationStatus = "pending"
	StatusRunning   SimulationStatus = "running"
	StatusCompleted SimulationStatus = "completed"
	StatusFailed    SimulationStatus = "failed"
)

type SimulationJob struct {
	ID          string           `json:"id"`
	Status      SimulationStatus `json:"status"`
	Result      *SimulationView  `json:"result,omitempty"`
	Error       string           `json:"error,omitempty"`
	CallbackURL string           `json:"callback_url"`
	CreatedAt   time.Time        `json:"created_at"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
}

type SimulatorConfig struct {
	GameName       string
	ReportPath     string
	Spins          int64
	Wager          int64
	Workers        int
	GenerateParams interface{}
	CallbackURL    string
}

type KeepGenerateWrapper func(engine.Context, engine.Spin, engine.SpinFactory) (engine.Spin, error)

type SimulatorService struct {
	boot             *engine.Bootstrap
	generateParams   interface{}
	keepGenerate     bool
	keepGenerateFunc KeepGenerateWrapper
	jobs             map[string]*SimulationJob
	jobsMux          sync.RWMutex
	client           *http.Client
	gameWrappers     map[string]KeepGenerateWrapper
}

func NewSimulatorService() *SimulatorService {
	return &SimulatorService{
		boot:         engine.GetFromContainer(),
		jobs:         make(map[string]*SimulationJob),
		client:       &http.Client{Timeout: 30 * time.Second},
		gameWrappers: make(map[string]KeepGenerateWrapper),
	}
}

func (s *SimulatorService) RegisterGameWrapper(gameName string, wrapper KeepGenerateWrapper) *SimulatorService {
	s.gameWrappers[gameName] = wrapper
	return s
}

func (s *SimulatorService) WithKeepGenerate(f KeepGenerateWrapper) *SimulatorService {
	s.keepGenerateFunc = f

	return s
}

func (s *SimulatorService) applyGameSpecificWrapper(gameName string) func() {
	originalWrapper := s.keepGenerateFunc

	if gameWrapper, exists := s.gameWrappers[gameName]; exists {
		zap.S().Infof("Using game-specific wrapper for %s", gameName)
		s.keepGenerateFunc = gameWrapper
	}

	return func() {
		s.keepGenerateFunc = originalWrapper
	}
}

type Payload struct {
	Result     *SimulationView `json:"result"`
	Rtp        string          `json:"rtp"`
	Volatility string          `json:"volatility"`
	JobID      string          `json:"job_id"`
}

type SimulatorResult struct {
	JobID       string          `json:"job_id"`
	Game        string          `json:"game"`
	Result      *SimulationView `json:"result"`
	RTP         string          `json:"rtp"`
	Volatility  string          `json:"volatility"`
	Status      string          `json:"status"`
	Error       string          `json:"error,omitempty"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

func (j *SimulationJob) toResult(rtp, volatility string) *SimulatorResult {
	var gameName string
	if j.Result != nil {
		gameName = j.Result.Game
	}

	return &SimulatorResult{
		JobID:       j.ID,
		Game:        gameName,
		Result:      j.Result,
		RTP:         rtp,
		Volatility:  volatility,
		Status:      string(j.Status),
		Error:       j.Error,
		CompletedAt: j.CompletedAt,
	}
}

func (s *SimulatorService) SimulateV2(cfg *SimulatorConfig, rtp, volatility string) error {
	restore := s.applyGameSpecificWrapper(cfg.GameName)
	defer restore()

	s.generateParams = cfg.GenerateParams

	result, err := s.Simulate(cfg.GameName, cfg.Spins, cfg.Wager, cfg.Workers)
	if err != nil {
		return fmt.Errorf("simulator error: %w", err)
	}

	reportPages := []utils.Page{{
		Name:  "Report",
		Table: utils.Transpose(utils.ExtractTable([]*SimulationView{result.View()}, "xlsx")),
	}}

	return saveReport(cfg, reportPages, rtp, volatility)
}

func (s *SimulatorService) CreateSimulation(cfg *SimulatorConfig, rtp, volatility string) (string, error) {
	job := s.createJob(cfg)

	go func() {
		s.updateJobStatus(job.ID, StatusRunning, nil)
		initialResult := job.toResult(rtp, volatility)
		initialResult.Game = cfg.GameName

		restore := s.applyGameSpecificWrapper(cfg.GameName)

		s.generateParams = cfg.GenerateParams
		result, err := s.Simulate(cfg.GameName, cfg.Spins, cfg.Wager, cfg.Workers)

		restore()

		if err != nil {
			s.updateJobStatus(job.ID, StatusFailed, fmt.Errorf("failed to simulate: %w", err))
			failedResult := job.toResult(rtp, volatility)
			failedResult.Game = cfg.GameName
			_ = s.sendResult(cfg.CallbackURL, failedResult)
			return
		}

		s.jobsMux.Lock()
		job.Result = result.View()
		s.jobsMux.Unlock()

		s.updateJobStatus(job.ID, StatusCompleted, nil)
		completedResult := job.toResult(rtp, volatility)

		if err := s.sendResult(cfg.CallbackURL, completedResult); err != nil {
			s.updateJobStatus(job.ID, StatusFailed, fmt.Errorf("failed to send final result: %w", err))
			return
		}

		s.updateJobStatus(job.ID, StatusCompleted, nil)
	}()

	return job.ID, nil
}

func (s *SimulatorService) notifyGameSimulator(job *SimulationJob, rtp, volatility string) error {
	result := job.toResult(rtp, volatility)

	maxRetries := 3
	backoff := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		if err := s.sendResult(job.CallbackURL, result); err != nil {
			zap.S().Warnf("Attempt %d to send results failed: %v", attempt+1, err)
			time.Sleep(backoff)
			backoff *= 2
			continue
		}
		zap.S().Infof("Successfully sent results for job %s", job.ID)
		return nil
	}

	return fmt.Errorf("failed to send results after %d attempts", maxRetries)
}

func (s *SimulatorService) sendResult(url string, result *SimulatorResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *SimulatorService) GetJob(id string) (*SimulationJob, bool) {
	s.jobsMux.RLock()
	defer s.jobsMux.RUnlock()

	job, exists := s.jobs[id]
	return job, exists
}

func (s *SimulatorService) createJob(config *SimulatorConfig) *SimulationJob {
	job := &SimulationJob{
		ID:          uuid.New().String(),
		Status:      StatusPending,
		CallbackURL: config.CallbackURL,
		CreatedAt:   time.Now(),
	}

	s.jobsMux.Lock()
	s.jobs[job.ID] = job
	s.jobsMux.Unlock()

	return job
}

func (s *SimulatorService) updateJobStatus(jobID string, status SimulationStatus, err error) {
	s.jobsMux.Lock()
	defer s.jobsMux.Unlock()

	if job, exists := s.jobs[jobID]; exists {
		job.Status = status
		if err != nil {
			job.Error = err.Error()
		}
		if status == StatusCompleted || status == StatusFailed {
			now := time.Now()
			job.CompletedAt = &now
		}
	}
}

func saveReport(cfg *SimulatorConfig, reportPages []utils.Page, rtp, volatility string) error {
	excel, err := utils.ExportMultiPageXLSX(reportPages)
	if err != nil {
		return err
	}

	abs, err := filepath.Abs(cfg.ReportPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(abs, os.ModePerm); err != nil {
		return err
	}

	withDash := func(str string) string {
		return lo.Ternary(str != "", fmt.Sprintf("-%s", str), "")
	}

	filename := cfg.GameName
	filename += withDash(rtp)
	filename += withDash(volatility)
	filename += withDash(time.Now().UTC().Format("2006-01-02-15-04-05"))
	filename += ".xlsx"

	file, err := os.Create(filepath.Join(abs, filename))
	if err != nil {
		return err
	}

	if err = excel.Write(file); err != nil {
		return err
	}

	return nil
}

func (s *SimulatorService) Simulate(game string, count int64, wager int64, workersCount int) (*SimulationResult, error) {
	res := &SimulationResult{
		Wager: wager,
		Count: count,
		Game:  game,

		BaseAward:  new(big.Int),
		BonusAward: new(big.Int),
		Award:      new(big.Int),
		Spent:      new(big.Int),

		BaseAwardSquareSum:  new(big.Int),
		BonusAwardSquareSum: new(big.Int),
		AwardSquareSum:      new(big.Int),

		BaseAwardStandardDeviation:  new(big.Float),
		BonusAwardStandardDeviation: new(big.Float),
		AwardStandardDeviation:      new(big.Float),
	}

	type result struct {
		Wager          int64
		BaseAward      int64
		BonusAward     int64
		BonusTriggered bool
	}
	now := time.Now()
	bar := progressbar.NewOptions64(count,
		progressbar.OptionThrottle(200*time.Millisecond),
		progressbar.OptionSetDescription("Simulating..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println("\nTime elapsed:", time.Since(now))
		}),
	)

	var (
		factory = s.boot.SpinFactory

		inputCh  = make(chan int64, workersCount)
		outputCh = make(chan result, workersCount)
		errCh    = make(chan error, 1)

		wg = new(sync.WaitGroup)
	)

	worker := func(wg *sync.WaitGroup, inputCh <-chan int64, outputCh chan<- result) {
		defer wg.Done()

		var prevSpin engine.Spin

		for {
			if _, ok := <-inputCh; !ok {
				return
			}

			ctx := engine.Context{Context: context.Background()}
			if prevSpin != nil {
				ctx.LastSpin = prevSpin
			}

			spin, _, err := factory.Generate(ctx, wager, s.generateParams)
			if err != nil {
				errCh <- err

				return
			}

			if s.keepGenerateFunc != nil {
				spin, err = s.keepGenerateFunc(ctx, spin, factory)
				if err != nil {
					errCh <- err

					return
				}
			}

			prevSpin = spin

			outputCh <- result{
				Wager:          spin.Wager(),
				BaseAward:      spin.BaseAward(),
				BonusAward:     spin.BonusAward(),
				BonusTriggered: spin.BonusTriggered(),
			}
		}
	}

	go func() {
		defer close(inputCh)

		for i := int64(0); i < count; i++ {
			inputCh <- i
		}
	}()

	go func() {
		for i := 0; i < workersCount; i++ {
			wg.Add(1)
			go worker(wg, inputCh, outputCh)
		}

		wg.Wait()
		close(outputCh)
	}()

	i := 0

Loop:
	for {
		select {
		case output, ok := <-outputCh:
			if !ok {
				break Loop
			}

			award := output.BaseAward + output.BonusAward

			res.BaseAward.Add(res.BaseAward, big.NewInt(output.BaseAward))
			res.BonusAward.Add(res.BonusAward, big.NewInt(output.BonusAward))
			res.Award.Add(res.Award, big.NewInt(award))
			res.Spent.Add(res.Spent, big.NewInt(output.Wager))

			if award > res.MaxExposure {
				res.MaxExposure = award
			}

			if award > 0 {
				res.BaseAwardCount++
			}

			if output.BonusTriggered {
				res.BonusGameCount++
			}

			if award >= wager*1 {
				res.X1Count++
			}

			if award >= wager*10 {
				res.X10Count++
			}

			if award >= wager*100 {
				res.X100Count++
			}

			res.BaseAwardSquareSum.Add(res.BaseAwardSquareSum, big.NewInt(0).Mul(big.NewInt(output.BaseAward), big.NewInt(output.BaseAward)))
			res.BonusAwardSquareSum.Add(res.BonusAwardSquareSum, big.NewInt(0).Mul(big.NewInt(output.BonusAward), big.NewInt(output.BonusAward)))
			res.AwardSquareSum.Add(res.AwardSquareSum, big.NewInt(0).Mul(big.NewInt(award), big.NewInt(award)))

			_ = bar.Add(1) // ignore error

			i++
		case err := <-errCh:
			return nil, err
		}
	}

	baseMeanB := new(big.Float).SetInt(res.BaseAward)
	bonusMeanB := new(big.Float).SetInt(res.BonusAward)
	totalMeanB := new(big.Float).SetInt(res.Award)

	baseMeanB.Quo(baseMeanB, big.NewFloat(float64(res.Count)))
	bonusMeanB.Quo(bonusMeanB, big.NewFloat(float64(res.Count)))
	totalMeanB.Quo(totalMeanB, big.NewFloat(float64(res.Count)))

	res.BaseAwardStandardDeviation = StandardDeviation(res.BaseAwardSquareSum, res.BaseAward, baseMeanB, res.Count)
	res.BonusAwardStandardDeviation = StandardDeviation(res.BonusAwardSquareSum, res.BonusAward, bonusMeanB, res.Count)
	res.AwardStandardDeviation = StandardDeviation(res.AwardSquareSum, res.Award, totalMeanB, res.Count)

	res.Volatility = new(big.Float).Quo(res.AwardStandardDeviation, new(big.Float).SetInt64(wager))

	awardF := new(big.Float).SetInt(res.Award)
	baseF := new(big.Float).SetInt(res.BaseAward)
	bonusF := new(big.Float).SetInt(res.BonusAward)
	spentF := new(big.Float).SetInt(res.Spent)

	res.RTP, _ = new(big.Float).Quo(awardF, spentF).Float64()
	res.RTPBaseGame, _ = new(big.Float).Quo(baseF, spentF).Float64()
	res.RTPBonusGame, _ = new(big.Float).Quo(bonusF, spentF).Float64()

	return res, nil
}

type SimulationResult struct {
	Game        string   `xlsx:"Game"`
	Count       int64    `xlsx:"Count"`
	Wager       int64    `xlsx:"Wager"`
	Spent       *big.Int `xlsx:"Spent"`
	MaxExposure int64    `xlsx:"Max Exposure"`

	BaseAwardCount int64 `xlsx:"Base Award Count"`
	BonusGameCount int64 `xlsx:"Bonus Game Count"`

	X1Count   int64 `xlsx:"X1 Count"`
	X10Count  int64 `xlsx:"X10 Count"`
	X100Count int64 `xlsx:"X100 Count"`

	BaseAward  *big.Int `xlsx:"Base BaseAward"`
	BonusAward *big.Int `xlsx:"Bonus BaseAward"`
	Award      *big.Int `xlsx:"Award"`

	BaseAwardSquareSum  *big.Int `xlsx:"Base BaseAward Square Sum"`
	BonusAwardSquareSum *big.Int `xlsx:"Bonus BaseAward Square Sum"`
	AwardSquareSum      *big.Int `xlsx:"BaseAward Square Sum"`

	BaseAwardStandardDeviation  *big.Float `xlsx:"Base BaseAward Standard Deviation"`
	BonusAwardStandardDeviation *big.Float `xlsx:"Bonus BaseAward Standard Deviation"`
	AwardStandardDeviation      *big.Float `xlsx:"BaseAward Standard Deviation"`

	Volatility *big.Float `xlsx:"Volatility"`

	RTP          float64 `xlsx:"RTP"`
	RTPBaseGame  float64 `xlsx:"RTP Base Game"`
	RTPBonusGame float64 `xlsx:"RTP Bonus Game"`
}

func (r SimulationResult) View() *SimulationView {
	return &SimulationView{
		Game:        r.Game,
		Count:       fmt.Sprint(r.Count),
		Wager:       fmt.Sprint(r.Wager),
		Spent:       fmt.Sprint(r.Spent),
		MaxExposure: fmt.Sprint(r.MaxExposure),

		AwardCount: fmt.Sprint(r.BaseAwardCount),
		AwardRate:  countToRate(r.BaseAwardCount, r.Count),

		BonusGameCount: fmt.Sprint(r.BonusGameCount),
		BonusGameRate:  countToRate(r.BonusGameCount, r.Count),

		X1Count:   fmt.Sprint(r.X1Count),
		X10Count:  fmt.Sprint(r.X10Count),
		X100Count: fmt.Sprint(r.X100Count),

		X1Rate:   countToRate(r.X1Count, r.Count),
		X10Rate:  countToRate(r.X10Count, r.Count),
		X100Rate: countToRate(r.X100Count, r.Count),

		BaseAward:  r.BaseAward.String(),
		BonusAward: r.BonusAward.String(),
		Award:      r.Award.String(),

		BaseAwardSquareSum:  r.BaseAwardSquareSum.String(),
		BonusAwardSquareSum: r.BonusAwardSquareSum.String(),
		AwardSquareSum:      r.AwardSquareSum.String(),

		BaseAwardStandardDeviation:  float64FromBigFloat(r.BaseAwardStandardDeviation, 3),
		BonusAwardStandardDeviation: float64FromBigFloat(r.BonusAwardStandardDeviation, 3),
		AwardStandardDeviation:      float64FromBigFloat(r.AwardStandardDeviation, 3),

		Volatility: float64FromBigFloat(r.Volatility, 3),

		RTP:          floatWithPrecision(r.RTP),
		RTPBaseGame:  floatWithPrecision(r.RTPBaseGame),
		RTPBonusGame: floatWithPrecision(r.RTPBonusGame),
	}
}

type SimulationView struct {
	Game        string `json:"game" xlsx:"Game"`
	Count       string `json:"count" xlsx:"Count"`
	Wager       string `json:"wager" xlsx:"Wager"`
	Spent       string `json:"spent" xlsx:"Spent"`
	MaxExposure string `json:"max_exposure" xlsx:"Max Exposure"`

	NewLine1 string `xlsx:""`

	AwardCount string `json:"award_count" xlsx:"BaseAward Count"`
	AwardRate  string `json:"award_rate" xlsx:"AwardRate (Hit Rate)"`

	NewLine2 string `xlsx:""`

	BonusGameCount string `json:"bonus_game_count" xlsx:"Bonus Game Count"`
	BonusGameRate  string `json:"bonus_game_rate" xlsx:"Bonus Game Rate"`

	NewLine3 string `xlsx:""`

	X1Count   string `json:"x1_count" xlsx:"X1 Count"`
	X10Count  string `json:"x10_count" xlsx:"X10 Count"`
	X100Count string `json:"x100_count" xlsx:"X100 Count"`

	NewLine4 string `xlsx:""`

	X1Rate   string `json:"x1_rate" xlsx:"X1 Rate"`
	X10Rate  string `json:"x10_rate" xlsx:"X10 Rate"`
	X100Rate string `json:"x100_rate" xlsx:"X100 Rate"`

	NewLine5 string `xlsx:""`

	BaseAward  string `json:"base_award" xlsx:"Base BaseAward"`
	BonusAward string `json:"bonus_award" xlsx:"Bonus BaseAward"`
	Award      string `json:"award" xlsx:"BaseAward"`

	NewLine6 string `xlsx:""`

	BaseAwardSquareSum  string `json:"base_award_square_sum" xlsx:"Base BaseAward Square Sum"`
	BonusAwardSquareSum string `json:"bonus_award_square_sum" xlsx:"Bonus BaseAward Square Sum"`
	AwardSquareSum      string `json:"award_square_sum" xlsx:"BaseAward Square Sum"`

	NewLine7 string `xlsx:""`

	BaseAwardStandardDeviation  float64 `json:"base_award_standard_deviation" xlsx:"Base BaseAward Standard Deviation"`
	BonusAwardStandardDeviation float64 `json:"bonus_award_standard_deviation" xlsx:"Bonus BaseAward Standard Deviation"`
	AwardStandardDeviation      float64 `json:"award_standard_deviation" xlsx:"BaseAward Standard Deviation"`

	NewLine8 string `xlsx:""`

	Volatility float64 `json:"volatility" xlsx:"Volatility"`

	NewLine9 string `xlsx:""`

	RTP          string `json:"rtp" xlsx:"RTP"`
	RTPBaseGame  string `json:"rtp_base_game" xlsx:"RTP Base Game"`
	RTPBonusGame string `json:"rtp_bonus_game" xlsx:"RTP Bonus Game"`
}

func countToRate(count, total int64) string {
	div := float64(count) / float64(total)

	return floatWithPrecision(div) + "%"
}

func floatWithPrecision(f float64) string {
	return fmt.Sprintf("%.3f", f*100)
}

func StandardDeviation(squareSum, award *big.Int, mean *big.Float, count int64) *big.Float {
	f1 := new(big.Float).SetInt(squareSum)

	f2 := big.NewFloat(-2)
	f2.Mul(f2, mean)
	f2.Mul(f2, new(big.Float).SetInt(award))

	f3 := new(big.Float).SetInt64(int64(count))
	f3.Mul(f3, mean)
	f3.Mul(f3, mean)

	sd := big.NewFloat(0)

	sd.Add(sd, f1)
	sd.Add(sd, f2)
	sd.Add(sd, f3)

	sd.Quo(sd, new(big.Float).SetInt64(int64(count)))
	sd.Sqrt(sd)

	return sd
}

func float64FromBigFloat(float *big.Float, i int) float64 {
	value, _ := float.Float64()

	round := math.Pow(10, float64(i))

	return math.Round(value*round) / round
}
