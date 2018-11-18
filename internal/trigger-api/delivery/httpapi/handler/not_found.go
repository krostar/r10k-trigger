package handler

import (
	"net/http"

	"github.com/krostar/httpw"
)

// NotFound handles unhandled routes.
func NotFound(r *http.Request) (*httpw.R, error) {
	return nil, &httpw.E{Status: http.StatusNotFound}
}
