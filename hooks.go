package main

// Hook is a command which is run before releasing.
type Hook struct {
	Name    string
	Command []string
}
