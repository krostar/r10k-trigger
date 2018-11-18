package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureGithubOrigin(t *testing.T) {
	var (
		secret  = "42universalanswer"
		bodyRaw = []byte(`{"hello": "world"}`)
		tests   = map[string]struct {
			headers        map[string]string
			expectedStatus int
		}{
			"success call": {
				headers: map[string]string{
					"User-Agent":      "GitHub-Hookshot/test",
					"X-GitHub-Event":  "push",
					"Content-Type":    "application/json",
					"X-Hub-Signature": "sha1=96ce0cdaf58b128b2068867d85955ab4d8c737b2",
				},
				expectedStatus: http.StatusOK,
			},
			"bad headers": {
				expectedStatus: http.StatusForbidden,
			},
			"bad received mac": {
				headers: map[string]string{
					"User-Agent":      "GitHub-Hookshot/test",
					"X-GitHub-Event":  "push",
					"Content-Type":    "application/json",
					"X-Hub-Signature": "bad",
				},
				expectedStatus: http.StatusForbidden,
			},
			"hmac not equal": {
				headers: map[string]string{
					"User-Agent":      "GitHub-Hookshot/test",
					"X-GitHub-Event":  "push",
					"Content-Type":    "application/json",
					"X-Hub-Signature": "sha1=96ce",
				},
				expectedStatus: http.StatusForbidden,
			},
		}
	)

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				next = func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					assert.True(t, r.Context().Value(ctxEnsureGithubOriginKey).(bool))
				}
				r, _ = http.NewRequest("", "/", bytes.NewReader(bodyRaw))
				w    = httptest.NewRecorder()
			)

			for key, value := range test.headers {
				r.Header.Set(key, value)
			}

			EnsureGithubOrigin(secret)(http.HandlerFunc(next)).ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
