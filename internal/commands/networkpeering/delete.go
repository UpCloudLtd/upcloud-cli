package networkpeering

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "networkpeering delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a network peering",
			"upctl network-peering delete 8abc8009-4325-4b23-4321-b1232cd81231",
			"upctl network-peering delete my-network-peering",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingNetworkPeering
	completion.NetworkPeering
}

// InitCommand implements Command.InitCommand
func (c *deleteCommand) InitCommand() {
	// Deprecating networkpeering in favour of network-peering
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"networkpeering"})
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Deprecating networkpeering in favour of network-peering
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"networkpeering"}, "network-peering")

	svc := exec.All()
	msg := fmt.Sprintf("Deleting network peering %v", arg)
	exec.PushProgressStarted(msg)

	err := svc.DeleteNetworkPeering(exec.Context(), &request.DeleteNetworkPeeringRequest{
		UUID: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
