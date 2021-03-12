package ipaddress

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
)

// ShowCommand creates the 'ip-address show' command
func ShowCommand(service service.IpAddress) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show current ip address"),
		service:     service,
	}
}

type showCommand struct {
	*commands.BaseCommand
	service service.IpAddress
}

// InitCommand implements Command.InitCommand
func (s *showCommand) InitCommand() error {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(getArgCompFn(s.service))

	return nil
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("one ip address or PTR Record is required")
		}
		ip, err := searchIPAddress(args[0], s.service, true)
		if err != nil {
			return nil, err
		}
		return ip[0], nil
	}
}

// HandleOutput implements Command.HandleOutput
func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	ip := out.(*upcloud.IPAddress)

	layout := ui.ListLayoutDefault
	l := ui.NewListLayout(layout)
	{
		dCommon := ui.NewDetailsView()
		dCommon.Append(
			table.Row{"Address:", ui.DefaultAddressColours.Sprint(ip.Address)},
			table.Row{"Access:", ip.Access},
			table.Row{"Family:", ip.Family},
			table.Row{"Part of Plan:", ui.FormatBool(ip.PartOfPlan.Bool())},
			table.Row{"PTR Record:", ip.PTRRecord},
			table.Row{"Server UUID:", ui.DefaultUUUIDColours.Sprint(ip.ServerUUID)},
			table.Row{"MAC:", ip.MAC},
			table.Row{"Floating:", ui.FormatBool(ip.Floating.Bool())},
			table.Row{"Zone:", ip.Zone},
		)
		l.AppendSection("", dCommon.Render())
	}
	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}
