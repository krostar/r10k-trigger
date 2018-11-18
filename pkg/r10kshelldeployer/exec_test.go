package r10kshelldeployer

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecCommand(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	var tests = map[string]struct {
		status          int
		output          string
		expectedFailure bool
	}{
		"success": {
			status:          0,
			output:          "",
			expectedFailure: false,
		},
		"failure": {
			status:          2,
			output:          "hello world!",
			expectedFailure: true,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				args = append([]string{"-test.run=TestHelperProcess", "--"}, strconv.Itoa(test.status), test.output)
				env  = []string{"GO_WANT_HELPER_PROCESS=1"}
			)
			status, output, err := execCommand(context.Background(), os.Args[0], args, env)
			if test.expectedFailure {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.status, status)
			assert.Equal(t, test.output, output)
		})
	}
}

func TestHelperProcess(*testing.T) {
	// expected command should look like
	//   go test -test.run=TestHelperProcess -- "<status>" "<output>"
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// filter original go test args until we find `--`
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}

	if len(args) != 2 {
		os.Exit(-1)
	}

	// write wanted output and exit with desired status
	os.Stderr.WriteString(args[1]) // nolint: errcheck, gosec
	n, err := strconv.Atoi(args[0])
	if err != nil {
		os.Exit(-2)
	}
	os.Exit(n)
}
