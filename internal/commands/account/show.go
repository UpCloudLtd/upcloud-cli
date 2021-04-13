package account

import (
	"fmt"
	"math"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
)

// ShowCommand creates the 'account show' command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show account", ""),
	}
}

type showCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *showCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.Account()
	account, err := svc.GetAccount()
	if err != nil {
		return nil, err
	}
	return output.Details{
		Sections: []output.DetailSection{
			{
				Rows: []output.DetailRow{
					{Title: "Username:", Key: "username", Value: account.UserName},
					{Title: "Credits:", Key: "credits", Value: formatCredits(account.Credits)},
				},
			},
			{
				Title: "Resource Limits:", Key: "resource_limits", Rows: []output.DetailRow{
					{
						Title: "Cores:",
						Key:   "cores",
						Value: account.ResourceLimits.Cores,
					},
					{
						Title: "Detached Floating IPs:",
						Key:   "detached_floating_ips",
						Value: account.ResourceLimits.DetachedFloatingIps,
					},
					{
						Title: "Memory:",
						Key:   "memory",
						Value: account.ResourceLimits.Memory,
					},
					{
						Title: "Networks:",
						Key:   "networks",
						Value: account.ResourceLimits.Networks,
					},
					{
						Title: "Public IPv4:",
						Key:   "public_ipv4",
						Value: account.ResourceLimits.PublicIPv4,
					},
					{
						Title: "Public IPv6:",
						Key:   "public_ipv6",
						Value: account.ResourceLimits.PublicIPv6,
					},
					{
						Title: "Storage HDD:",
						Key:   "storage_hdd",
						Value: account.ResourceLimits.StorageHDD,
					},
					{
						Title: "Storage SSD:",
						Key:   "storage_ssd",
						Value: account.ResourceLimits.StorageSSD,
					},
				},
			},
		},
	}, nil
}

func formatCredits(credits float64) string {
	if math.Abs(credits) < 0.001 {
		return "Denied"
	}
	return fmt.Sprintf("%.2f$", credits/100)
}
