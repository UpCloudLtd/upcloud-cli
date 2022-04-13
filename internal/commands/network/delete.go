package network

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type deleteCommand struct {
	*commands.BaseCommand
	completion.Network
	resolver.CachingNetwork
}

// DeleteCommand creates the 'network delete' command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a network",
			"upctl network delete 037f260c-9568-4d9b-97e5-44cf52440ccb",
			"upctl network delete 03d7b5c2-b80a-4636-88d4-f9911185c975 0312a237-8204-4c1c-9fd1-2314013ec687",
			`upctl network delete "My Network 1"`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
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
