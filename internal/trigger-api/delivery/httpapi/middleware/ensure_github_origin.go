package middleware

import (
	"context"
	"crypto/hmac"
	"encoding/hex"
	"hash"
	"io/ioutil"
	"net/http"
	"strings"

	// this is only used for signature decryption, not for any sensitive content encryption
	// nolint: gosec
	"crypto/sha1"

	"github.com/krostar/logger/logmid"
	"github.com/pkg/errors"
)

type ctxEnsureGithubOrigin string

// ctxEnsureGithubOriginKey is the context key that holds the fact that the origin has been ensured.
var ctxEnsureGithubOriginKey = ctxEnsureGithubOrigin("github-ensure-origin") // nolint: gochecknoglobals

// EnsureGithubOrigin is a middleware that ensure that the caller is github.
func EnsureGithubOrigin(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			if err := ensureGithubOrigin(r, secret); err != nil {
				logmid.AddErrorInContext(ctx, err)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ctx = context.WithValue(ctx, ctxEnsureGithubOriginKey, true)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ensureGithubOrigin(r *http.Request, secret string) error {
	if err := checkGithubHeaders(r); err != nil {
		return errors.Wrap(err, "header does not contain what we expect them to")
	}

	macReceived, hash, err := getGithubSideGeneratedMAC(r)
	if err != nil {
		return errors.Wrap(err, "MAC given by github cannot be retrieved")
	}

	macGenerated, err := computeGithubRequestMAC(r, hash, secret)
	if err != nil {
		return errors.Wrap(err, "unable to generate request MAC")
	}

	if !hmac.Equal(macGenerated, macReceived) {
		return errors.New("received and generated MAC are not the same")
	}

	return nil
}

func checkGithubHeaders(r *http.Request) error {
	// the user agent is not one of github's one
	if !strings.HasPrefix(r.Header.Get("User-Agent"), "GitHub-Hookshot/") {
		return errors.Errorf("wrong user-agent prefix: %s", r.Header.Get("User-Agent"))
	}

	// the event is not a push event
	if r.Header.Get("X-GitHub-Event") != "push" {
		return errors.Errorf("wrong event: %s", r.Header.Get("X-GitHub-Event"))
	}

	// the content type is not JSON
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.Errorf("wrong content type: %s", r.Header.Get("Content-Type"))
	}

	return nil
}

func getGithubSideGeneratedMAC(r *http.Request) ([]byte, func() hash.Hash, error) {
	// X-Hub-Signature should look like alg=sig
	signatureReceived := strings.SplitN(r.Header.Get("X-Hub-Signature"), "=", 2)
	if len(signatureReceived) != 2 {
		return nil, nil, errors.Errorf("signature should look like alg=sig")
	}

	// get the hash function from the the signature
	//   and the expected length of the output of the hash function
	var hashFunc func() hash.Hash
	var hashLength int
	switch alg := signatureReceived[0]; alg {
	case "sha1":
		hashFunc = sha1.New
		hashLength = 20
	default:
		return nil, nil, errors.Errorf("unhandled alg %s", alg)
	}

	var macReceived = make([]byte, hashLength)
	if _, err := hex.Decode(macReceived, []byte(signatureReceived[1])); err != nil {
		return nil, nil, errors.Wrap(err, "unable to hex decode received mac")
	}

	return macReceived, hashFunc, nil
}

func computeGithubRequestMAC(r *http.Request, hashFunc func() hash.Hash, secret string) ([]byte, error) {
	// get a copy of the body to let handler read it later
	bodyReader, err := r.GetBody()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get a copy of the body")
	}
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "unable to ready body")
	}

	// compute the mac of the body
	hasher := hmac.New(hashFunc, []byte(secret))
	if _, err := hasher.Write(body); err != nil {
		return nil, errors.Wrap(err, "unable to write to hmac")
	}
	macGenerated := hasher.Sum(nil)

	return macGenerated, nil
}
