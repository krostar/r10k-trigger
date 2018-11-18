package httpapi

import (
	"net/http"
)

// Option defines a function prototype to apply options to the http instance.
type Option func(h *HTTP)

// WithStatsHandler set the stat endpoint handler to the http instance.
func WithStatsHandler(endpoint string, handler http.HandlerFunc) Option {
	return func(h *HTTP) {
		h.statsEndpoint = endpoint
		h.statsHandler = handler
	}
}
