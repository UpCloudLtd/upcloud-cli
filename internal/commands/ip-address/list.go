package ip_address

import (
  "fmt"
  "github.com/UpCloudLtd/cli/internal/commands"
  "github.com/UpCloudLtd/cli/internal/ui"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
  "github.com/jedib0t/go-pretty/v6/table"
  "github.com/spf13/pflag"
  "io"
)

func ListCommand(service service.IpAddress) commands.Command {
  return &listCommand{
    BaseCommand: commands.New("list", "List ip addresses"),
    service:     service,
  }
}

type listCommand struct {
  *commands.BaseCommand
  service        service.IpAddress
  header         table.Row
  columnKeys     []string
  visibleColumns []string
}

func (s *listCommand) InitCommand() {
  s.header = table.Row{"Access", "Address", "Family", "PartOfPlan", "PTRRecord", "Server", "MAC", "Floating", "Zone"}
  s.columnKeys = []string{"access", "address", "family", "partofplan", "ptrrecord", "server", "mac", "floating", "zone"}
  s.visibleColumns = []string{"access", "address", "family", "partofplan", "ptrrecord", "server", "mac", "floating", "zone"}
  flags := &pflag.FlagSet{}
  s.AddVisibleColumnsFlag(flags, &s.visibleColumns, s.columnKeys, s.visibleColumns)
  s.AddFlags(flags)
}

func (s *listCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
  return func(args []string) (interface{}, error) {
    ips, err := s.service.GetIPAddresses()
    if err != nil {
      return nil, err
    }
    return ips, nil
  }
}

func (s *listCommand) HandleOutput(writer io.Writer, out interface{}) error {
  ips := out.(*upcloud.IPAddresses)

  t := ui.NewDataTable(s.columnKeys...)
  t.OverrideColumnKeys(s.visibleColumns...)
  t.SetHeader(s.header)

  for _, ip := range ips.IPAddresses {
    t.AppendRow(table.Row{
      ip.Access,
      ip.Address,
      ip.Family,
      ip.PartOfPlan == 1,
      ip.PTRRecord,
      ip.ServerUUID,
      ip.MAC,
      ip.Floating == 1,
      ip.Zone})
  }

  fmt.Fprintln(writer)
  fmt.Fprintln(writer, t.Render())
  fmt.Fprintln(writer)
  return nil

  return nil
}
