package account

import (
	"fmt"
	"math"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
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

func (s *showCommand) MakeExecutor() commands.CommandExecutor {
	return func(args []string) (output.Command, error) {
		svc := s.Config().Service.(service.Account)
		account, err := svc.GetAccount()
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
}

func formatCredits(credits float64) string {
	if math.Abs(credits) < 0.001 {
		return "Denied"
	}
	return fmt.Sprintf("%.2f$", credits/100)
}
