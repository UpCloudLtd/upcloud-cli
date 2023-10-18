package commands

import (
	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

// HandleError updates error details to progress log message identified by given key. Returns (nil, err), where err is the err passed in as input.
func HandleError(exec Executor, key string, err error) (output.Output, error) {
	exec.PushProgressUpdate(messages.Update{
		Key:     key,
		Status:  messages.MessageStatusError,
		Details: "Error: " + err.Error(),
	})

	return nil, handledError{err}
}

type handledError struct {
	err error
}

func (h handledError) Error() string {
	return h.err.Error()
}

// outputError outputs given error to progress log, if the error has not been already handled by HandleError
func outputError(arg string, err error, exec Executor) {
	if err == nil {
		return
	}

	if _, ok := err.(handledError); ok {
		return
	}

	msg := "Command execution failed"
	if arg != "" {
		msg += " for " + arg
	}

	exec.PushProgressUpdate(messages.Update{
		Message: msg,
		Status:  messages.MessageStatusError,
		Details: "Error: " + err.Error(),
	})
}
