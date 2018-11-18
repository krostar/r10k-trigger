package r10kshelldeployer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShell_DeployR10kEnv(t *testing.T) {
	var tests = map[string]struct {
		environment     string
		expectedArgs    []string
		execerError     error
		expectedFailure bool
	}{
		"success": {
			environment:  "toto",
			expectedArgs: []string{"arg1", "toto", "arg3"},
		},
		"exec failed": {
			environment:     "toto",
			expectedArgs:    []string{"arg1", "toto", "arg3"},
			execerError:     errors.New("eww"),
			expectedFailure: true,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			shell := New()
			shell.cfg.Args = []string{"arg1", EnvironmentArgToReplace, "arg3"}
			shell.exec = func(ctx context.Context, cmd string, args []string, env []string) (int, string, error) {
				assert.Equal(t, test.expectedArgs, args)
				return 0, "", test.execerError
			}

			err := shell.DeployR10KEnv(context.Background(), test.environment)
			if test.expectedFailure {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
