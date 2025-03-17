package handlers

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/services"
	httpPackage "bitbucket.org/play-workspace/base-slot-server/pkg/kernel/transport/http"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/validator"
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
	"os/exec"
	"reflect"
)

type Request struct {
	GameName       string                 `json:"gameName" validate:"required"`
	ReportPath     string                 `json:"reportPath" validate:"required"`
	Spins          int64                  `json:"spins" validate:"required"`
	Wager          int64                  `json:"wager" validate:"required"`
	Workers        int                    `json:"workers" validate:"required"`
	GenerateParams map[string]interface{} `json:"generateParams"`
	RTP            string                 `json:"rtp" validate:"required"`
	Volatility     string                 `json:"volatility" validate:"required"`
	CallbackURL    string                 `json:"callbackURL" validate:"required"`
}

var gamesWithBonusChoice = []string{
	"cleos-riches-flexiways",
	"fortune-777-respin",
	"coral-reef-flexiways",
}

const (
	configPath   = "/app/config.yml"
	appPath      = "/app/app"
	tempFilePath = "/app/temp.yml"
)

type simulatorHandler struct {
	ctn              di.Container
	validationEngine *validator.Validator
	simulatorService *services.SimulatorService
}

func NewSimulatorHandler(ctn di.Container, validationEngine *validator.Validator, service *services.SimulatorService) httpPackage.Handler {
	rand := ctn.Get(constants.RNGName).(rng.Client)
	bonusChoiceWrapper := createBonusChoiceWrapper(rand)

	for _, game := range gamesWithBonusChoice {
		service.RegisterGameWrapper(game, bonusChoiceWrapper)
	}

	return &simulatorHandler{
		ctn:              ctn,
		validationEngine: validationEngine,
		simulatorService: service,
	}
}

func createBonusChoiceWrapper(rand rng.Client) services.KeepGenerateWrapper {
	return func(ctx engine.Context, spin engine.Spin, factory engine.SpinFactory) (engine.Spin, error) {
		sp, ok := spin.(interface{})
		if !ok {
			return spin, nil
		}

		v := reflect.ValueOf(sp)
		if v.Kind() == reflect.Ptr && v.IsValid() {
			v = v.Elem()

			bonusChoiceField := v.FieldByName("BonusChoice")
			if bonusChoiceField.IsValid() && !bonusChoiceField.IsNil() && bonusChoiceField.Len() > 0 {
				ctx.LastSpin = spin

				index, err := rand.Rand(uint64(bonusChoiceField.Len()))
				if err != nil {
					return nil, err
				}

				bonusChoice := bonusChoiceField.Index(int(index)).Interface()

				newSpin, ok, err := factory.KeepGenerate(ctx, bonusChoice)
				if err != nil {
					return nil, err
				}

				if ok {
					spin = newSpin
				}
			}
		}

		return spin, nil
	}
}

func (s *simulatorHandler) Shutdown() {}

func (s *simulatorHandler) Register(router *gin.RouterGroup) {
	router.POST("simulator", s.RunSimulate)
	router.GET("simulator/:jobId/status", s.GetSimulationStatus)
}

func (s *simulatorHandler) RunSimulate(ctx *gin.Context) {
	req, err := s.parseAndValidateRequests(ctx)
	if err != nil {
		httpPackage.BadRequest(ctx, err, nil)
		return
	}

	cfg := &services.SimulatorConfig{
		GameName:       req.GameName,
		ReportPath:     req.ReportPath,
		Spins:          req.Spins,
		Wager:          req.Wager,
		Workers:        req.Workers,
		GenerateParams: req.GenerateParams,
		CallbackURL:    req.CallbackURL,
	}

	jobID, err := s.simulatorService.CreateSimulation(cfg, req.RTP, req.Volatility)
	if err != nil {
		zap.S().Errorf("failed to start simulation: %v", err)
		httpPackage.ServerError(ctx, fmt.Errorf("failed to start simulation: %w", err), nil)
		return
	}

	response := map[string]interface{}{
		"job_id":    jobID,
		"status":    string(services.StatusPending),
		"game_name": req.GameName,
	}

	httpPackage.OK(ctx, "Simulation started", response)
}

func (s *simulatorHandler) GetSimulationStatus(ctx *gin.Context) {
	jobID := ctx.Param("jobId")

	job, exists := s.simulatorService.GetJob(jobID)
	if !exists {
		httpPackage.NotFound(ctx, fmt.Errorf("job not found"), nil)
		return
	}

	response := map[string]interface{}{
		"job_id":    job.ID,
		"status":    string(job.Status),
		"game_name": job.Result.Game,
	}

	if job.Error != "" {
		response["error"] = job.Error
	}

	httpPackage.OK(ctx, "Job status retrieved", response)
}

func (s *simulatorHandler) parseAndValidateRequests(ctx *gin.Context) (*Request, error) {
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewInternalValidationError(err)
	}

	if err := s.validationEngine.ValidateStruct(&req); err != nil {
		return nil, errs.NewInternalValidationError(err)
	}

	if _, err := url.Parse(req.CallbackURL); err != nil {
		return nil, errs.NewInternalValidationError(fmt.Errorf("invalid callback URL: %w", err))
	}

	return &req, nil
}

func (s *simulatorHandler) updateConfigWithRequest(originConfig map[string]interface{}, req *Request) (string, error) {
	generateParams := make(map[string]interface{})
	for key, value := range req.GenerateParams {
		generateParams[key] = value
	}
	originConfig["simulator"] = map[string]interface{}{
		"gameName":       req.GameName,
		"reportPath":     req.ReportPath,
		"spins":          req.Spins,
		"wager":          req.Wager,
		"workers":        req.Workers,
		"GenerateParams": generateParams,
	}

	originConfig["engine"] = map[string]interface{}{
		"rtp":               req.RTP,
		"volatility":        req.Volatility,
		"buildVersion":      1,
		"isCheatsAvailable": true,
		"mockRng":           false,
	}

	configBytes, err := yaml.Marshal(originConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	return string(configBytes), nil
}

func createTempConfigFile(configYaml string) (string, error) {
	zap.S().Info("Starting to create temp config file")

	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		zap.S().Error("Error when creating a temporary file", "error", err)
		return "", fmt.Errorf("could not create temp file: %w", err)
	}

	zap.S().Info("Temporary file created", "file", tempFile.Name())

	if _, err := tempFile.Write([]byte(configYaml)); err != nil {
		zap.S().Error("Error when writing to a file", "error", err)
		return "", fmt.Errorf("could not write to temp file: %w", err)
	}

	zap.S().Info("Successfully wrote to temp file")

	if err := tempFile.Close(); err != nil {
		zap.S().Error("Error when closing a file", "error", err)
		return "", fmt.Errorf("could not close temp file: %w", err)
	}

	zap.S().Info("Temporary config file created and closed", "file", tempFile.Name())
	return tempFile.Name(), nil
}

func (s *simulatorHandler) runSingleSimulation(req *Request) error {
	originConfig := loadConfig(configPath)

	updatedConfig, err := s.updateConfigWithRequest(originConfig, req)
	if err != nil {
		zap.S().Error("simulator: failed to update config", err)
		return err
	}

	tempFileName, err := createTempConfigFile(updatedConfig)
	if err != nil {
		zap.S().Error("Failed to create temporary config file", "error", err)
		return err
	}
	defer os.Remove(tempFileName)

	cmd := exec.Command(appPath, fmt.Sprintf("--config=%s", tempFileName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func loadConfig(path string) map[string]interface{} {
	config := make(map[string]interface{})

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}
