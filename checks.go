package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Check defines a single check.
type Check struct {
	Name        string
	Description string
	Run         func(CheckConfig) error
}

// CheckConfig contains information needed to perform checks.
type CheckConfig struct {
	Dir     string
	Version string
}

// CheckResult bundles a check with its result after running.
type CheckResult struct {
	Check
	Result error
}

// AllChecks contains a list of all checks with descriptions.
var AllChecks = []Check{
	{
		Name:        "branch-master",
		Description: "the current branch is master",
		Run:         CheckBranchMaster,
	},
	{
		Name:        "uncommitted-changes",
		Description: "no uncommitted changes or files exist",
		Run:         CheckUncommittedChanges,
	},
	{
		Name:        "gofmt",
		Description: "code is formatted with 'gofmt'",
		Run:         CheckGofmt,
	},
	{
		Name:        "tag-exists",
		Description: "version tag does not exist",
		Run:         CheckTagExists,
	},
}

// FilterChecks returns a list of checks without the ones listed in reject. For
// invalid names, an error is returned.
func FilterChecks(list []Check, reject []string) (result []Check, err error) {
	all := make(map[string]struct{})
	for _, check := range list {
		all[check.Name] = struct{}{}
	}

	disabled := make(map[string]struct{})

	for _, name := range reject {
		if _, ok := all[name]; !ok {
			return nil, fmt.Errorf("invalid check name %q", name)
		}

		disabled[name] = struct{}{}
	}

	for _, check := range list {
		if _, ok := disabled[check.Name]; ok {
			continue
		}

		result = append(result, check)
	}

	return result, nil
}

// RunChecks runs all checks.
func RunChecks(cfg CheckConfig, checks []Check) (result []CheckResult, err error) {
	merr := &MultiError{}

	for _, check := range checks {
		err := check.Run(cfg)
		if err != nil {
			merr.Insert(fmt.Errorf("%v failed: %w", check.Name, err))
		}

		result = append(result, CheckResult{
			Check:  check,
			Result: err,
		})
	}

	if merr.Length() == 0 {
		return result, nil
	}

	return result, merr
}

func CheckBranchMaster(cfg CheckConfig) error {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = cfg.Dir

	name, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to find current branch: %w", err)
	}

	branch := strings.TrimRight(string(name), "\n")
	if branch != "master" {
		return fmt.Errorf("current branch is %q instead of master", branch)
	}

	return nil
}

func CheckUncommittedChanges(cfg CheckConfig) error {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = cfg.Dir

	buf, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to check for uncommitted changes: %w", err)
	}

	status := strings.TrimRight(string(buf), "\n")
	if len(status) > 0 {
		return fmt.Errorf("repository contains uncommitted changes or additional files")
	}

	return nil
}

func CheckGofmt(cfg CheckConfig) error {
	cmd := exec.Command("gofmt", "-l", ".")
	cmd.Dir = cfg.Dir

	buf, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("running 'gofmt' failed: %w", err)
	}

	text := strings.TrimRight(string(buf), "\n")
	text = strings.ReplaceAll(text, "\n", ", ")

	if len(text) > 0 {
		return fmt.Errorf("repository contains files not formatted with gofmt: %v", text)
	}

	return nil
}

func CheckTagExists(cfg CheckConfig) error {
	cmd := exec.Command("git", "tag", "-l", "v"+cfg.Version)
	cmd.Dir = cfg.Dir

	buf, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("running 'git tag -l %v' failed: %w", cfg.Version, err)
	}

	tag := strings.TrimRight(string(buf), "\n")
	if tag != "" {
		return fmt.Errorf("tag %q already exists", cfg.Version)
	}

	return nil
}
