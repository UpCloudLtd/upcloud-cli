package networkpeering

import (
	"context"
	"fmt"
	"time"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// BaseNetworkPeeringCommand creates the base "networkpeering" command
func BaseNetworkPeeringCommand() commands.Command {
	baseCmd := commands.New("network-peering", "Manage network peerings")
	baseCmd.SetDeprecatedAliases([]string{"networkpeering"})

	return &networkpeeringCommand{
		BaseCommand: baseCmd,
	}
}

type networkpeeringCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (np *networkpeeringCommand) InitCommand() {
	np.Cobra().Aliases = []string{"np", "networkpeering"}
	// Deprecating networkpeering in favour of network-peering
	// TODO: Remove this in the future
	commands.SetDeprecationHelp(np.Cobra(), np.DeprecatedAliases())
}

// waitForNetworkPeeringState waits for network peering to reach given state and updates progress message with key matching given msg. Finally, progress message is updated back to given msg and either done state or timeout warning.
func waitForNetworkPeeringState(uuid string, state upcloud.NetworkPeeringState, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for network peering %s to be in %s state", uuid, state))

	ctx, cancel := context.WithTimeout(exec.Context(), 15*time.Minute)
	defer cancel()

	if _, err := exec.All().WaitForNetworkPeeringState(ctx, &request.WaitForNetworkPeeringStateRequest{
		UUID:         uuid,
		DesiredState: state,
	}); err != nil {
		exec.PushProgressUpdate(messages.Update{
			Key:     msg,
			Message: msg,
			Status:  messages.MessageStatusWarning,
			Details: "Error: " + err.Error(),
		})
		return
	}

	exec.PushProgressUpdateMessage(msg, msg)
	exec.PushProgressSuccess(msg)
}
