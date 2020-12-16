package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"scbunn.org/tmp/gps-tracking-service/pkg/models"
	"scbunn.org/tmp/gps-tracking-service/pkg/models/inmem"
	"scbunn.org/tmp/gps-tracking-service/pkg/service"
)

func main() {

	addr := flag.String("addr", ":5000", "HTTP network address")
	datastore := flag.String("datastore", "inmemdb", "backend datastore to use")
	ttl := flag.Duration("Object TTL", 60*time.Second, "TTL of Object Telemetry")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldUnit = time.Second
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123})

	var db models.TelemetryReaderWriterChecker

	switch *datastore {
	case "inmemdb":
		db = createInMemoryDatabase(*ttl)
	case "redis":
		log.Fatal().Str("datastore", *datastore).Err(errors.New("datastore not implemented")).Msg("")
	default:
		log.Fatal().Str("datastore", *datastore).Err(errors.New("unknown datastore")).Msg("")
	}
	log.Info().Str("datastore", *datastore).Msg("datastore created")

	log.Info().Msg("starting location tracking service")
	service := service.New(*addr, db, &log.Logger)

	svr := http.Server{
		Addr:         *addr,
		Handler:      service.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("host", *addr).Msg("starting http server")
		err := svr.ListenAndServe()
		if err != nil {
			log.Fatal().Err(err)
		}
	}()

	// Trap signals so we can get a clean exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Block until we get a shutdown signal
	sig := <-c
	log.Info().Str("signal", sig.String()).Msg("recived a signal to shutdown...")

	// Try and shutdown the telemety service cleanly
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

}

func createInMemoryDatabase(expire time.Duration) *inmem.InMemoryDB {
	dbLogger := log.With().Str("component", "database").Logger()
	memdb := inmem.New(&dbLogger)

	dbLogger.Info().Dur("Object TTL", expire).Msg("Starting expiration goroutine")
	go func() {
		tick := time.Tick(expire)
		for range tick {
			memdb.Expire()
		}
	}()

	return memdb
}

func init() {
	// Register the prometheus build info collector
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}
