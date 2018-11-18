package r10kshelldeployer

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// DeployR10KEnv call the deploy command for r10k via a local shell call.
func (sh *Shell) DeployR10KEnv(ctx context.Context, env string) error {
	ctx, span := trace.StartSpan(ctx, "r10kshelldeployer.deploy-r10k-env")
	defer span.End()

	// replace EnvironmentArgToReplace by the actual env
	var args []string
	for _, arg := range sh.cfg.Args {
		if arg == EnvironmentArgToReplace {
			arg = env
		}
		args = append(args, arg)
	}

	_, stderr, err := sh.exec(ctx, sh.cfg.Command, args, sh.cfg.Environment)
	if err != nil {
		span.SetStatus(trace.Status{Code: trace.StatusCodeInternal, Message: err.Error()})
		return errors.Wrapf(err, "unable to execute shell command: %q", stderr)
	}

	span.SetStatus(trace.Status{Code: trace.StatusCodeOK})
	return nil
}
