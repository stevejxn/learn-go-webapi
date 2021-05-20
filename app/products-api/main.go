package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/stevejxn/learn-go-webapi/app/products-api/handlers"
)

const (
	ApiName = "PRODUCTS"
)

// version number of the program - to be set in part of ci/cd pipeline
var build = "develop"

func main() {
	logger := log.New(os.Stdout, "PRODUCTS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(logger); err != nil {
		logger.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	// ===
	// Configuration
	var cfg struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
	}
	cfg.SVN = build
	cfg.Version.Desc = "products api"

	if err := conf.Parse(os.Args[1:], ApiName, &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage(ApiName, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString(ApiName, &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// ===
	// App Starting
	log.Printf("main: Started: Application initializing: version %q", build)
	defer log.Println("main: Completed")

	configDetails, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config output for logging")
	}
	log.Printf("main: Config:\n%v\n", configDetails)

	// TODO: init auth, database, tracing, debug support

	// ===
	// Start API Service
	log.Println("main: Initializing API service")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// A buffered channel is required by the signal package
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      handlers.API(build, shutdown, log),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Make a channel to receive for errors coming from the listener.
	// Use a buffered channel so allow the goroutine to exit if the error is never collected
	serverErrors := make(chan error, 1)

	// Start the service listening for requests
	go func() {
		log.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// ===
	// Shutdown
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")
	case sig := <-shutdown:
		log.Printf("main: %v: Start shutdown", sig)

		// Give any outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			_ = api.Close()
			return errors.Wrap(err, "could not stop the server gracefully")
		}

		log.Printf("main: %v: Completed shutdown", sig)
	}

	return nil
}
