package tc

import (
	"github.com/go-kit/kit/metrics"

	"github.com/jho4us/tc/test"

	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) CreateTest(name string) (test.ID, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "create").Add(1)
		s.requestLatency.With("method", "create").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.CreateTest(name)
}

func (s *instrumentingService) LoadTest(id test.ID) (Test, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "load").Add(1)
		s.requestLatency.With("method", "load").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.LoadTest(id)
}

func (s *instrumentingService) PutTest(t *Test) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "put").Add(1)
		s.requestLatency.With("method", "put").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PutTest(t)
}

func (s *instrumentingService) DeleteTest(id test.ID) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "delete").Add(1)
		s.requestLatency.With("method", "delete").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.DeleteTest(id)
}

func (s *instrumentingService) Tests() []Test {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list_tests").Add(1)
		s.requestLatency.With("method", "list_tests").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Tests()
}
