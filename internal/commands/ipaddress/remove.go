package ipaddress

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type removeCommand struct {
	*commands.BaseCommand
	completion.IPAddress
	resolver.CachingIPAddress
}

// RemoveCommand creates the 'ip-address remove' command
func RemoveCommand() commands.Command {
	// TODO: should this be 'release'? inconsistent with libs now
	return &removeCommand{
		BaseCommand: commands.New(
			"remove",
			"Remove an IP address",
			"upctl ip-address remove 185.70.197.44",
			"upctl ip-address remove 2a04:3544:8000:1000:d40e:4aff:fe6f:2c85",
		),
	}
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *removeCommand) MaximumExecutions() int {
	return maxIPAddressActions
}

// InitCommand implements Command.MakeExecuteCommand
func (s *removeCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
func (s *removeCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Removing ip-address %v", arg)
	exec.PushProgressStarted(msg)

	err := exec.IPAddress().ReleaseIPAddress(&request.ReleaseIPAddressRequest{
		IPAddress: arg,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
