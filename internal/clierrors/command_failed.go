package clierrors

import "fmt"

type CommandFailedError struct {
	FailedCount int
}

func (err CommandFailedError) Error() string {
	return fmt.Sprintf("Command execution failed for %d resource(s)", err.FailedCount)
}
