package r10kshelldeployer

// Shell implements the R10K deployer interface with a local shell instance.
type Shell struct {
	cfg  Config
	exec execFunc
}

// EnvironmentArgToReplace defines the arg key that will be replaces during the deploy command.
const EnvironmentArgToReplace = "<environment>"

// New creates a new r10k deployer using a local shell.
func New(opts ...Option) *Shell {
	var s = &Shell{
		cfg: Config{
			Command: "r10k",
			Args: []string{
				"deploy", "environment", EnvironmentArgToReplace, "--puppetfile", "--verbose",
			},
		},
		exec: execCommand,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
