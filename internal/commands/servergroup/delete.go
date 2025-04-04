package servergroup

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "servergroup delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a server group",
			"upctl server-group delete 8abc8009-4325-4b23-4321-b1232cd81231",
			"upctl server-group delete my-server-group",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingServerGroup
	completion.ServerGroup
}

// InitCommand implements Command.InitCommand
func (c *deleteCommand) InitCommand() {
	// Deprecating servergroup in favour of server-group
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"servergroup"})
}

func Delete(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting server group %v", uuid)
	exec.PushProgressStarted(msg)

	err := svc.DeleteServerGroup(exec.Context(), &request.DeleteServerGroupRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Deprecating servergroup in favour of server-group
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"servergroup"}, "server-group")

	return Delete(exec, arg)
}
