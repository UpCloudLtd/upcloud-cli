package ipaddress

import (
	"errors"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"
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
	// TODO: reimplmement
	// s.SetPositionalArgHelp(positionalArgHelp)
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.mac, "mac", "", "MAC address of server interface to attach floating IP to.")
	fs.StringVar(&s.ptrrecord, "ptr-record", "", "A fully qualified domain name.")
	s.AddFlags(fs)
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *modifyCommand) MaximumExecutions() int {
	return maxIPAddressActions
}

// Execute implements Command.Execute
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
	return output.Marshaled{Value: res}, nil
}
