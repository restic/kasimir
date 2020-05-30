package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Hook is a command which is run before releasing.
type Hook struct {
	Name        string
	Description string
	Command     []string
}

var AllHooks = []Hook{
	{
		Name:        "go-mod-download",
		Description: "run 'go mod download' to make sure all Go modules are accessible",
		Command:     []string{"go", "mod", "download"},
	},
	{
		Name:        "go-generate",
		Description: "run 'go generate ./...' to make sure all generated code is up to date",
		Command:     []string{"go", "generate", "./..."},
	},
	{
		Name:        "gofmt",
		Description: "run 'gofmt -w .' to format all source code",
		Command:     []string{"gofmt", "-w", "."},
	},
}

// RunHooks run all hooks.
func RunHooks(cfg CheckConfig) error {
	// check for uncommitted changes before running hooks
	err := CheckUncommittedChanges(cfg)
	if err != nil {
		return err
	}

	for _, hook := range AllHooks {
		fmt.Printf("run %v\n", hook.Name)
		cmd := exec.Command(hook.Command[0], hook.Command[1:]...)
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			return fmt.Errorf("hook %v failed: %w", hook.Name, err)
		}
	}

	// afterwards, check if the repository contains uncommitted changes
	return CheckUncommittedChanges(cfg)
}
