package router

import (
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/spf13/pflag"
)

const maxRouterActions = 10

// BaseRouterCommand creates the base "router" command
func BaseRouterCommand() commands.Command {
	return &routerCommand{
		commands.New("router", "Manage routers"),
	}
}

type routerCommand struct {
	*commands.BaseCommand
}

func handleStaticRoute(in string) (upcloud.StaticRoute, error) {
	result := upcloud.StaticRoute{}

	args, err := commands.Parse(in)
	if err != nil {
		return result, err
	}

	fs := &pflag.FlagSet{}
	fs.StringVar(&result.Name, "name", result.Name, "Name or description of the route.")
	fs.StringVar(&result.Nexthop, "nexthop", result.Nexthop, "Next hop address. NOTE: For static route to be active the next hop has to be an address of a reachable running Cloud Server in one of the Private Networks attached to the router.")
	fs.StringVar(&result.Route, "route", result.Route, "Destination prefix of the route.")

	return result, fs.Parse(args)
}
