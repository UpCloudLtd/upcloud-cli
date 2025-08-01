package ipaddress

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
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
	for _, flag := range []string{"mac", "ptr-record"} {
		commands.Must(s.Cobra().RegisterFlagCompletionFunc(flag, cobra.NoFileCompletions))
	}
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *modifyCommand) MaximumExecutions() int {
	return maxIPAddressActions
}

// Execute implements commands.MultipleArgumentCommand
func (s *modifyCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Modifying ip-address %v", arg)
	exec.PushProgressStarted(msg)

	res, err := exec.IPAddress().ModifyIPAddress(exec.Context(), &request.ModifyIPAddressRequest{
		IPAddress: arg,
		MAC:       s.mac,
		PTRRecord: s.ptrrecord,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
