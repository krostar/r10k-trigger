package r10kshelldeployer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_SetDefault(t *testing.T) {
	var cfg Config
	cfg.SetDefault()

	assert.Equal(t, "r10k", cfg.Command)
}

func TestConfig_Validate(t *testing.T) {
	t.Run("default is valid", func(t *testing.T) {
		var cfg Config
		cfg.SetDefault()
		require.NoError(t, cfg.Validate())
	})

	t.Run("missing command", func(t *testing.T) {
		var cfg Config
		cfg.SetDefault()

		cfg.Command = ""
		require.Error(t, cfg.Validate())
	})
}
