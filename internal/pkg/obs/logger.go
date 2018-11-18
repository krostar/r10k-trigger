package obs

import (
	"github.com/krostar/logger"
	"github.com/krostar/logger/zap"
	"github.com/pkg/errors"
)

func initLogger(cfg logger.Config) (logger.Logger, func(), error) {
	log, flush, err := zap.New(zap.WithConfig(cfg))
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create logger")
	}
	restore := logger.RedirectStdLog(log, logger.LevelError)

	return log, func() {
		restore()
		flush() // nolint: errcheck, gosec
	}, nil
}
