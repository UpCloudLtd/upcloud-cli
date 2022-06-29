package ipaddress

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
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
		BaseCommand: commands.New(
			"modify",
			"Modify an IP address",
			"upctl ip-address modify 185.70.196.225 --ptr-record myapp.com",
			"upctl ip-address modify 185.70.197.175 --mac d6:0e:4a:6f:2b:06",
		),
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
	msg := fmt.Sprintf("Modifying ip-address %v", arg)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()
	logline.SetMessage(fmt.Sprintf("%s: sending request", msg))
	res, err := exec.IPAddress().ModifyIPAddress(&request.ModifyIPAddressRequest{
		IPAddress: arg,
		MAC:       s.mac,
		PTRRecord: s.ptrrecord,
	})
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
