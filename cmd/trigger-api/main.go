package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/krostar/cleaner"
	"github.com/pkg/errors"

	"github.com/krostar/r10k-trigger/internal/pkg/app"
)

func main() {
	defer cleaner.Clean(onCleanError)

	var (
		flags                            = initFlags(os.Args[1:])
		cfg                              = initConfig(flags.configFile)
		log, statsEndpoint, statsHandler = initObs(cfg.Observability)
	)

	log.WithField("version", app.Version()).Info("starting app")

	var (
		r10kdeployer = initR10KDeployer(cfg.R10KDeployer)
		httpUsecases = initHTTPUsecases(r10kdeployer)
		srv          = initHTTP(cfg.HTTPServer, httpUsecases, log, statsEndpoint, statsHandler)
	)

	if err := srv.Run(cfg.GracefulShutdownTimeout, syscall.SIGINT, syscall.SIGTERM); err != nil {
		panic(errors.Wrap(err, "unable to start or stop the server"))
	}
}

func onCleanError(err error) {
	fmt.Fprintf(os.Stderr, "a fatal error occured: %v\n", err)
	os.Exit(2)
}
