package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krostar/httpw"
	"github.com/stretchr/testify/assert"

	"github.com/krostar/r10k-trigger/internal/trigger-api/delivery/httpapi/handler"
)

func TestHandler_notFound(t *testing.T) {
	var (
		r, _ = http.NewRequest("", "", nil)
		w    = httptest.NewRecorder()
	)

	// call the handler
	httpw.WrapF(handler.NotFound).ServeHTTP(w, r)

	// match expectations
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Empty(t, w.Body)
}
