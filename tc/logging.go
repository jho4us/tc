package tc

import (
	"github.com/go-kit/kit/log"

	"github.com/jho4us/tc/test"

	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) CreateTest(name string) (id test.TID, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "create",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CreateTest(name)
}

func (s *loggingService) LoadTest(id test.TID) (t Test, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "load",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.LoadTest(id)
}

func (s *loggingService) PutTest(t *Test) (err error) {
	defer func(begin time.Time) {
		id := "nil"
		if t != nil {
			id = t.ID
		}
		_ = s.logger.Log(
			"method", "put",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.PutTest(t)
}

func (s *loggingService) DeleteTest(id test.TID) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "delete",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.DeleteTest(id)
}

func (s *loggingService) Tests() []Test {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "list_tests",
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.Tests()
}
