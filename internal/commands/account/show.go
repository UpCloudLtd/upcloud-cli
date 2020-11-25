package account

import (
  "fmt"
  "github.com/UpCloudLtd/cli/internal/commands"
  "github.com/UpCloudLtd/cli/internal/ui"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
  "github.com/jedib0t/go-pretty/v6/table"
  "io"
)

func ShowCommand(service *service.Service) commands.Command {
  return &showCommand{
    BaseCommand: commands.New("show", "Show account"),
    service: service,
  }
}

type showCommand struct {
  *commands.BaseCommand
  service *service.Service
}

func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
  return func(args []string) (interface{}, error) {
    account, err := s.service.GetAccount()
    if err != nil {
      return nil, err
    }
    return account, nil
  }
}

func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
  account := out.(*upcloud.Account)

  layout := ui.ListLayoutDefault
  layout.MarginLeft = false
  layout.MarginTop = false
  l := ui.NewListLayout(layout)
  {
    dCommon := ui.NewDetailsView()
    dCommon.AppendRows([]table.Row{
      {"Username:", account.UserName},
      {"Credits:", account.Credits},
      {"Cores:", account.ResourceLimits.Cores},
      {"DetachedFloatingIps:", account.ResourceLimits.DetachedFloatingIps},
      {"Memory:", account.ResourceLimits.Memory},
      {"Networks:", account.ResourceLimits.Networks},
      {"PublicIPv4:", account.ResourceLimits.PublicIPv4},
      {"PublicIPv6:", account.ResourceLimits.PublicIPv6},
      {"StorageHDD:", account.ResourceLimits.StorageHDD},
      {"StorageSSD:", account.ResourceLimits.StorageSSD},
    })
    l.AppendSection("", dCommon.Render())
  }
  _, _ = fmt.Fprintln(writer, l.Render())
  return nil
}
