package gateway

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "gateway delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a gateway",
			"upctl gateway delete 8abc8009-4325-4b23-4321-b1232cd81231",
			"upctl gateway delete my-gateway",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingGateway
	completion.Gateway
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting gateway %v", arg)
	exec.PushProgressStarted(msg)

	err := svc.DeleteGateway(exec.Context(), &request.DeleteGatewayRequest{
		UUID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
