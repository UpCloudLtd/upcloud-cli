package networkinterface

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
	commands.Must(s.Cobra().MarkFlagRequired("index"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("index", cobra.NoFileCompletions))
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *deleteCommand) MaximumExecutions() int {
	return maxNetworkInterfaceActions
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *deleteCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Deleting network interface %d of server %s", s.interfaceIndex, arg)
	exec.PushProgressStarted(msg)

	err := exec.Network().DeleteNetworkInterface(exec.Context(), &request.DeleteNetworkInterfaceRequest{
		ServerUUID: arg,
		Index:      s.interfaceIndex,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
