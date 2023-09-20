package router

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	name         string
	staticRoutes []string
}

// CreateCommand creates the "router create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a router",
			"upctl router create --name my_router",
			`upctl router create --name "My Router" --static-route name=my_static_route,nexthop=10.0.0.100,route=0.0.0.0/0"`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.name, "name", s.name, "Router name.")
	fs.StringArrayVar(
		&s.staticRoutes,
		"static-route",
		[]string{},
		"Static route for this router, multiple can be declared.\n\n "+
			"Fields: \n"+
			"  name: string \n"+
			"  nexthop: string \n"+
			"  route: string")

	s.AddFlags(fs)
	_ = s.Cobra().MarkFlagRequired("name")
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *createCommand) MaximumExecutions() int {
	return maxRouterActions
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	msg := fmt.Sprintf("Creating router %s", s.name)
	exec.PushProgressStarted(msg)

	var staticRoutes []upcloud.StaticRoute
	for _, v := range s.staticRoutes {
		staticRoute, err := handleStaticRoute(v)
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}

		staticRoutes = append(staticRoutes, staticRoute)
	}

	res, err := exec.Network().CreateRouter(exec.Context(), &request.CreateRouterRequest{
		Name:         s.name,
		StaticRoutes: staticRoutes,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
