package account

import (
	"fmt"
	"math"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ShowCommand creates the 'account show' command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New("show", "Show account", "upctl account show"),
	}
}

type showCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *showCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.Account()
	account, err := svc.GetAccount(exec.Context())
	if err != nil {
		return nil, err
	}

	details := output.Details{
		Sections: []output.DetailSection{
			{
				Rows: []output.DetailRow{
					{Title: "Username:", Key: "username", Value: account.UserName},
					{Title: "Credits:", Key: "credits", Value: account.Credits, Format: formatCredits},
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
						Title: "Load balancers:",
						Key:   "load_balancers",
						Value: account.ResourceLimits.LoadBalancers,
					},
					{
						Title: "Managed object storages:",
						Key:   "managed_object_storages",
						Value: account.ResourceLimits.ManagedObjectStorages,
					},
					{
						Title: "Memory:",
						Key:   "memory",
						Value: account.ResourceLimits.Memory,
					},
					{
						Title: "Network peerings:",
						Key:   "network_peerings",
						Value: account.ResourceLimits.NetworkPeerings,
					},
					{
						Title: "Networks:",
						Key:   "networks",
						Value: account.ResourceLimits.Networks,
					},
					{
						Title: "NTP excess GiB:",
						Key:   "ntp_excess_gib",
						Value: account.ResourceLimits.NTPExcessGiB,
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
						Title: "Storage MaxIOPS:",
						Key:   "storage_maxiops",
						Value: account.ResourceLimits.StorageMaxIOPS,
					},
					{
						Title: "Storage SSD:",
						Key:   "storage_ssd",
						Value: account.ResourceLimits.StorageSSD,
					},
					{
						Title: "GPUs:",
						Key:   "gpus",
						Value: account.ResourceLimits.GPUs,
					},
				},
			},
		},
	}

	return output.MarshaledWithHumanOutput{
		Value:  account,
		Output: details,
	}, nil
}

func formatCredits(val interface{}) (text.Colors, string, error) {
	credits, ok := val.(float64)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse %T, expected float64", val)
	}

	if math.Abs(credits) < 0.001 {
		return nil, "Denied", nil
	}

	// Format does not follow european standards, but this is in sync with UI
	return nil, fmt.Sprintf("€%.2f", credits/100), nil
}
