package stress

import (
	"github.com/samber/lo"
	"sync"
	"time"
)

type SpamReport struct {
	Errs            []error
	InitStart       time.Time
	InitEnd         time.Time
	WageringStart   time.Time
	WageringEnd     time.Time
	InitLatencies   []time.Duration
	UpdateLatencies []time.Duration
	WagerLatencies  []time.Duration

	errMu, wagerLatMu, initLatMu, updateLatMu sync.Mutex
}

func (s *SpamReport) addError(err error) {
	s.errMu.Lock()
	s.Errs = append(s.Errs, err)
	s.errMu.Unlock()
}

func (s *SpamReport) startInitLatency() func() {
	now := time.Now()

	return func() {
		lat := time.Since(now)

		s.initLatMu.Lock()
		s.InitLatencies = append(s.InitLatencies, lat)
		s.initLatMu.Unlock()
	}
}

func (s *SpamReport) startWagerLatency() func() {
	now := time.Now()

	return func() {
		lat := time.Since(now)

		s.wagerLatMu.Lock()
		s.WagerLatencies = append(s.WagerLatencies, lat)
		s.wagerLatMu.Unlock()
	}
}

func (s *SpamReport) startUpdateLatencies() func() {
	now := time.Now()

	return func() {
		lat := time.Since(now)

		s.updateLatMu.Lock()
		s.UpdateLatencies = append(s.UpdateLatencies, lat)
		s.updateLatMu.Unlock()
	}
}

func (s *SpamReport) InitDuration() time.Duration {
	return s.InitEnd.Sub(s.InitStart)
}

func (s *SpamReport) WageringDuration() time.Duration {
	return s.WageringEnd.Sub(s.WageringStart)
}

func (s *SpamReport) AvgInitLatency() time.Duration {
	return avgLatency(s.InitLatencies)
}

func (s *SpamReport) AvgWagerLatency() time.Duration {
	return avgLatency(s.WagerLatencies)
}

func (s *SpamReport) AvgUpdateLatency() time.Duration {
	return avgLatency(s.UpdateLatencies)
}

func avgLatency(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	return lo.Sum(latencies) / time.Duration(len(latencies))
}
