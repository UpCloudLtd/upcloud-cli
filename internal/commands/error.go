package commands

import (
	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
)

// HandleError updates error details to progress log message identified by given key. Returns (nil, err), where err is the err passed in as input.
func HandleError(exec Executor, key string, err error) (output.Output, error) {
	exec.PushProgressUpdate(messages.Update{
		Key:     key,
		Status:  messages.MessageStatusError,
		Details: "Error: " + err.Error(),
	})

	return nil, err
}
