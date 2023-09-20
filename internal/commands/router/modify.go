package router

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	resolver.CachingRouter
	completion.Router
	name         string
	staticRoutes []string
}

// ModifyCommand creates the "router modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New(
			"modify",
			"Modify a router",
			"upctl router modify 04d031ab-4b85-4cbc-9f0e-6a2977541327 --name my_super_router",
			`upctl router modify "My Router" --name "My Turbo Router" --static-route name=my_static_route,nexthop=10.0.0.100,route=0.0.0.0/0"`,
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.name, "name", "", "New router name.")
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

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *modifyCommand) ExecuteSingleArgument(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Modifying router %s", s.name)
	exec.PushProgressStarted(msg)

	var staticRoutes []upcloud.StaticRoute
	for _, v := range s.staticRoutes {
		staticRoute, err := handleStaticRoute(v)
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}

		staticRoutes = append(staticRoutes, staticRoute)
	}

	res, err := exec.Network().ModifyRouter(exec.Context(), &request.ModifyRouterRequest{UUID: arg, Name: s.name, StaticRoutes: &staticRoutes})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
