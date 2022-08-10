package networkinterface

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"
)

type deleteCommand struct {
	*commands.BaseCommand
	interfaceIndex int
	resolver.CachingServer
	completion.Server
}

// DeleteCommand creates the "network-interface delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a network interface",
			"upctl server network-interface delete 009d7f4e-99ce-4c78-88f1-e695d4c37743 --index 1",
			"upctl server network-interface delete my_server --index 7",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.IntVar(&s.interfaceIndex, "index", 0, "Interface index.")

	s.AddFlags(fs)
	s.Cobra().MarkFlagRequired("index") //nolint:errcheck
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *deleteCommand) MaximumExecutions() int {
	return maxNetworkInterfaceActions
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *deleteCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Deleting network interface %d of server %s", s.interfaceIndex, arg)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	err := exec.Network().DeleteNetworkInterface(&request.DeleteNetworkInterfaceRequest{
		ServerUUID: arg,
		Index:      s.interfaceIndex,
	})
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.None{}, nil
}
