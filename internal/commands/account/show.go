package account

import (
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"math"
)

// ShowCommand creates the 'account show' command
func ShowCommand(service service.Account) commands.NewCommand {
	return &showCommand{
		BaseCommand: commands.New("show", "Show account"),
		service:     service,
	}
}

type showCommand struct {
	*commands.BaseCommand
	service service.Account
}

func (s *showCommand) NewParent() commands.NewCommand {
	return s.Parent().(commands.NewCommand)
}

func (s *showCommand) Execute(exec commands.Executor, args []string) (output.Command, error) {
	account, err := s.service.GetAccount()
	if err != nil {
		return nil, err
	}
	return output.Details{
		Sections: []output.DetailSection{
			{"", []output.DetailRow{
				{"Username:", account.UserName},
				{"Credits:", formatCredits(account.Credits)},
			}},
			{"Resource Limits:", []output.DetailRow{
				{"Cores:", account.ResourceLimits.Cores},
				{"Detached Floating IPs:", account.ResourceLimits.DetachedFloatingIps},
				{"Memory:", account.ResourceLimits.Memory},
				{"Networks:", account.ResourceLimits.Networks},
				{"Public IPv4:", account.ResourceLimits.PublicIPv4},
				{"Public IPv6:", account.ResourceLimits.PublicIPv6},
				{"Storage HDD:", account.ResourceLimits.StorageHDD},
				{"Storage SSD:", account.ResourceLimits.StorageSSD},
			}},
		},
	}, nil
}

func formatCredits(credits float64) string {
	if math.Abs(credits) < 0.001 {
		return "Denied"
	}
	return fmt.Sprintf("%.2f$", credits/100)
}
