package gateway

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/paging"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/pflag"
)

// ListCommand creates the "gateway list" command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List gateways", "upctl gateway list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
	paging.PageParameters
}

func (c *listCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	c.ConfigureFlags(fs)
	c.AddFlags(fs)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	gateways, err := svc.GetGateways(exec.Context(), c.Page())
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, gtw := range gateways {
		rows = append(rows, output.TableRow{
			gtw.UUID,
			gtw.Name,
			gtw.Routers,
			gtw.OperationalState,
			gtw.Zone,
		})
	}

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: gateways,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "uuid", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "name", Header: "Name"},
				{Key: "routers", Header: "Routers", Format: formatRouters},
				{Key: "status", Header: "Status"},
				{Key: "zone", Header: "Zone"},
			},
			Rows: rows,
		},
	}, nil
}

func formatRouters(val interface{}) (text.Colors, string, error) {
	routers, ok := val.([]upcloud.GatewayRouter)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse routers from %T, expected []upcloud.GatewayRouter", val)
	}

	var rows []string
	for _, rt := range routers {
		rows = append(rows, ui.DefaultUUUIDColours.Sprint(rt.UUID))
	}

	return nil, strings.Join(rows, ",\n"), nil
}
