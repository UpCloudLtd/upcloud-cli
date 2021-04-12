package router

import (
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
	name string
	resolver.CachingRouter
	completion.Router
}

// ModifyCommand creates the "router modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modify a router"),
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	// TODO: reimplmement
	// s.SetPositionalArgHelp(positionalArgHelp)
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.name, "name", "", "New router name.")
	s.AddFlags(fs)
}

// Execute implements command.Command
func (s *modifyCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if arg == "" {
		return nil, fmt.Errorf("router is required")
	}
	if s.name == "" {
		return nil, fmt.Errorf("name is required")
	}
	msg := fmt.Sprintf("Modifying router %s", s.name)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	res, err := exec.Network().ModifyRouter(&request.ModifyRouterRequest{UUID: arg, Name: s.name})
	if err != nil {
		logline.SetMessage(ui.LiveLogEntryErrorColours.Sprintf("%s: failed (%v)", msg, err.Error()))
		logline.SetDetails(err.Error(), "error: ")
		return nil, err
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.Marshaled{Value: res}, nil
}
