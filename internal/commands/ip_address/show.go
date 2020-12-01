package ip_address

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
)

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

func (s *showCommand) InitCommand() {
	s.SetPositionalArgHelp(positionalArgHelp)
	s.ArgCompletion(GetArgCompFn(s.service))
}

func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("one ip address or PTR Record is required")
		}
		ip, err := searchIpAddress(args[0], s.service, true)
		if err != nil {
			return nil, err
		}
		return ip[0], nil
	}
}

func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	ip := out.(*upcloud.IPAddress)

	layout := ui.ListLayoutDefault
	l := ui.NewListLayout(layout)
	{
		dCommon := ui.NewDetailsView()
		dCommon.AppendRows([]table.Row{
			{"Address:", ui.DefaultAddressColours.Sprint(ip.Address)},
			{"Access:", ip.Access},
			{"Family:", ip.Family},
			{"Part of Plan:", ui.FormatBool(ip.PartOfPlan.Bool())},
			{"PTR Record:", ip.PTRRecord},
			{"Server UUID:", ui.DefaultUuidColours.Sprint(ip.ServerUUID)},
			{"MAC:", ip.MAC},
			{"Floating:", ui.FormatBool(ip.Floating.Bool())},
			{"Zone:", ip.Zone},
		})
		l.AppendSection("", dCommon.Render())
	}
	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}
