package account

import (
	"fmt"
	"math"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/mapper"
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

func (s *showCommand) MaximumExecutions() int {
	return 1
}

func (s *showCommand) ArgumentMapper() (mapper.Argument, error) {
	return nil, nil
}

func (s *showCommand) Execute(_ commands.Executor, _ string) (output.Command, error) {
	svc := s.Config().Service.(service.Account)
	account, err := svc.GetAccount()
	if err != nil {
		return nil, err
	}
	return output.Details{
		Sections: []output.DetailSection{
			{Title: "", Rows: []output.DetailRow{
				{Title: "Username:", Value: account.UserName},
				{Title: "Credits:", Value: formatCredits(account.Credits)},
			}},
			{Title: "Resource Limits:", Rows: []output.DetailRow{
				{Title: "Cores:", Value: account.ResourceLimits.Cores},
				{Title: "Detached Floating IPs:", Value: account.ResourceLimits.DetachedFloatingIps},
				{Title: "Memory:", Value: account.ResourceLimits.Memory},
				{Title: "Networks:", Value: account.ResourceLimits.Networks},
				{Title: "Public IPv4:", Value: account.ResourceLimits.PublicIPv4},
				{Title: "Public IPv6:", Value: account.ResourceLimits.PublicIPv6},
				{Title: "Storage HDD:", Value: account.ResourceLimits.StorageHDD},
				{Title: "Storage SSD:", Value: account.ResourceLimits.StorageSSD},
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
