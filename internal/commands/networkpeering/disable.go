package networkpeering

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DisableCommand creates the "networkpeering disable" command
func DisableCommand() commands.Command {
	return &disableCommand{
		BaseCommand: commands.New(
			"disable",
			"Disable a network peering",
			"upctl networkpeering disable 8abc8009-4325-4b23-4321-b1232cd81231",
			"upctl networkpeering disable my-network-peering",
		),
	}
}

type disableCommand struct {
	*commands.BaseCommand
	resolver.CachingNetworkPeering
	completion.NetworkPeering
}

// Execute implements commands.MultipleArgumentCommand
func (c *disableCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Disabling network peering %v", arg)
	exec.PushProgressStarted(msg)

	peering, err := svc.ModifyNetworkPeering(exec.Context(), &request.ModifyNetworkPeeringRequest{
		UUID: arg,
		NetworkPeering: request.ModifyNetworkPeering{
			ConfiguredStatus: upcloud.NetworkPeeringConfiguredStatusDisabled,
		},
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: peering}, nil
}
