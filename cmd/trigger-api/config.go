package main

import (
	"time"

	"github.com/krostar/config"
	sourceenv "github.com/krostar/config/source/env"
	sourcefile "github.com/krostar/config/source/file"
	"github.com/pkg/errors"

	"github.com/krostar/r10k-trigger/internal/pkg/app"
	"github.com/krostar/r10k-trigger/internal/pkg/obs"
	"github.com/krostar/r10k-trigger/internal/trigger-api/delivery/httpapi"
	"github.com/krostar/r10k-trigger/pkg/r10kshelldeployer"
)

type appConfig struct {
	Observability           obs.Config         `json:"observability"             yaml:"observability"`
	R10KDeployer            r10kDeployerConfig `json:"r10k-deployer"             yaml:"r10k-deployer"`
	HTTPServer              httpapi.Config     `json:"http-server"               yaml:"http-server"`
	GracefulShutdownTimeout time.Duration      `json:"graceful-shutdown-timeout" yaml:"graceful-shutdown-timeout"`
}

type r10kDeployerConfig struct {
	Shell r10kshelldeployer.Config `json:"shell" yaml:"shell"`
}

func (c *appConfig) SetDefault() {
	c.GracefulShutdownTimeout = 10 * time.Second
}

func initConfig(configFile string) *appConfig {
	var cfg appConfig

	if err := config.New(config.WithSources(
		sourcefile.New(configFile, sourcefile.FailOnUnknownFields()),
		sourceenv.New(app.AlphaNumericName()),
	)).Load(&cfg); err != nil {
		panic(errors.Wrap(err, "unable to load config"))
	}

	if err := config.Validate(&cfg); err != nil {
		panic(errors.Wrap(err, "config is invalid"))
	}

	return &cfg
}
