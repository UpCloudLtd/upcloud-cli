package router

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/completion"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

// ShowCommand creates the "router show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show current router"),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingRouter
	completion.Router
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() {
	// TODO: reimplmement
	// s.SetPositionalArgHelp(positionalArgHelp)
}

// Execute implements command.Command
func (s *showCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if arg == "" {
		return nil, fmt.Errorf("one router uuid or name is required")
	}
	router, err := s.CachingRouter.GetCached(arg)
	if err != nil {
		return nil, err
	}
	networks, err := getNetworks(exec, router.AttachedNetworks)
	if err != nil {
		return nil, err
	}
	networkRows := make([]output.TableRow, len(networks))
	for i, network := range networks {
		networkRows[i] = output.TableRow{network.UUID, network.Name, network.Router, network.Type, network.Zone}
	}
	return output.Combined{
		output.CombinedSection{
			Key:   "",
			Title: "Common",
			Contents: output.Details{
				Sections: []output.DetailSection{
					{Rows: []output.DetailRow{
						{Key: "uuid", Title: "UUID:", Color: ui.DefaultUUUIDColours, Value: router.UUID},
						{Key: "name", Title: "Name:", Value: router.Name},
						{Key: "type", Title: "Type:", Value: router.Type},
					}},
				},
			},
		},
		output.CombinedSection{
			Key:   "networks",
			Title: "Networks:",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "uuid", Header: "UUID", Color: ui.DefaultUUUIDColours},
					{Key: "name", Header: "Name"},
					{Key: "router", Header: "Router", Color: ui.DefaultUUUIDColours},
					{Key: "type", Header: "Type"},
					{Key: "zone", Header: "Zone"},
				},
				Rows: networkRows,
			},
		},
	}, nil
}

func getNetworks(exec commands.Executor, attached upcloud.RouterNetworkSlice) ([]upcloud.Network, error) {
	if len(attached) == 0 {
		return []upcloud.Network{}, nil
	}
	var idleWorkers = make(chan int, maxRouterActions)
	for n := 0; n < maxRouterActions; n++ {
		idleWorkers <- n
	}
	results := make(chan *upcloud.Network)
	errors := make(chan error)
	for _, routerNetwork := range attached {
		go func(uuid string) {
			// get worker
			workerID := <-idleWorkers
			nw, err := exec.Network().GetNetworkDetails(&request.GetNetworkDetailsRequest{UUID: uuid})
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
