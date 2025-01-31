package router

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	name string
}

// CreateCommand creates the "router create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a router",
			"upctl router create --name my_router",
			`upctl router create --name "My Router"`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.name, "name", s.name, "Router name.")

	s.AddFlags(fs)
	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("name", cobra.NoFileCompletions))
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *createCommand) MaximumExecutions() int {
	return maxRouterActions
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	msg := fmt.Sprintf("Creating router %s", s.name)
	exec.PushProgressStarted(msg)

	res, err := exec.Network().CreateRouter(exec.Context(), &request.CreateRouterRequest{Name: s.name})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
