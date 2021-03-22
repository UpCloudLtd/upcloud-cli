package account

import (
	"fmt"
	"io"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/jedib0t/go-pretty/v6/table"
)

// ShowCommand creates the 'account show' command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show account"),
	}
}

type showCommand struct {
	*commands.BaseCommand
}

// MakeExecuteCommand implements command.MakeExecuteCommand
func (s *showCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		svc := s.Config().Service.(service.Account)
		account, err := svc.GetAccount()
		if err != nil {
			return nil, err
		}
		return account, nil
	}
}

// HandleOutput implements command.HandleOutput
func (s *showCommand) HandleOutput(writer io.Writer, out interface{}) error {
	account := out.(*upcloud.Account)

	var credits string

	if account.Credits == 0.0 {
		credits = "Denied"
	} else {
		credits = fmt.Sprintf("%.2f$", account.Credits/100)
	}

	layout := ui.ListLayoutDefault
	l := ui.NewListLayout(layout)
	{
		dCommon := ui.NewDetailsView()
		dCommon.Append(
			table.Row{"Username:", account.UserName},
			table.Row{"Credits:", credits},
		)
		l.AppendSection("", dCommon.Render())
	}

	{
		dCommon := ui.NewDetailsView()
		dCommon.SetHeaderWidth(25)
		dCommon.Append(
			table.Row{"Cores:", account.ResourceLimits.Cores},
			table.Row{"Detached Floating IPs:", account.ResourceLimits.DetachedFloatingIps},
			table.Row{"Memory:", account.ResourceLimits.Memory},
			table.Row{"Networks:", account.ResourceLimits.Networks},
			table.Row{"Public IPv4:", account.ResourceLimits.PublicIPv4},
			table.Row{"Public IPv6:", account.ResourceLimits.PublicIPv6},
			table.Row{"Storage HDD:", account.ResourceLimits.StorageHDD},
			table.Row{"Storage SSD:", account.ResourceLimits.StorageSSD},
		)
		l.AppendSection("Resource Limits:", dCommon.Render())
	}
	_, _ = fmt.Fprintln(writer, l.Render())
	return nil
}
