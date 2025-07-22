package router

import (
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
)

// ShowCommand creates the "router show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show current router",
			"upctl router show 04d0a7f6-ee78-42b5-8077-6947f9e67c5a",
			"upctl router show 04d0a7f6-ee78-42b5-8077-6947f9e67c5a 04d031ab-4b85-4cbc-9f0e-6a2977541327",
			`upctl router show "My Turbo Router" my_super_router`,
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingRouter
	completion.Router
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	router, err := s.GetCached(arg)
	if err != nil {
		return nil, err
	}
	exec.Debug("got router", "uuid", router.UUID)
	networks, err := getNetworks(exec, router.AttachedNetworks)
	if err != nil {
		return nil, err
	}
	exec.Debug("got router networks", "networks", len(networks))
	networkRows := make([]output.TableRow, len(networks))
	for i, network := range networks {
		networkRows[i] = output.TableRow{
			network.UUID,
			network.Name,
			network.Type,
			network.Zone,
		}
	}
	staticRouteRows := make([]output.TableRow, len(router.StaticRoutes))
	for i, staticRoute := range router.StaticRoutes {
		staticRouteRows[i] = output.TableRow{
			staticRoute.Name,
			staticRoute.Route,
			staticRoute.Nexthop,
			staticRoute.Type,
		}
	}
	combined := output.Combined{
		output.CombinedSection{
			Key:   "",
			Title: "Common",
			Contents: output.Details{
				Sections: []output.DetailSection{
					{Rows: []output.DetailRow{
						{Key: "uuid", Title: "UUID:", Colour: ui.DefaultUUUIDColours, Value: router.UUID},
						{Key: "name", Title: "Name:", Value: router.Name},
						{Key: "type", Title: "Type:", Value: router.Type},
					}},
				},
			},
		},
		labels.GetLabelsSectionWithResourceType(router.Labels, "router"),
		output.CombinedSection{
			Key:   "networks",
			Title: "Networks:",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
					{Key: "name", Header: "Name"},
					{Key: "type", Header: "Type"},
					{Key: "zone", Header: "Zone"},
				},
				Rows: networkRows,
			},
		},
		output.CombinedSection{
			Key:   "static_routes",
			Title: "Static routes:",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "name", Header: "Name"},
					{Key: "route", Header: "Route", Colour: ui.DefaultAddressColours},
					{Key: "nexthop", Header: "Nexthop", Colour: ui.DefaultAddressColours},
					{Key: "type", Header: "Type"},
				},
				Rows:         staticRouteRows,
				EmptyMessage: "No static routes defined for this router.",
			},
		},
	}

	return output.MarshaledWithHumanOutput{
		Value:  router,
		Output: combined,
	}, nil
}

func getNetworks(exec commands.Executor, attached upcloud.RouterNetworkSlice) ([]upcloud.Network, error) {
	if len(attached) == 0 {
		return []upcloud.Network{}, nil
	}
	idleWorkers := make(chan int, maxRouterActions)
	for n := 0; n < maxRouterActions; n++ {
		idleWorkers <- n
	}
	results := make(chan *upcloud.Network)
	errors := make(chan error)
	for _, routerNetwork := range attached {
		go func(uuid string) {
			// get worker
			workerID := <-idleWorkers
			nw, err := exec.Network().GetNetworkDetails(exec.Context(), &request.GetNetworkDetailsRequest{UUID: uuid})
			if err != nil {
				errors <- err
			} else {
				results <- nw
			}
			// return worker
			idleWorkers <- workerID
		}(routerNetwork.NetworkUUID)
	}
	// collect results
	returns := make([]upcloud.Network, 0, len(attached))
	for {
		select {
		case err := <-errors:
			return nil, err
		case result := <-results:
			returns = append(returns, *result)
			if len(returns) >= len(attached) {
				return returns, nil
			}
		}
	}
}
