package obs

import (
	"net/http"

	"github.com/krostar/logger"
	"github.com/pkg/errors"

	"github.com/krostar/r10k-trigger/internal/pkg/app"
)

// Config stores the whole observability related configuration.
type Config struct {
	Logger  logger.Config `json:"logger"  yaml:"logger"`
	Tracer  TracerConfig  `json:"tracer"  yaml:"tracer"`
	Monitor MonitorConfig `json:"monitor" yaml:"monitor"`
}

// Init initializes all the components related to observability.
func Init(cfg Config) (logger.Logger, string, http.HandlerFunc, func(), error) {
	log, flushLogs, errLogger := initLogger(cfg.Logger)
	if errLogger != nil {
		return nil, "", nil, nil, errors.Wrap(errLogger, "logger init failed")
	}

	stopTracer, errTracer := initTracer(cfg.Tracer)
	if errTracer != nil {
		return log, "", nil, flushLogs, errors.Wrap(errTracer, "tracer init failed")
	}

	monitorEndpoint, monitorHandler, stopMonitor, errMonitor := initMonitor(
		cfg.Monitor, app.AlphaNumericName(),
		logger.WriterLevel(log.WithField("source", "monitor"), logger.LevelError),
	)
	if errMonitor != nil {
		return log, "", nil, func() {
			stopTracer()
			flushLogs()
		}, errors.Wrap(errMonitor, "monitor init failed")
	}

	return log, monitorEndpoint, monitorHandler, func() {
		stopTracer()
		stopMonitor()
		flushLogs()
	}, nil
}
