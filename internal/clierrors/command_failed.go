package clierrors

import "fmt"

type CommandFailedError struct {
	FailedCount int
}

var _ ClientError = CommandFailedError{}

func (err CommandFailedError) ErrorCode() int {
	return min(err.FailedCount, 99)
}

func (err CommandFailedError) Error() string {
	return fmt.Sprintf("Command execution failed for %d resource(s)", err.FailedCount)
}
