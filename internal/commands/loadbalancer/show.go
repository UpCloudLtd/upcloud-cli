package loadbalancer

import (
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ShowCommand creates the "loadbalancer show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show load balancer details",
			"upctl loadbalancer show 55199a44-4751-4e27-9394-7c7661910be3",
			"upctl loadbalancer show my-load-balancer",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingLoadBalancer
	completion.LoadBalancer
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	lb, err := svc.GetLoadBalancer(&request.GetLoadBalancerRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	var networkName string
	if network, err := svc.GetNetworkDetails(&request.GetNetworkDetailsRequest{UUID: lb.NetworkUUID}); err != nil {
		networkName = text.FgHiBlack.Sprint("Unknown")
	} else {
		networkName = network.Name
	}

	backEndRows := []output.TableRow{}
	for _, backEnd := range lb.Backends {
		resolver := text.FgHiBlack.Sprint("None")
		if backEnd.Resolver != "" {
			resolver = backEnd.Resolver
		}

		backEndRows = append(backEndRows, output.TableRow{
			backEnd.Name,
			resolver,
			len(backEnd.Members),
		})
	}

	frontEndRows := []output.TableRow{}
	for _, frontEnd := range lb.Frontends {
		frontEndRows = append(frontEndRows, output.TableRow{
			frontEnd.Name,
			frontEnd.Mode,
			frontEnd.Port,
			len(frontEnd.TLSConfigs),
			frontEnd.DefaultBackend,
			len(frontEnd.Rules),
		})
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: lb,
		Output: output.Combined{
			output.CombinedSection{
				Contents: output.Details{
					Sections: []output.DetailSection{
						{
							Title: "Overview:",
							Rows: []output.DetailRow{
								{Title: "UUID:", Value: lb.UUID, Colour: ui.DefaultUUUIDColours},
								{Title: "Name:", Value: lb.Name},
								{Title: "Plan:", Value: lb.Plan},
								{Title: "Zone:", Value: lb.Zone},
								{Title: "DNS name", Value: lb.DNSName},
								{Title: "Network name", Value: networkName},
								{Title: "Network UUID", Value: lb.NetworkUUID},
								{Title: "Operational state:", Value: lb.OperationalState, Colour: commands.LoadBalancerOperationalStateColour(lb.OperationalState)},
							},
						},
					},
				},
			},
			output.CombinedSection{
				Title: "Backends:",
				Contents: output.Table{
					Columns: []output.TableColumn{
						{Key: "name", Header: "Name"},
						{Key: "resolver", Header: "Resolver"},
						{Key: "members", Header: "Members"},
					},
					Rows: backEndRows,
				},
			},
			output.CombinedSection{
				Title: "Frontends:",
				Contents: output.Table{
					Columns: []output.TableColumn{
						{Key: "name", Header: "Name"},
						{Key: "mode", Header: "Mode"},
						{Key: "port", Header: "Port"},
						{Key: "tls_configs", Header: "TLS configs"},
						{Key: "default_backend", Header: "Default Backend"},
						{Key: "rules", Header: "Rules"},
					},
					Rows: frontEndRows,
				},
			},
		},
	}, nil
}
