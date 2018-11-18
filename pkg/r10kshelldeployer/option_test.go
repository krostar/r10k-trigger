package r10kshelldeployer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithConfig(t *testing.T) {
	var (
		defaultConfig = Config{
			Command: "r10k",
			Args: []string{
				"deploy", "environment", EnvironmentArgToReplace, "--puppetfile", "--verbose",
			},
		}
		tests = map[string]struct {
			cfg             *Config
			expectedCommand string
			expectedArgs    []string
			expectedEnv     []string
		}{
			"nil config": {
				cfg:             nil,
				expectedCommand: defaultConfig.Command,
				expectedArgs:    defaultConfig.Args,
				expectedEnv:     defaultConfig.Environment,
			},
			"different command": {
				cfg:             &Config{Command: "different"},
				expectedCommand: "different",
				expectedArgs:    nil,
				expectedEnv:     defaultConfig.Environment,
			},
			"different args": {
				cfg:             &Config{Args: []string{"bli", "blu"}},
				expectedCommand: defaultConfig.Command,
				expectedArgs:    []string{"bli", "blu"},
				expectedEnv:     defaultConfig.Environment,
			},
			"different env": {
				cfg:             &Config{Environment: []string{"blum"}},
				expectedCommand: defaultConfig.Command,
				expectedArgs:    defaultConfig.Args,
				expectedEnv:     append(defaultConfig.Environment, "blum"),
			},
		}
	)

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var shell = New(WithConfig(test.cfg))

			assert.Equal(t, test.expectedCommand, shell.cfg.Command)
			assert.Equal(t, test.expectedArgs, shell.cfg.Args)
			assert.ElementsMatch(t, test.expectedEnv, shell.cfg.Environment)
		})
	}
}
