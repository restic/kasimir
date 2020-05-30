package main

import (
	"github.com/spf13/cobra"
)

func addCommandHooks(root *cobra.Command, gopts *GlobalOptions, cfg *Config) {
	var cmd = &cobra.Command{
		Use:           "hooks",
		Short:         "run hooks and print the result",
		Long:          "run all hooks and print the current result for each",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := CheckConfig{
				Dir:     ".",
				Version: gopts.Version,
			}

			return RunHooks(cfg)
		},
	}

	root.AddCommand(cmd)
}
