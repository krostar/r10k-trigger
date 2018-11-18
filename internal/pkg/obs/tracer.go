package obs

import (
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"

	"github.com/krostar/r10k-trigger/internal/pkg/app"
)

func initTracer(cfg TracerConfig) (func(), error) {
	var stopFunc = func() {}

	if !cfg.Enabled {
		return stopFunc, nil
	}

	// by default all traces are dropped as the sampling decision
	// should be cautiously set to avoid a useless number of traces
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.NeverSample()})
	zipkinEndpoint, err := openzipkin.NewEndpoint(app.Name(), "")
	if err != nil {
		return stopFunc, errors.Wrap(err, "unable to create zipkin endpoint")
	}

	zipkinReporter := zipkinHTTP.NewReporter(cfg.ZipkinURL)
	stopFunc = func() {
		zipkinReporter.Close() // nolint: errcheck, gosec
	}

	zipkinExporter := zipkin.NewExporter(zipkinReporter, zipkinEndpoint)
	trace.RegisterExporter(zipkinExporter)

	return stopFunc, nil
}
