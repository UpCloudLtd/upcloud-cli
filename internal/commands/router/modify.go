package router

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
	name string
	resolver.CachingRouter
	completion.Router
}

// ModifyCommand creates the "router modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New(
			"modify",
			"Modify a router",
			"upctl router modify 04d031ab-4b85-4cbc-9f0e-6a2977541327 --name my_super_router",
			`upctl router modify "My Router" --name "My Turbo Router"`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.name, "name", "", "New router name.")
	s.AddFlags(fs)
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *modifyCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	if s.name == "" {
		return nil, fmt.Errorf("name is required")
	}
	msg := fmt.Sprintf("Modifying router %s", s.name)
	logline := exec.NewLogEntry(msg)
	logline.StartedNow()
	res, err := exec.Network().ModifyRouter(&request.ModifyRouterRequest{UUID: arg, Name: s.name})
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: done", msg))
	logline.MarkDone()

	return output.OnlyMarshaled{Value: res}, nil
}
