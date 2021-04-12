package network

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type deleteCommand struct {
	*commands.BaseCommand
	completion.Network
	resolver.CachingNetwork
}

// DeleteCommand creates the 'network delete' command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New("delete", "Delete a network"),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
}

// Execute implements Command.Execute
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	svc := exec.Network()
	msg := fmt.Sprintf("deleting network %v", arg)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	err := svc.DeleteNetwork(&request.DeleteNetworkRequest{
		UUID: arg,
	})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}
	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()
	return output.None{}, nil
}
