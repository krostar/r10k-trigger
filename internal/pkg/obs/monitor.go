package obs

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
)

func initMonitor(cfg MonitorConfig, namespace string, writerErr io.Writer) (string, http.HandlerFunc, func(), error) {
	if !cfg.Enabled {
		return "", nil, nil, nil
	}

	prometheusExporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: namespace,
		OnError: func(err error) {
			if writerErr != nil {
				io.WriteString(writerErr, err.Error()) // nolint: errcheck, gosec
			}
		},
	})
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "unable to create prometheus exporter")
	}
	view.RegisterExporter(prometheusExporter)

	stopMonitorReport := MonitorProcess(cfg.ProcessReportInterval)

	return cfg.PrometheusEndpoint, prometheusExporter.ServeHTTP, func() {
		stopMonitorReport()
		view.UnregisterExporter(prometheusExporter)
	}, nil
}
