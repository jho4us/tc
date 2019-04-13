package main

import (
	"flag"
	"fmt"

	"github.com/jho4us/tc/handlers"
	"github.com/jho4us/tc/repo"
	"github.com/jho4us/tc/tc"
	"github.com/jho4us/tc/test"

	"github.com/kelseyhightower/envconfig"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
)

var (
	REPOSITORY = "UNKNOWN"
	REVISION   = "UNKNOWN"
)

type Specification struct {
	Port     string `envconfig:"port" default:"8081"`                  // TC_PORT
	DbConfig string `envconfig:"db_config" default:"./db-config.json"` // TC_DB_CONFIG
}

func main() {
	var spec Specification

	err := envconfig.Process("tc", &spec)
	if err != nil {
		panic(err)
	}
	var (
		httpAddr = flag.String("http.addr", ":"+spec.Port, "HTTP listen address")
		dbConfig = flag.String("db.config", spec.DbConfig, "Database configuration file")
	)

	flag.Parse()

	var logger log.Logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var tests test.Repository

	tests, err = repo.NewTestRepository(*dbConfig)
	if err != nil {
		panic(err)
	}
	fieldKeys := []string{"method"}

	var tcs tc.Service = tc.NewService(tests)
	tcs = tc.NewLoggingService(log.With(logger, "component", "tc"), tcs)
	tcs = tc.NewInstrumentingService(
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

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/tc/v1/", tc.MakeHandler(tcs, httpLogger))
	mux.Handle("/resilience/", handlers.MakeHandler("/resilience", REPOSITORY, REVISION))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		_ = logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
