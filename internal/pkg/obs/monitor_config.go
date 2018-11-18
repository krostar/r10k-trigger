package obs

import (
	"errors"
	"strings"
	"time"
)

// MonitorConfig stores the configuration for the monitor.
type MonitorConfig struct {
	Enabled               bool          `json:"enabled"                 yaml:"enabled"`
	PrometheusEndpoint    string        `json:"prometheus-endpoint"     yaml:"prometheus-endpoint"`
	ProcessReportInterval time.Duration `json:"process-report-interval" yaml:"process-report-interval"`
}

// SetDefault sets sane default for monitor's config.
func (c *MonitorConfig) SetDefault() {
	c.Enabled = false
	c.PrometheusEndpoint = "/metrics"
	c.ProcessReportInterval = 20 * time.Second
}

// Validate makes sure config has valid values.
func (c *MonitorConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	if c.PrometheusEndpoint == "" {
		return errors.New("prometheus endpoint should not be empty")
	}

	if !strings.HasPrefix(c.PrometheusEndpoint, "/") {
		return errors.New("prometheus endpoint should start with /")
	}

	if c.ProcessReportInterval < time.Second {
		return errors.New("process report interval should be higher than a second")
	}

	return nil
}
