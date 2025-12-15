package billing

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// SummaryCommand creates the 'billing summary' command
func SummaryCommand() commands.Command {
	return &summaryCommand{
		BaseCommand: commands.New(
			"summary",
			"View billing summary for a specific period",
			"upctl billing summary --period 2024-01",
			"upctl billing summary --period 'last month'",
			"upctl billing summary --period '2months from 2024-06'",
			"upctl billing summary --period '+3months from 2024-01' --detailed",
			"upctl billing summary", // defaults to current month
		),
	}
}

type summaryCommand struct {
	*commands.BaseCommand
	period     string
	resourceID string
	detailed   bool
}

// InitCommand implements commands.Command
func (c *summaryCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}

	flagSet.StringVar(&c.period, "period", "month", "Billing period: 'month', 'quarter', 'year', 'YYYY-MM', relative like '3months', 'last month', or '2months from 2024-06'")
	flagSet.StringVar(&c.resourceID, "resource", "", "Filter by specific resource UUID")
	flagSet.BoolVar(&c.detailed, "detailed", false, "Show detailed breakdown of all resources")

	c.AddFlags(flagSet)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *summaryCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Parse period into YYYY-MM format
	yearMonth, periodDesc, err := parsePeriod(c.period)
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	msg := fmt.Sprintf("Fetching billing summary for %s", periodDesc)
	exec.PushProgressStarted(msg)

	req := &request.GetBillingSummaryRequest{
		YearMonth:  yearMonth,
		ResourceID: c.resourceID,
	}

	summary, err := exec.Account().GetBillingSummary(exec.Context(), req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	if c.detailed {
		return buildDetailedOutput(summary, periodDesc), nil
	}
	return buildSummaryOutput(summary, periodDesc), nil
}

func buildSummaryOutput(summary *upcloud.BillingSummary, month string) output.Output {
	rows := []output.TableRow{}

	// Add category rows
	if summary.Servers != nil && summary.Servers.TotalAmount > 0 {
		rows = append(rows, output.TableRow{
			"Servers",
			fmt.Sprintf("%.2f", summary.Servers.TotalAmount),
		})
	}

	if summary.Storages != nil && summary.Storages.TotalAmount > 0 {
		rows = append(rows, output.TableRow{
			"Storages",
			fmt.Sprintf("%.2f", summary.Storages.TotalAmount),
		})
	}

	if summary.ManagedDatabases != nil && summary.ManagedDatabases.TotalAmount > 0 {
		rows = append(rows, output.TableRow{
			"Managed Databases",
			fmt.Sprintf("%.2f", summary.ManagedDatabases.TotalAmount),
		})
	}

	if summary.ManagedObjectStorages != nil && summary.ManagedObjectStorages.TotalAmount > 0 {
		rows = append(rows, output.TableRow{
			"Object Storage",
			fmt.Sprintf("%.2f", summary.ManagedObjectStorages.TotalAmount),
		})
	}

	if summary.ManagedLoadbalancers != nil && summary.ManagedLoadbalancers.TotalAmount > 0 {
		rows = append(rows, output.TableRow{
			"Load Balancers",
			fmt.Sprintf("%.2f", summary.ManagedLoadbalancers.TotalAmount),
		})
	}

	if summary.ManagedKubernetes != nil && summary.ManagedKubernetes.TotalAmount > 0 {
		rows = append(rows, output.TableRow{
			"Kubernetes",
			fmt.Sprintf("%.2f", summary.ManagedKubernetes.TotalAmount),
		})
	}

	if summary.NetworkGateways != nil && summary.NetworkGateways.TotalAmount > 0 {
		rows = append(rows, output.TableRow{
			"Network Gateways",
			fmt.Sprintf("%.2f", summary.NetworkGateways.TotalAmount),
		})
	}

	if summary.Networks != nil && summary.Networks.TotalAmount > 0 {
		rows = append(rows, output.TableRow{
			"Networks",
			fmt.Sprintf("%.2f", summary.Networks.TotalAmount),
		})
	}

	// Add total with separator
	rows = append(rows, output.TableRow{
		"────────────────────",
		"────────────",
	})

	rows = append(rows, output.TableRow{
		"TOTAL",
		fmt.Sprintf("%.2f", summary.TotalAmount),
	})

	return output.MarshaledWithHumanOutput{
		Value: summary,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "category", Header: "Category"},
				{Key: "amount", Header: fmt.Sprintf("Amount (%s)", summary.Currency)},
			},
			Rows: rows,
		},
	}
}

