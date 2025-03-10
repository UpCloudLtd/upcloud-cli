package networkpeering

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DisableCommand creates the "networkpeering disable" command
func DisableCommand() commands.Command {
	return &disableCommand{
		BaseCommand: commands.New(
			"disable",
			"Disable a network peering",
			"upctl network-peering disable 8abc8009-4325-4b23-4321-b1232cd81231",
			"upctl network-peering disable my-network-peering",
		),
	}
}

type disableCommand struct {
	*commands.BaseCommand
	resolver.CachingNetworkPeering
	completion.NetworkPeering

	wait config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (c *disableCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &c.wait, "wait", false, "Wait for network peering to be in disabled state before returning.")
	c.AddFlags(flags)
	// Deprecating networkpeering in favour of network-peering
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"networkpeering"})
}

func disablePeering(exec commands.Executor, uuid string, wait bool) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Disabling network peering %v", uuid)

	// Use UUID as message key to allow updating the message later
	exec.PushProgressStarted(msg)

	peering, err := svc.ModifyNetworkPeering(exec.Context(), &request.ModifyNetworkPeeringRequest{
		UUID: uuid,
		NetworkPeering: request.ModifyNetworkPeering{
			ConfiguredStatus: upcloud.NetworkPeeringConfiguredStatusDisabled,
		},
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if wait {
		waitForNetworkPeeringState(uuid, upcloud.NetworkPeeringStateDisabled, exec, msg)
	} else {
		exec.PushProgressSuccess(msg)
	}

	return output.OnlyMarshaled{Value: peering}, err
}

// Execute implements commands.MultipleArgumentCommand
func (c *disableCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	// Deprecating networkpeering in favour of network-peering
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"networkpeering"}, "network-peering")

	return disablePeering(exec, arg, c.wait.Value())
}
