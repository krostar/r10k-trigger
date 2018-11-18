package httpapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Defaults(t *testing.T) {
	var cfg Config
	cfg.SetDefault()

	assert.NotEmpty(t, cfg.ListenAddress)
	assert.True(t, cfg.RequestTimeout > 500*time.Millisecond)
}

func TestConfig_Validate(t *testing.T) {
	t.Run("default is not enough as a secret is required", func(t *testing.T) {
		t.Parallel()

		var cfg Config
		cfg.SetDefault()

		cfg.DeployerMACSecret = "blbilbi"
		require.NoError(t, cfg.Validate())
	})
	t.Run("empty listen address", func(t *testing.T) {
		t.Parallel()

		var cfg Config
		cfg.SetDefault()
		cfg.DeployerMACSecret = "blbilbi"

		cfg.ListenAddress = ""

		require.Error(t, cfg.Validate())
	})
	t.Run("timeout is not too small", func(t *testing.T) {
		t.Parallel()

		var cfg Config
		cfg.SetDefault()

		cfg.RequestTimeout = 0

		require.Error(t, cfg.Validate())
	})
}

func TestTLSConfig_Validate(t *testing.T) {
	var cfg TLSConfig

	assert.Error(t, cfg.Validate())

	cfg.CertFile = "/bli"
	assert.Error(t, cfg.Validate())

	cfg.KeyFile = "/bla"
	assert.NoError(t, cfg.Validate())
}
