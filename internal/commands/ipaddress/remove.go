package ipaddress

import (
	"errors"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type removeCommand struct {
	*commands.BaseCommand
}

// RemoveCommand creates the 'ip-address remove' command
func RemoveCommand() commands.Command {
	// TODO: should this be 'release'? inconsistent with libs now
	return &removeCommand{
		BaseCommand: commands.New("remove", "Remove an IP address"),
	}
}

// MaximumExecutions implements NewCommand.MaximumExecutions
func (s *removeCommand) MaximumExecutions() int {
	return maxIPAddressActions
}

// InitCommand implements Command.MakeExecuteCommand
func (s *removeCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	// TODO: reimplement
	// s.ArgCompletion(getArgCompFn(s.service))
}

// Execute implements NewCommand.Execute
func (s *removeCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if arg == "" {
		return nil, errors.New("need ip address to remove")
	}
	msg := fmt.Sprintf("removing ip-address %v", arg)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))
	err := exec.IPAddress().ReleaseIPAddress(&request.ReleaseIPAddressRequest{
		IPAddress: arg,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}
	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()
	return output.Marshaled{Value: nil}, nil
}
