package networkpeering

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/spf13/pflag"

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

	disable config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (c *deleteCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &c.disable, "disable", false, "Disable the network peering before deleting it. This is equivalent to running 'upctl network-peering disable --wait` before 'upctl network-peering delete'.")
	c.AddFlags(flags)

	// Deprecating networkpeering in favour of network-peering
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"networkpeering"})
}

func Delete(exec commands.Executor, uuid string, disable bool) (output.Output, error) {
	if disable {
		_, err := disablePeering(exec, uuid, true)
		if err != nil {
			return nil, err
		}
	}

	svc := exec.All()
	msg := fmt.Sprintf("Deleting network peering %s", uuid)
	exec.PushProgressStarted(msg)

	err := svc.DeleteNetworkPeering(exec.Context(), &request.DeleteNetworkPeeringRequest{
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
	// Deprecating networkpeering in favour of network-peering
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"networkpeering"}, "network-peering")

	return Delete(exec, arg, c.disable.Value())
}
