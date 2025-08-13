package types

import "time"

type CommandResult struct {
	CommandID     string
	Status        string
	ExitCode      int32
	StandardOut   string
	StandardErr   string
	ExecutionTime time.Duration
}