func buildDetailedOutput(summary *upcloud.BillingSummary, month string) output.Output {
	rows := []output.TableRow{}

	// Process servers
	if summary.Servers != nil && summary.Servers.Server != nil {
		for _, resource := range summary.Servers.Server.Resources {
			rows = append(rows, output.TableRow{
				"Server",
				resource.ResourceID,
				fmt.Sprintf("%.2f", resource.Amount),
				fmt.Sprintf("%d", resource.Hours),
				summary.Currency,
			})
		}
	}

	// Process storages
	if summary.Storages != nil && summary.Storages.Storage != nil {
		for _, resource := range summary.Storages.Storage.Resources {
			rows = append(rows, output.TableRow{
				"Storage",
				resource.ResourceID,
				fmt.Sprintf("%.2f", resource.Amount),
				fmt.Sprintf("%d", resource.Hours),
				summary.Currency,
			})
		}
	}

	// Process managed databases
	if summary.ManagedDatabases != nil && summary.ManagedDatabases.ManagedDatabase != nil {
		for _, resource := range summary.ManagedDatabases.ManagedDatabase.Resources {
			rows = append(rows, output.TableRow{
				"Database",
				resource.ResourceID,
				fmt.Sprintf("%.2f", resource.Amount),
				fmt.Sprintf("%d", resource.Hours),
				summary.Currency,
			})
		}
	}

	// Process object storage
	if summary.ManagedObjectStorages != nil && summary.ManagedObjectStorages.ManagedObjectStorage != nil {
		for _, resource := range summary.ManagedObjectStorages.ManagedObjectStorage.Resources {
			rows = append(rows, output.TableRow{
				"Object Storage",
				resource.ResourceID,
				fmt.Sprintf("%.2f", resource.Amount),
				fmt.Sprintf("%d", resource.Hours),
				summary.Currency,
			})
		}
	}

	// Process load balancers
	if summary.ManagedLoadbalancers != nil && summary.ManagedLoadbalancers.ManagedLoadbalancer != nil {
		for _, resource := range summary.ManagedLoadbalancers.ManagedLoadbalancer.Resources {
			rows = append(rows, output.TableRow{
				"Load Balancer",
				resource.ResourceID,
				fmt.Sprintf("%.2f", resource.Amount),
				fmt.Sprintf("%d", resource.Hours),
				summary.Currency,
			})
		}
	}

	// Process Kubernetes
	if summary.ManagedKubernetes != nil && summary.ManagedKubernetes.ManagedKubernetes != nil {
		for _, resource := range summary.ManagedKubernetes.ManagedKubernetes.Resources {
			rows = append(rows, output.TableRow{
				"Kubernetes",
				resource.ResourceID,
				fmt.Sprintf("%.2f", resource.Amount),
				fmt.Sprintf("%d", resource.Hours),
				summary.Currency,
			})
		}
	}

	// Process network gateways
	if summary.NetworkGateways != nil && summary.NetworkGateways.NetworkGateway != nil {
		for _, resource := range summary.NetworkGateways.NetworkGateway.Resources {
			rows = append(rows, output.TableRow{
				"Network Gateway",
				resource.ResourceID,
				fmt.Sprintf("%.2f", resource.Amount),
				fmt.Sprintf("%d", resource.Hours),
				summary.Currency,
			})
		}
	}

	// Process networks (IPv4 addresses)
	if summary.Networks != nil && summary.Networks.IPv4Address != nil {
		for _, resource := range summary.Networks.IPv4Address.Resources {
			rows = append(rows, output.TableRow{
				"IPv4 Address",
				resource.ResourceID,
				fmt.Sprintf("%.2f", resource.Amount),
				fmt.Sprintf("%d", resource.Hours),
				summary.Currency,
			})
		}
	}

	// Add total
	rows = append(rows, output.TableRow{
		"────────────────",
		"────────────────────────────────────",
		"────────────",
		"──────",
		"────────",
	})

	rows = append(rows, output.TableRow{
		"TOTAL",
		"",
		fmt.Sprintf("%.2f", summary.TotalAmount),
		"",
		summary.Currency,
	})

	return output.MarshaledWithHumanOutput{
		Value: summary,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "type", Header: "Type"},
				{Key: "resource_id", Header: "Resource ID"},
				{Key: "amount", Header: fmt.Sprintf("Amount (%s)", summary.Currency)},
				{Key: "hours", Header: "Hours"},
			},
			Rows: rows,
		},
	}
}

// MaximumExecutions implements commands.Command
func (c *summaryCommand) MaximumExecutions() int {
	return 1
}
