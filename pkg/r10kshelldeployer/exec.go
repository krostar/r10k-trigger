package r10kshelldeployer

import (
	"context"
	"os/exec"
	"syscall"
)

type execFunc func(ctx context.Context, command string, args []string, env []string) (int, string, error)

func execCommand(ctx context.Context, command string, args []string, env []string) (int, string, error) {
	// The nolint below is due to the fact the command is a variable
	// but it's assumed as it's chose server-side.
	// nolint: gosec
	var (
		cmd    = exec.CommandContext(ctx, command, args...)
		status = -1
	)
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if wstatus, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				status = wstatus.ExitStatus()
			}
		}
	} else {
		status = 0
	}
	return status, string(out), err
}
