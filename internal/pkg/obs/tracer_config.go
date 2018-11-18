package obs

import "errors"

// TracerConfig stores tracer configuration.
type TracerConfig struct {
	Enabled   bool   `json:"enabled"    yaml:"enabled"`
	ZipkinURL string `json:"zipkin-url" yaml:"zipkin-url"`
}

// SetDefault sets sane default for tracer's config.
func (c *TracerConfig) SetDefault() {
	c.Enabled = false
}

// Validate makes sure config has valid values.
func (c *TracerConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	if c.ZipkinURL == "" {
		return errors.New("zipkin url should not be empty")
	}

	return nil
}
