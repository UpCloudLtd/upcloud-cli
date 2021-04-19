package ipaddress

import (
	"errors"
	"fmt"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	mac       string
	ptrrecord string
	resolver.CachingIPAddress
	completion.IPAddress
}

// ModifyCommand creates the 'ip-address modify' command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify an IP address"),
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.mac, "mac", "", "MAC address of server interface to attach floating IP to.")
	fs.StringVar(&s.ptrrecord, "ptr-record", "", "New fully qualified domain name to set as the PTR record for the IP address.")
	s.AddFlags(fs)
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *modifyCommand) MaximumExecutions() int {
	return maxIPAddressActions
}

// Execute implements commands.MultipleArgumentCommand
func (s *modifyCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if arg == "" {
		return nil, errors.New("need ip address to modify")
	}
	msg := fmt.Sprintf("modifying ip-address %v", arg)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))
	res, err := exec.IPAddress().ModifyIPAddress(&request.ModifyIPAddressRequest{
		IPAddress: arg,
		MAC:       s.mac,
		PTRRecord: s.ptrrecord,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}
	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()
	return output.OnlyMarshaled{Value: res}, nil
}
