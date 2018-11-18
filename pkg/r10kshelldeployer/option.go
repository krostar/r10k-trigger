package r10kshelldeployer

// Option defines a function prototype to apply options to the Shell instance.
type Option func(*Shell)

// WithConfig sets all the shell configurations.
func WithConfig(cfg *Config) Option {
	return func(s *Shell) {
		if cfg == nil {
			return
		}

		if cfg.Command != "" {
			s.cfg.Command = cfg.Command
			s.cfg.Args = nil
		}

		if cfg.Args != nil {
			s.cfg.Args = cfg.Args
		}

		s.cfg.Environment = append(s.cfg.Environment, cfg.Environment...)
	}
}
