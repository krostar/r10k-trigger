package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/krostar/r10k-trigger/internal/pkg/app"
)

type flags struct {
	help       bool
	version    bool
	configFile string
}

func initFlags(args []string) *flags {
	var f flags

	flagset := flag.NewFlagSet(app.Name(), flag.ContinueOnError)
	flagset.Usage = flagUsage(flagset)

	flagset.BoolVar(&f.help, "help", false, "displays this help and quit")
	flagset.BoolVar(&f.version, "version", false, `absolute path to the configuration file to use`)
	flagset.StringVar(&f.configFile, "config", "", "configuration file to use")

	if err := flagset.Parse(args); err != nil {
		panic(err)
	}

	if f.help {
		flagset.Usage()
		os.Exit(0)
	}

	if f.version {
		fmt.Printf("%s version %s\n", app.Name(), app.Version())
		os.Exit(0)
	}

	if f.configFile == "" {
		fmt.Printf("Usage error: -config flag is required and can't be empty\n\n")
		flagset.Usage()
		os.Exit(2)
	}

	return &f
}

func flagUsage(flagset *flag.FlagSet) func() {
	return func() {
		var (
			flags bytes.Buffer
			w     = tabwriter.NewWriter(&flags, 10, 1, 5, ' ', 0)
		)

		flagset.VisitAll(func(f *flag.Flag) {
			w.Write([]byte(fmt.Sprintf("   -%s\t%s\n", f.Name, f.Usage))) // nolint: errcheck, gosec
		})
		w.Flush() // nolint: errcheck, gosec

		fmt.Fprintf(os.Stderr, `%s exposes a hook to trigger r10k from an http call.

Find more information at: https://github.com/krostar/r10k-trigger.

Usage
   %s -config FILE

Example
   %s -version
   %s -config /etc/r10k/trigger.yml

Flags
%s
`, app.Name(), app.Name(), app.Name(), app.Name(), flags.String())
	}
}
