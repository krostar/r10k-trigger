package httpapi

import "github.com/krostar/r10k-trigger/internal/trigger-api/delivery/httpapi/handler"

// Usecases defines all the usecases required by handlers.
//go:generate mockery -inpkg -testonly -name Usecases
type Usecases interface {
	handler.DeployUsecases
}
