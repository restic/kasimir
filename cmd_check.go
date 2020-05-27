package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func addCommandCheck(root *cobra.Command, gopts *GlobalOptions, cfg *Config) {
	var cmd = &cobra.Command{
		Use:           "check",
		Short:         "run checks and print the result",
		Long:          "run all checks and print the current result for each",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck(*gopts, *cfg, args)
		},
	}

	root.AddCommand(cmd)
}

func runCheck(gopts GlobalOptions, _ Config, _ []string) error {
	checks, err := FilterChecks(AllChecks, gopts.DisableChecks)
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(os.Stdout, 3, 8, 2, ' ', 0)

	results, err := RunChecks(checks)

	for _, result := range results {
		text := ""
		status := "✓"

		if result.Result != nil {
			text = result.Result.Error()
			status = "✗"
		}

		fmt.Fprintf(tw, "%s\t%v\t", status, result.Check.Name)

		if gopts.Verbose {
			fmt.Fprintf(tw, "%v\t", result.Check.Description)
		}

		fmt.Fprintln(tw, text)
	}

	tw.Flush()
	fmt.Println()

	return err
}
