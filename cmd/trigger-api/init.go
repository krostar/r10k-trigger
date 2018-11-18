package main

import (
	"net/http"

	"github.com/krostar/cleaner"
	"github.com/krostar/logger"
	"github.com/pkg/errors"

	"github.com/krostar/r10k-trigger/internal/pkg/obs"
	"github.com/krostar/r10k-trigger/internal/trigger-api/delivery/httpapi"
	"github.com/krostar/r10k-trigger/internal/trigger-api/usecase"
	"github.com/krostar/r10k-trigger/pkg/r10kshelldeployer"
)

func initObs(cfg obs.Config) (logger.Logger, string, http.HandlerFunc) {
	logger, statsEndpoint, statsHandler, stopFunc, err := obs.Init(cfg)
	cleaner.Add(stopFunc)
	if err != nil {
		panic(errors.Wrap(err, "unable to initialize tracer"))
	}

	return logger, statsEndpoint, statsHandler
}

func initR10KDeployer(cfg r10kDeployerConfig) *r10kshelldeployer.Shell {
	return r10kshelldeployer.New(
		r10kshelldeployer.WithConfig(&cfg.Shell),
	)
}

func initHTTPUsecases(deployer usecase.R10KEnvDeployer) httpapi.Usecases {
	return struct {
		*usecase.R10KDeploy
	}{
		R10KDeploy: &usecase.R10KDeploy{Deployer: deployer},
	}
}

func initHTTP(
	cfg httpapi.Config, usecases httpapi.Usecases, log logger.Logger,
	statsEndpoint string, statsHandler http.HandlerFunc,
) *httpapi.HTTP {
	http, err := httpapi.New(cfg, usecases, log,
		httpapi.WithStatsHandler(statsEndpoint, statsHandler),
	)
	if err != nil {
		panic(errors.Wrap(err, "unable to create http server"))
	}
	return http
}
