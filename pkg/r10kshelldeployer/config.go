package r10kshelldeployer

import "errors"

// Config stores the shell deployer configuration.
type Config struct {
	Command     string   `json:"command"     yaml:"command"`
	Args        []string `json:"args"        yaml:"args"`
	Environment []string `json:"environment" yaml:"environment"`
}

// SetDefault sets sane default for monitor's config.
func (c *Config) SetDefault() {
	c.Command = "r10k"
	c.Args = []string{
		"deploy", "environment", EnvironmentArgToReplace, "--puppetfile", "--verbose",
	}
}

// Validate makes sure config has valid values.
func (c *Config) Validate() error {
	if c.Command == "" {
		return errors.New("command should not be empty")
	}
	return nil
}
