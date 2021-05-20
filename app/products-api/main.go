package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
)

const (
	ApiName = "PRODUCTS"
)

// version number of the program - to be set in part of ci/cd pipeline
var build = "develop"

func main() {
	log := log.New(os.Stdout, "PRODUCTS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	// ===
	// Configuration
	var cfg struct {
		conf.Version
		Web struct {
			APIHost         string `conf:"default:0.0.0.0:3000"`
			ReadTimeout     string `conf:"default:5s"`
			WriteTimeout    string `conf:"default:5s"`
			ShutdownTimeout string `conf:"default:5s"`
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

	return nil
}
