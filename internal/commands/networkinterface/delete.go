package networkinterface

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

type deleteCommand struct {
	*commands.BaseCommand
	interfaceIndex int
	resolver.CachingServer
	completion.Server
}

// DeleteCommand creates the "network-interface delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a network interface"),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	// TODO: reimplmement
	// s.SetPositionalArgHelp(server.PositionalArgHelp)
	fs := &pflag.FlagSet{}
	fs.IntVar(&s.interfaceIndex, "index", 0, "Interface index.")
	s.AddFlags(fs)
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *deleteCommand) MaximumExecutions() int {
	return maxNetworkInterfaceActions
}

// Execute implements command.Command
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if s.interfaceIndex == 0 {
		return nil, fmt.Errorf("interface index is required")
	}
	if arg == "" {
		return nil, errors.New("single server uuid is required")
	}
	msg := fmt.Sprintf("Deleting network interface %d of server %s", s.interfaceIndex, arg)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	err := exec.Network().DeleteNetworkInterface(&request.DeleteNetworkInterfaceRequest{
		ServerUUID: arg,
		Index:      s.interfaceIndex,
	})

	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.None{}, nil
}
