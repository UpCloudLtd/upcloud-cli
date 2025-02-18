package token

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "token delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete an API token",
			"upctl account token delete 0c0e2abf-cd89-490b-abdb-d06db6e8d816",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingToken
	completion.Token
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	svc := exec.Token()
	msg := fmt.Sprintf("Deleting API token %v", arg)
	exec.PushProgressStarted(msg)

	err := svc.DeleteToken(exec.Context(), &request.DeleteTokenRequest{
		ID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
