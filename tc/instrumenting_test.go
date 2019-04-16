package tc

import (
	"testing"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func TestInstrumenting(t *testing.T) {
	rep, err := testCreateTestRepo(t)
	if err != nil {
		t.Fatal(err)
	}
	fieldKeys := []string{"method"}

	tcs := NewService(rep)
	tcs = NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "tc_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "tc_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		tcs,
	)

	_, err = tcs.CreateTest("")
	if err == nil {
		t.Fatal(err)
	}
	_, err = tcs.LoadTest("")
	if err == nil {
		t.Fatal(err)
	}
	_, err = tcs.LoadTest("unexistent")
	if err == nil {
		t.Fatal(err)
	}
	err = tcs.DeleteTest("")
	if err == nil {
		t.Fatal(err)
	}
	ta := tcs.Tests()
	if ta == nil {
		t.Errorf("cant query all tests")
	}
	testInsertTestAndGetTest(tcs, t, 0)
	testUpdateAndGetTest(tcs, t)
}
