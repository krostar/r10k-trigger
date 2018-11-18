package httpapi

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/krostar/httpx"
	"github.com/krostar/logger"
	"github.com/pkg/errors"
)

// HTTP handles the http transport layer.
type HTTP struct {
	cfg Config
	log logger.Logger

	statsEndpoint string
	statsHandler  http.HandlerFunc

	listener net.Listener
	server   *http.Server
}

// New returns a new http instance.
func New(cfg Config, usecases Usecases, log logger.Logger, opts ...Option) (*HTTP, error) {
	var (
		err error
		api = HTTP{cfg: cfg, log: log}
	)

	for _, opt := range opts {
		opt(&api)
	}

	api.server = httpx.NewServer(api.initRouter(usecases))
	api.server.ErrorLog = logger.StdLog(api.log.WithField("source", "http-error"), logger.LevelWarn)

	api.listener, err = httpx.NewListener(cfg.ListenAddress)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create listener")
	}

	return &api, nil
}

// Run starts and gracefully stops the running http server.
func (h *HTTP) Run(timeout time.Duration, stopSignals ...os.Signal) error {
	return httpx.StartAndStopWithSignal(h.server, h.listener, timeout, stopSignals...)
}
