package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Options collects all configuration for a single run.
type Options struct {
	Verbose             bool
	Debug               bool
	Config              string
	ConfigFileSpecified bool
}

func main() {
	var opts Options

	var cmd = &cobra.Command{
		Short:         "pondi [options]",
		Long:          "pondi builds release assets (binaries, source code archive) and creates a new release on GitHub ",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("config") {
				// signal that it is an error if the config file does not exist
				opts.ConfigFileSpecified = true
			}

			return run(opts, args)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.Verbose, "verbose", false, "be verbose")
	flags.BoolVar(&opts.Debug, "debug", false, "print debug messages")
	flags.StringVar(&opts.Config, "config", ".pondi.yml", "load configuration from `file`")

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "exiting, error: %v\n", err)
	}
}

func run(opts Options, args []string) error {
	var cfg Config
	var err error

	// try to load the config file, fall back to the default config if the file
	// does not exist and the file name has not been set manually
	cfg, err = LoadConfig(opts.Config)
	if errors.Is(err, os.ErrNotExist) && !opts.ConfigFileSpecified {
		cfg = DefaultConfig
		err = nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("main, config: %#v\n", cfg)
	return nil
}
