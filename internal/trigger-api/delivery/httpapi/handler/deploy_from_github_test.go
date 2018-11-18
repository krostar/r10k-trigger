package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/krostar/httpw"
	"github.com/krostar/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/krostar/r10k-trigger/internal/trigger-api/delivery/httpapi/handler"
)

func TestHandler_DeployFromGithub(t *testing.T) {
	var tests = map[string]struct {
		requestPayload         string
		mockDeployR10KAsyncEnv string
		expectedStatus         int
	}{
		"deploy successfull": {
			requestPayload:         `{ "ref": "refs/head/toto" }`,
			mockDeployR10KAsyncEnv: "toto",
			expectedStatus:         http.StatusAccepted,
		},
		"wrong payload syntax": {
			requestPayload: `{`,
			expectedStatus: http.StatusBadRequest,
		},
		"empty payload": {
			requestPayload: ``,
			expectedStatus: http.StatusBadRequest,
		},
		"empty payload environment": {
			requestPayload: `{ "ref": "" }`,
			expectedStatus: http.StatusNoContent,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				usecases handler.MockDeployUsecases
				r, _     = http.NewRequest("", "", strings.NewReader(test.requestPayload))
				w        = httptest.NewRecorder()
			)

			// prepare the mock
			if test.mockDeployR10KAsyncEnv != "" {
				usecases.
					On("DeployR10KEnvAsync", mock.Anything, test.mockDeployR10KAsyncEnv, mock.Anything, mock.Anything).
					Return().
					Once()
			}

			// call the handler
			httpw.Wrap(handler.DeployFromGithub(logger.Noop{}, &usecases)).ServeHTTP(w, r)

			// match expectations
			assert.Equal(t, test.expectedStatus, w.Code)
			assert.Empty(t, w.Body)
			usecases.AssertExpectations(t)
		})
	}
}
