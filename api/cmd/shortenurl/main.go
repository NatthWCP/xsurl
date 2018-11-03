package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"xsurl/api/inmem"
	"xsurl/api/server"
	"xsurl/api/shortening"
	"xsurl/api/shortenurl"
)

const (
	defaultPort              = "8080"
	defaultRoutingServiceURL = "http://localhost:1234"
	defaultDBName            = "xsurl"
)

func main() {
	var (
		addr = envString("PORT", defaultPort)
		// rsURL = envString("ROUTING_SERVICE_URL", defaultRoutingServiceURL)

		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		// routingServiceURL = flag.String("service.routing", rsURL, "routing service URL")
		inmemory = flag.Bool("inmem", true, "use in-memory repositories")
	)
	flag.Parse()

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	var urls shortenurl.URLRepository

	if *inmemory {
		urls = inmem.NewURLRepository()
	} else {
		// TODO: mongodb setup
	}

	fieldKeys := []string{"method"}

	var ss shortening.Service
	ss = shortening.NewService(urls)
	ss = shortening.NewLoggingService(kitlog.With(logger, "component", "shortening"), ss)
	ss = shortening.NewInstrumentingService(
		kitprometheus.NewCounterFrom(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "shortening_service",
			Name:      "request_count",
			Help:      "Number of request received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "shortening_service",
			Name:      "request_latenacy_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		ss,
	)

	srv := server.New(ss, kitlog.With(logger, "component", "http"))

	errs := make(chan error, 2)

	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, srv)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
