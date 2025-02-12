package loadbalancer

import (
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ShowCommand creates the "loadbalancer show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show load balancer details",
			"upctl load-balancer show 55199a44-4751-4e27-9394-7c7661910be3",
			"upctl load-balancer show my-load-balancer",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
	resolver.CachingLoadBalancer
	completion.LoadBalancer
}

func (s *showCommand) InitCommand() {
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(s, []string{"loadbalancer"})
}

// Execute implements commands.MultipleArgumentCommand
func (s *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	// Deprecating loadbalancer in favour of load-balancer
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(s, []string{"loadbalancer"}, "load-balancer")

	svc := exec.All()
	lb, err := svc.GetLoadBalancer(exec.Context(), &request.GetLoadBalancerRequest{UUID: uuid})
	if err != nil {
		return nil, err
	}

	var networkName string
	if network, err := svc.GetNetworkDetails(exec.Context(), &request.GetNetworkDetailsRequest{UUID: lb.NetworkUUID}); err != nil {
		networkName = ""
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
			len(backEnd.TLSConfigs),
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

	resolverRows := []output.TableRow{}
	for _, resolver := range lb.Resolvers {
		var nameservers []string
		for _, nameserver := range resolver.Nameservers {
			nameservers = append(nameservers, ui.DefaultAddressColours.Sprint(nameserver))
		}

		resolverRows = append(resolverRows, output.TableRow{
			resolver.Name,
			strings.Join(nameservers, ", "),
		})
	}

	combined := output.Combined{
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
							{Title: "Network name", Value: networkName, Format: format.PossiblyUnknownString},
							{Title: "Network UUID", Value: lb.NetworkUUID},
							{Title: "Operational state:", Value: lb.OperationalState, Format: format.LoadBalancerState},
						},
					},
				},
			},
		},
		labels.GetLabelsSection(lb.Labels),
		output.CombinedSection{
			Title: "Backends:",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "name", Header: "Name"},
					{Key: "resolver", Header: "Resolver"},
					{Key: "members", Header: "Members"},
					{Key: "tls_configs", Header: "TLS configs"},
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
		output.CombinedSection{
			Title: "Resolvers:",
			Contents: output.Table{
				Columns: []output.TableColumn{
					{Key: "name", Header: "Name"},
					{Key: "nameservers", Header: "Nameservers"},
				},
				Rows: resolverRows,
			},
		},
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value:  lb,
		Output: combined,
	}, nil
}
