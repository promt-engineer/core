package stress

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type RequestGenerator interface {
	WagerParams() interface{}
	UpdateParams(spin engine.Spin, restoring engine.RestoringIndexes) interface{}
}

type Spammer struct {
	rg          RequestGenerator
	spinFactory engine.SpinFactory

	logger  *zap.Logger
	client  http.Client
	urlBase url.URL
	game    string
}

func NewSpammer(rg RequestGenerator, spinFactory engine.SpinFactory, logger *zap.Logger, urlBase url.URL, game string) *Spammer {
	logger.Sugar().Info("Building spammer")

	return &Spammer{
		rg:          rg,
		spinFactory: spinFactory,

		logger:  logger,
		client:  http.Client{Transport: http.DefaultTransport},
		urlBase: urlBase,
		game:    game}
}

func (s *Spammer) Run(userCount, requestPerUser int) (SpamReport, error) {
	s.logger.Sugar().Info("Start spamming")
	initBar := progressbar.Default(int64(userCount), "init")

	rep := SpamReport{
		InitStart: time.Now(),
	}
	mu := &sync.Mutex{}

	initRequests, err := s.InitUsersRequests(userCount)
	if err != nil {
		return rep, err
	}

	states := []*State{}
	wg := &sync.WaitGroup{}

	for _, req := range initRequests {
		wg.Add(1)

		go func(r *http.Request) {
			defer wg.Done()
			defer initBar.Add(1)

			endInit := rep.startInitLatency()
			resp, err := s.client.Do(r)
			endInit()

			if err != nil {
				rep.addError(err)

				return
			}

			state, err := ParseStateResponse(resp)
			if err != nil {
				rep.addError(err)

				return
			}

			mu.Lock()
			states = append(states, state)
			mu.Unlock()
		}(req)
	}

	wg.Wait()

	rep.InitEnd = time.Now()
	s.logger.Sugar().Infof("Init took %v", rep.InitDuration())

	if len(rep.Errs) != 0 {
		return rep, errors.Join(rep.Errs...)
	}

	s.logger.Sugar().Info("Start wagering")

	wageringBar := progressbar.Default(int64(userCount*requestPerUser), "wagering")
	rep.WageringStart = time.Now()

	for _, state := range states {
		wg.Add(1)

		go func(state *State) {
			defer wg.Done()

			for i := 0; i < requestPerUser; i++ {
				wageringBar.Add(1)

				req, err := WagerRequest(s.urlBase, state.SessionToken, DefaultWager, s.rg.WagerParams())
				if err != nil {
					rep.addError(err)

					return
				}

				endWager := rep.startWagerLatency()
				resp, err := s.client.Do(req)
				endWager()

				if err != nil {
					rep.addError(err)

					return
				}

				state, err = ParseStateResponse(resp)
				if err != nil {
					rep.addError(err)

					return
				}

				spinBytes, err := json.Marshal(state.GameResult.Spin)
				if err != nil {
					rep.addError(err)

					return
				}

				restoringBytes, err := json.Marshal(state.GameResult.RestoringIndexes)
				if err != nil {
					rep.addError(err)

					return
				}

				spin, err := s.spinFactory.UnmarshalJSONSpin(spinBytes)
				if err != nil {
					rep.addError(err)

					return
				}

				restoring, err := s.spinFactory.UnmarshalJSONRestoringIndexes(restoringBytes)
				if err != nil {
					rep.addError(err)

					return
				}

				req, err = UpdateRestoringRequest(s.urlBase, state.SessionToken, s.rg.UpdateParams(spin, restoring))
				if err != nil {
					rep.addError(err)

					return
				}

				endUpdate := rep.startUpdateLatencies()
				resp, err = s.client.Do(req)
				endUpdate()

				if err != nil {
					rep.addError(err)

					return
				}

				if resp.StatusCode > 299 || resp.StatusCode < 200 {
					rep.addError(fmt.Errorf("bad status code: %d", resp.StatusCode))

					return
				}
			}
		}(state)
	}

	wg.Wait()

	rep.WageringEnd = time.Now()
	s.logger.Sugar().Infof("Wagering took %v", rep.WageringDuration())

	return rep, nil
}

func (s *Spammer) InitUsersRequests(userCount int) ([]*http.Request, error) {
	reqs := make([]*http.Request, userCount)
	var err error

	for i := 0; i < userCount; i++ {
		reqs[i], err = GenerateInitRequest(s.urlBase, s.game)
		if err != nil {
			return nil, err
		}
	}

	return reqs, nil
}
