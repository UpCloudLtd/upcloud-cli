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
  service        service.IpAddress
}

func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
  return func(args []string) (interface{}, error) {
    if len(args) != 1 {return nil, fmt.Errorf("one ip address or  is required")}
    ip, err := searchIpAddress(args[0], s.service)
    if err != nil {return nil, err}
    return ip, nil
  }
}

func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
  ip := out.(*upcloud.IPAddress)

  layout := ui.ListLayoutDefault
  layout.MarginLeft = false
  layout.MarginTop = false
  l := ui.NewListLayout(layout)
  {
    dCommon := ui.NewDetailsView()
    dCommon.AppendRows([]table.Row{
      {"Access:", ip.Access},
      {"Address:", ip.Address},
      {"Family:", ip.Family},
      {"PartOfPlan:", ip.PartOfPlan == 1},
      {"PTRRecord:", ip.PTRRecord},
      {"ServerUUID:", ip.ServerUUID},
      {"MAC:", ip.MAC},
      {"Floating:", ip.Floating == 1},
      {"Zone:", ip.Zone},
    })
    l.AppendSection("", dCommon.Render())
  }
  _, _ = fmt.Fprintln(writer, l.Render())
  return nil
}
