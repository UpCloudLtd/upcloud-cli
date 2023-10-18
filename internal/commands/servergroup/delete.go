package servergroup

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// DeleteCommand creates the "servergroup delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a server group",
			"upctl servergroup delete 8abc8009-4325-4b23-4321-b1232cd81231",
			"upctl servergroup delete my-server-group",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingServerGroup
	completion.ServerGroup
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting server group %v", arg)
	exec.PushProgressStarted(msg)

	err := svc.DeleteServerGroup(exec.Context(), &request.DeleteServerGroupRequest{
		UUID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
