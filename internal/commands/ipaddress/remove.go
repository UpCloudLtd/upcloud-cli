package ipaddress

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
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
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))
	err := exec.IPAddress().ReleaseIPAddress(&request.ReleaseIPAddressRequest{
		IPAddress: arg,
	})
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.None{}, nil
}
