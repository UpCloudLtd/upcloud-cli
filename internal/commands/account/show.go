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

func ShowCommand(service service.Account) commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show account"),
		service:     service,
	}
}

type showCommand struct {
	*commands.BaseCommand
	service service.Account
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

	var credits = "0"

	if account.Credits == 0.0 {
		credits = "Denied"
	} else {
		credits = fmt.Sprintf("%.f$", account.Credits / 100)
	}

	layout := ui.ListLayoutDefault
	l := ui.NewListLayout(layout)
	{
		dCommon := ui.NewDetailsView()
		dCommon.AppendRows([]table.Row{
			{"Username:", account.UserName},
			{"Credits:", credits},
		})
		l.AppendSection("", dCommon.Render())
	}

	{
		dCommon := ui.NewDetailsView()
		dCommon.SetHeaderWidth(25)
		dCommon.AppendRows([]table.Row{
			{"Cores:", account.ResourceLimits.Cores},
			{"Detached Floating IPs:", account.ResourceLimits.DetachedFloatingIps},
			{"Memory:", account.ResourceLimits.Memory},
			{"Networks:", account.ResourceLimits.Networks},
			{"Public IPv4:", account.ResourceLimits.PublicIPv4},
			{"Public IPv6:", account.ResourceLimits.PublicIPv6},
			{"Storage HDD:", account.ResourceLimits.StorageHDD},
			{"Storage SSD:", account.ResourceLimits.StorageSSD},
		})
		l.AppendSection("Resource Limits:", dCommon.Render())
	}
	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}
