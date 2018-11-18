package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// nolint: gochecknoglobals
var (
	r10kDeployStatsTagKeyEnvironment, _ = tag.NewKey("environment")
	r10kDeployStatsTagKeyStatus, _      = tag.NewKey("status")

	r10kDeployStatsMetricDeployCount = stats.Int64("r10k/deploy/count",
		"Number of r10k deployment", stats.UnitDimensionless,
	)
	r10kDeployStatsMetricDeployDuration = stats.Int64("r10k/deploy/duration",
		"Duration of the r10k deployment", stats.UnitMilliseconds,
	)

	r10kDeployStatsDistributionDuration = view.Distribution(
		1, 300, 500, 800, 1000, 1250, 1500, 1750, 2000, 2500, 3000, 3500, 5000, 7000, 10000,
	)
)

// nolint: gochecknoinits
func init() {
	if err := view.Register(
		&view.View{
			Name:        r10kDeployStatsMetricDeployCount.Name(),
			Description: r10kDeployStatsMetricDeployCount.Description(),
			Measure:     r10kDeployStatsMetricDeployCount,
			Aggregation: view.Count(),
			TagKeys: []tag.Key{
				r10kDeployStatsTagKeyEnvironment,
				r10kDeployStatsTagKeyStatus,
			},
		},
		&view.View{
			Name:        r10kDeployStatsMetricDeployDuration.Name(),
			Description: r10kDeployStatsMetricDeployDuration.Description(),
			Measure:     r10kDeployStatsMetricDeployDuration,
			Aggregation: r10kDeployStatsDistributionDuration,
			TagKeys: []tag.Key{
				r10kDeployStatsTagKeyEnvironment,
				r10kDeployStatsTagKeyStatus,
			},
		},
	); err != nil {
		panic(errors.Wrap(err, "unable to register deploy stats view"))
	}
}

func recordR10KDeployStat(ctx context.Context, startedAt time.Time, environment string, err *error) {
	var (
		duration = time.Since(startedAt).Nanoseconds() / int64(time.Millisecond)
		status   = "succeed"

		mesurements = []stats.Measurement{
			r10kDeployStatsMetricDeployCount.M(1),
		}
	)

	if err != nil && *err != nil {
		status = "failed"
	} else {
		mesurements = append(mesurements, r10kDeployStatsMetricDeployDuration.M(duration))
	}

	ctx, _ = tag.New(ctx, // nolint: errcheck
		tag.Upsert(r10kDeployStatsTagKeyEnvironment, environment),
		tag.Upsert(r10kDeployStatsTagKeyStatus, status),
	)

	stats.Record(ctx, mesurements...)
}
