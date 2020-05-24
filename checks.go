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
	Run         func() error
}

// CheckResult bundles a check with its result after running.
type CheckResult struct {
	Check
	Result error
}

// AllChecks contains a list of all checks with descriptions.
var AllChecks = []Check{
	{
		Name:        "check-branch-master",
		Description: "test if the current branch is master",
		Run:         CheckBranchMaster,
	},
	{
		Name:        "check-uncommitted-changes",
		Description: "test if uncommitted changes or files exist",
		Run:         CheckUncommittedChanges,
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
func RunChecks(checks []Check) (result []CheckResult, err error) {
	merr := &MultiError{}

	for _, check := range checks {
		err := check.Run()
		if err != nil {
			merr.Insert(fmt.Errorf("check %v failed: %w", check.Name, err))
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

func CheckBranchMaster() error {
	name, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		return fmt.Errorf("unable to find current branch: %w", err)
	}

	branch := strings.TrimRight(string(name), "\n")
	if branch != "master" {
		return fmt.Errorf("current branch is %q instead of master", branch)
	}

	return nil
}

func CheckUncommittedChanges() error {
	buf, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return fmt.Errorf("unable to check for uncommitted changes: %w", err)
	}

	status := strings.TrimRight(string(buf), "\n")
	if len(status) > 0 {
		return fmt.Errorf("repository contains uncommitted changes or additional files")
	}

	return nil
}
