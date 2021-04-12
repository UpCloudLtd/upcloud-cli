package router

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	name string
}

// CreateCommand creates the "router create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New("create", "Create a router"),
	}
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.name, "name", s.name, "Router name.")
	s.AddFlags(fs)
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *createCommand) MaximumExecutions() int {
	return maxRouterActions
}

// Execute implements command.Command
func (s *createCommand) Execute(exec commands.Executor, _ string) (output.Output, error) {
	// TODO: should this be a regular argument?
	if s.name == "" {
		return nil, fmt.Errorf("name is required")
	}
	msg := fmt.Sprintf("Creating router %s", s.name)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()

	res, err := exec.Network().CreateRouter(&request.CreateRouterRequest{Name: s.name})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
