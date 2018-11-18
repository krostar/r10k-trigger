package handler

import (
	"context"
	"strings"
)

// DeployUsecases defines the trigger handler usecases.
//go:generate mockery -inpkg -testonly -name DeployUsecases
type DeployUsecases interface {
	DeployR10KEnvAsync(ctx context.Context, env string, onSuccess func(string), onError func(string, error))
}

func getEnvironmentFromGITRef(ref string) string {
	const refWantedPrefix = "refs/head/"

	tmp := strings.Split(ref, refWantedPrefix)
	if len(tmp) != 2 {
		return ""
	}

	return tmp[1]
}
