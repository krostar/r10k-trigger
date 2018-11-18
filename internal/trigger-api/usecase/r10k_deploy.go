package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// R10KEnvDeployer defines the contract of the r10k deployer.
//go:generate mockery -inpkg -testonly -name R10KEnvDeployer
type R10KEnvDeployer interface {
	DeployR10KEnv(ctx context.Context, env string) (err error)
}

// R10KDeploy returns a new usecase which deploies an r10k environment.
type R10KDeploy struct{ Deployer R10KEnvDeployer }

// DeployR10KEnv handles the trigger of the r10k deployment.
func (r *R10KDeploy) DeployR10KEnv(ctx context.Context, env string) error {
	ctx, span := trace.StartSpan(ctx, "usecase.r10k-deploy")
	defer span.End()

	var err error
	defer recordR10KDeployStat(ctx, time.Now(), env, &err)

	if err = r.Deployer.DeployR10KEnv(ctx, env); err != nil {
		err = errors.Wrapf(err, "unable to deploy r10k on environment %s", env)
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
		return err
	}

	span.SetStatus(trace.Status{Code: trace.StatusCodeOK})
	return nil
}

// DeployR10KEnvAsync calls DeployR10KEnvironment in a asynchronous way.
func (r *R10KDeploy) DeployR10KEnvAsync(ctx context.Context, env string,
	onDeploy func(env string), onError func(env string, err error),
) {
	_, span := trace.StartSpan(ctx, "usecase.r10k-deploy-async")
	go func() {
		defer span.End()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		ctx = trace.NewContext(ctx, span)

		if err := r.DeployR10KEnv(ctx, env); err != nil {
			err = errors.Wrap(err, "unable to deploy asynchronously")
			if onError != nil {
				onError(env, err)
			}
			span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
			return
		}
		if onDeploy != nil {
			onDeploy(env)
		}
		span.SetStatus(trace.Status{Code: trace.StatusCodeOK})
	}()
}
