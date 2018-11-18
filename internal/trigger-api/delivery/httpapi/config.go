package httpapi

import (
	"time"

	"github.com/pkg/errors"
)

// Config defines the http configuration.
type Config struct {
	ListenAddress        string        `json:"listen-address"         yaml:"listen-address"`
	RequestTimeout       time.Duration `json:"request-timeout"        yaml:"request-timeout"`
	TLS                  *TLSConfig    `json:"tls"                    yaml:"tls"`
	EnsureDeployerOrigin bool          `json:"ensure-deployer-origin" yaml:"ensure-deployer-origin"`
	DeployerMACSecret    string        `json:"deployer-mac-secret"    yaml:"deployer-mac-secret"`
}

// SetDefault sets sane default for monitor's config.
func (c *Config) SetDefault() {
	c.ListenAddress = ":8080"
	c.RequestTimeout = 3 * time.Second
}

// Validate makes sure config has valid values.
func (c *Config) Validate() error {
	if c.ListenAddress == "" {
		return errors.New("listen address can't be empty")
	}
	if c.RequestTimeout < 500*time.Millisecond {
		return errors.New("request timeout should be higher than 500ms")
	}
	if c.EnsureDeployerOrigin {
		if c.DeployerMACSecret == "" {
			return errors.New("since deployer origin will be ensured, deployer secret can't be empty")
		}
	}
	return nil
}

// TLSConfig defines the tls configuration.
type TLSConfig struct {
	CertFile string `json:"cert-file" yaml:"cert-file"`
	KeyFile  string `json:"key-file"  yaml:"key-file"`
}

// Validate makes sure config has valid values.
func (c *TLSConfig) Validate() error {
	if c.CertFile == "" {
		return errors.New("cert file can't be empty")
	}
	if c.KeyFile == "" {
		return errors.New("cert file can't be empty")
	}
	return nil
}
