package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

// GlobalOptions collects all configuration for a single run.
type GlobalOptions struct {
	Verbose       bool
	Debug         bool
	Config        string
	DisableChecks []string

	Version string
}

func main() {
	var (
		gopts GlobalOptions
		cfg   Config
	)

	var cmd = &cobra.Command{
		Use:           "pondi [flags]",
		Short:         "pondi",
		Long:          "pondi builds release assets (binaries, source code archive) and creates a new release on GitHub ",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if gopts.Version != "" {
				matched, err := regexp.MatchString(`^\d+\.\d+\.\d+$`, gopts.Version)
				if err != nil {
					panic(err)
				}

				if !matched {
					return fmt.Errorf("version %q is invalid (format: 1.2.3)", gopts.Version)
				}
			}

			cfg, err = LoadConfig(gopts.Config)

			// if the config file was not explicitly passed and it does not
			// exist, just use the default config.
			if !cmd.Flags().Changed("config") && errors.Is(err, os.ErrNotExist) {
				cfg = DefaultConfig
				err = nil
			}

			if err != nil {
				return err
			}

			return nil
		},
	}

	addCommandCheck(cmd, &gopts, &cfg)
	addCommandHooks(cmd, &gopts, &cfg)

	flags := cmd.PersistentFlags()
	flags.BoolVar(&gopts.Verbose, "verbose", false, "be verbose")
	flags.BoolVar(&gopts.Debug, "debug", false, "print debug messages")
	flags.StringVar(&gopts.Config, "config", ".pondi.yml", "load configuration from `file`")
	flags.StringSliceVar(&gopts.DisableChecks, "disable-checks", nil, "disable checks `name1,name2,[...]`")

	flags.StringVar(&gopts.Version, "version", "", "release version (format: `1.2.3`)")

	err := cmd.MarkPersistentFlagRequired("version")
	if err != nil {
		panic(err)
	}

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
