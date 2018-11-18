package httpapi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithStatsHandler(t *testing.T) {
	var api HTTP

	WithStatsHandler("/metrics", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))(&api)
	assert.NotEmpty(t, api.statsEndpoint)
	assert.NotNil(t, api.statsHandler)
}
