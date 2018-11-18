package middleware

import (
	"bytes"
	"crypto/hmac"
	"hash"
	"net/http"
	"testing"

	// this is only used for signature decryption, not for any sensitive content encryption
	// nolint: gosec
	"crypto/sha1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureGithubOrigin_checkGithubHeaders(t *testing.T) {
	var tests = map[string]struct {
		headers         map[string]string
		expectedFailure bool
	}{
		"all set": {
			headers: map[string]string{
				"User-Agent":     "GitHub-Hookshot/toto",
				"X-GitHub-Event": "push",
				"Content-Type":   "application/json",
			}, expectedFailure: false,
		},
		"missing content type": {
			headers: map[string]string{
				"User-Agent":     "GitHub-Hookshot/toto",
				"X-GitHub-Event": "push",
				"Content-Type":   "not json",
			}, expectedFailure: true,
		},
		"missing event type": {
			headers: map[string]string{
				"User-Agent":     "GitHub-Hookshot/toto",
				"X-GitHub-Event": "not push",
				"Content-Type":   "application/json",
			}, expectedFailure: true,
		},
		"wrong user agent": {
			headers: map[string]string{
				"User-Agent":     "not github",
				"X-GitHub-Event": "push",
				"Content-Type":   "application/json",
			}, expectedFailure: true,
		},
	}

	for name, test := range tests {
		// since name and test variables will keep the same addresses across the tests
		//   this shadow declaration is needed
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var r, _ = http.NewRequest("", "/", nil)
			for key, value := range test.headers {
				r.Header.Set(key, value)
			}

			err := checkGithubHeaders(r)

			if test.expectedFailure {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnsureGithubOrigin_getGithubSideGeneratedMAC(t *testing.T) {
	var tests = map[string]struct {
		header           string
		expectedFailure  bool
		expectedMAC      []byte
		expectedMACLen   int
		expectedHashFunc func() hash.Hash
	}{
		"correctly set": {
			header:           "sha1=68656c6c6f",
			expectedFailure:  false,
			expectedMAC:      []byte("hello"),
			expectedMACLen:   20,
			expectedHashFunc: sha1.New,
		},
		"wrong syntax": {
			header:          "sha1,68656c6c6f",
			expectedFailure: true,
		},
		"unhandled alg": {
			header:          "md5=68656c6c6f",
			expectedFailure: true,
		},
		"unable to decode hex": {
			header:          "sha1=toto",
			expectedFailure: true,
		},
	}

	for name, test := range tests {
		// since name and test variables will keep the same addresses across the tests
		//   this shadow declaration is needed
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var r, _ = http.NewRequest("", "/", nil)
			r.Header.Set("X-Hub-Signature", test.header)

			mac, hashFunc, err := getGithubSideGeneratedMAC(r)

			if test.expectedFailure {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if test.expectedHashFunc != nil {
				assert.IsType(t, test.expectedHashFunc(), hashFunc())
			} else {
				assert.Nil(t, hashFunc)
			}

			if test.expectedMAC != nil {
				assert.Equal(t, test.expectedMAC, mac[:len(test.expectedMAC)]) // expectedMAC is not padded with 0
				assert.Len(t, mac, test.expectedMACLen)
			} else {
				assert.Nil(t, mac)
			}
		})
	}
}

func TestEnsureGithubOrigin_computeGithubRequestMAC(t *testing.T) {
	var (
		secret      = "42universalanswer"
		hashFunc    = sha1.New
		bodyRaw     = []byte(`{"hello": "world"}`)
		expectedMAC = []byte{
			150, 206, 12, 218, 245, 139, 18,
			139, 32, 104, 134, 125, 133, 149,
			90, 180, 216, 199, 55, 178,
		}
		r, _ = http.NewRequest("", "/", bytes.NewReader(bodyRaw))
	)

	computedMAC, err := computeGithubRequestMAC(r, hashFunc, secret)
	require.NoError(t, err)
	assert.True(t, hmac.Equal(expectedMAC, computedMAC), "generated and expected mac differs")
}
