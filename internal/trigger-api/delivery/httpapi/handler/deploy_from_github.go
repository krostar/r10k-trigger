package handler

import (
	"encoding/json"
	"net/http"

	"github.com/krostar/httpw"
	"github.com/krostar/logger"
	"github.com/krostar/logger/logmid"
)

// DeployFromGithub deploies a r10k environment from a github trigger.
func DeployFromGithub(log logger.Logger, usecase DeployUsecases) httpw.HandlerFunc {
	return func(r *http.Request) (*httpw.R, error) {
		var (
			ctx     = r.Context()
			payload struct {
				Ref string `json:"ref"`
			}
		)

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			return nil, &httpw.E{
				Status: http.StatusBadRequest,
				Err:    err,
			}
		}

		environment := getEnvironmentFromGITRef(payload.Ref)
		logmid.AddFieldInContext(ctx, "environment", environment)
		if environment == "" {
			return &httpw.Response{Status: http.StatusNoContent}, nil
		}

		usecase.DeployR10KEnvAsync(ctx, environment, nil,
			func(environment string, err error) { // called on error
				log.
					WithField("environment", environment).
					WithError(err).
					Error("r10k deployment failed")
			},
		)

		return &httpw.R{Status: http.StatusAccepted}, nil
	}
}
