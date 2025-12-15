package billing

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// ListCommand creates the 'billing list' command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List billing details with resource names",
			"upctl billing list --period 2024-01",
			"upctl billing list --period 'last month' --match web",
			"upctl billing list --period '3months from 2024-06' --category server",
			"upctl billing list", // defaults to current month
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
	period   string
	match    string
	category string
}

// resourceInfo holds resource information
type resourceInfo struct {
	Type     string
	UUID     string
	Name     string
	Amount   float64
	Hours    int
	Currency string
}

// InitCommand implements commands.Command
func (c *listCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}

	flagSet.StringVar(&c.period, "period", "month", "Billing period: 'month', 'quarter', 'year', 'YYYY-MM', relative like '3months', 'last month', or '2months from 2024-06'")
	flagSet.StringVar(&c.match, "match", "", "Filter resources by name (case-insensitive substring match)")
	flagSet.StringVar(&c.category, "category", "", "Filter by resource category (server, storage, database, etc.)")

	c.AddFlags(flagSet)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Parse period into YYYY-MM format
	yearMonth, periodDesc, err := parsePeriod(c.period)
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	msg := fmt.Sprintf("Fetching billing details for %s", periodDesc)
	exec.PushProgressStarted(msg)

	// Get billing summary
	req := &request.GetBillingSummaryRequest{
		YearMonth: yearMonth,
	}

	summary, err := exec.Account().GetBillingSummary(exec.Context(), req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	// Fetch resource names
	exec.PushProgressUpdateMessage(msg, msg+" (fetching resource names)")

	// Get all servers
	servers, err := exec.Server().GetServers(exec.Context())
	if err != nil {
		// Continue even if we can't get server names
		servers = &upcloud.Servers{}
	}

	// Get all storages
	storages, err := exec.Storage().GetStorages(exec.Context(), &request.GetStoragesRequest{})
	if err != nil {
		storages = &upcloud.Storages{}
	}

	// Get all databases - check if GetManagedDatabases exists
	// For now, we'll skip database names since the method might not be available
	databaseNames := make(map[string]string)

	// Create name lookup maps
	serverNames := make(map[string]string)
	for _, server := range servers.Servers {
		serverNames[server.UUID] = server.Title
		if server.Title == "" {
			serverNames[server.UUID] = server.Hostname
		}
	}

	storageNames := make(map[string]string)
	for _, storage := range storages.Storages {
		storageNames[storage.UUID] = storage.Title
	}

	exec.PushProgressSuccess(msg)

	// Build resource list
	resources := c.collectResources(summary, serverNames, storageNames, databaseNames)

	// Apply filters
	if c.match != "" {
		filtered := []resourceInfo{}
		for _, r := range resources {
			if strings.Contains(strings.ToLower(r.Name), strings.ToLower(c.match)) ||
				strings.Contains(strings.ToLower(r.UUID), strings.ToLower(c.match)) {
				filtered = append(filtered, r)
			}
		}
		resources = filtered
	}

	if c.category != "" {
		filtered := []resourceInfo{}
		categoryLower := strings.ToLower(c.category)
		for _, r := range resources {
			if strings.Contains(strings.ToLower(r.Type), categoryLower) {
				filtered = append(filtered, r)
			}
		}
		resources = filtered
	}

	return c.buildOutput(resources, summary.TotalAmount, summary.Currency, periodDesc), nil
}

func (c *listCommand) collectResources(summary *upcloud.BillingSummary, serverNames, storageNames, databaseNames map[string]string) []resourceInfo {
	resources := []resourceInfo{}

	// Process servers
	if summary.Servers != nil && summary.Servers.Server != nil {
		for _, resource := range summary.Servers.Server.Resources {
			name := serverNames[resource.ResourceID]
			if name == "" {
				name = "<unnamed>"
			}
			resources = append(resources, resourceInfo{
				Type:     "Server",
				UUID:     resource.ResourceID,
				Name:     name,
				Amount:   resource.Amount,
				Hours:    resource.Hours,
				Currency: summary.Currency,
			})
		}
	}

	// Process storages
	if summary.Storages != nil && summary.Storages.Storage != nil {
		for _, resource := range summary.Storages.Storage.Resources {
			name := storageNames[resource.ResourceID]
			if name == "" {
				name = "<unnamed storage>"
			}
			resources = append(resources, resourceInfo{
				Type:     "Storage",
				UUID:     resource.ResourceID,
				Name:     name,
				Amount:   resource.Amount,
				Hours:    resource.Hours,
				Currency: summary.Currency,
			})
		}
	}

	// Process backups
	if summary.Storages != nil && summary.Storages.Backup != nil {
		for _, resource := range summary.Storages.Backup.Resources {
			// Backups might be related to storages
			name := storageNames[resource.ResourceID]
			if name == "" {
				name = "<backup>"
			}
			resources = append(resources, resourceInfo{
				Type:     "Backup",
				UUID:     resource.ResourceID,
				Name:     name,
				Amount:   resource.Amount,
				Hours:    resource.Hours,
				Currency: summary.Currency,
			})
		}
	}

	// Process managed databases
	if summary.ManagedDatabases != nil && summary.ManagedDatabases.ManagedDatabase != nil {
		for _, resource := range summary.ManagedDatabases.ManagedDatabase.Resources {
			name := databaseNames[resource.ResourceID]
			if name == "" {
				name = "<unnamed database>"
			}
			resources = append(resources, resourceInfo{
				Type:     "Database",
				UUID:     resource.ResourceID,
				Name:     name,
				Amount:   resource.Amount,
				Hours:    resource.Hours,
				Currency: summary.Currency,
			})
		}
	}

	// Process object storage
	if summary.ManagedObjectStorages != nil && summary.ManagedObjectStorages.ManagedObjectStorage != nil {
		for _, resource := range summary.ManagedObjectStorages.ManagedObjectStorage.Resources {
			resources = append(resources, resourceInfo{
				Type:     "Object Storage",
				UUID:     resource.ResourceID,
				Name:     resource.ResourceID, // Object storage might not have names
				Amount:   resource.Amount,
				Hours:    resource.Hours,
				Currency: summary.Currency,
			})
		}
	}

	// Process load balancers
	if summary.ManagedLoadbalancers != nil && summary.ManagedLoadbalancers.ManagedLoadbalancer != nil {
		for _, resource := range summary.ManagedLoadbalancers.ManagedLoadbalancer.Resources {
			resources = append(resources, resourceInfo{
				Type:     "Load Balancer",
				UUID:     resource.ResourceID,
				Name:     resource.ResourceID, // We'd need to fetch load balancer names separately
				Amount:   resource.Amount,
				Hours:    resource.Hours,
				Currency: summary.Currency,
			})
		}
	}

	// Process Kubernetes
	if summary.ManagedKubernetes != nil && summary.ManagedKubernetes.ManagedKubernetes != nil {
		for _, resource := range summary.ManagedKubernetes.ManagedKubernetes.Resources {
			resources = append(resources, resourceInfo{
				Type:     "Kubernetes",
				UUID:     resource.ResourceID,
				Name:     resource.ResourceID, // We'd need to fetch K8s cluster names separately
				Amount:   resource.Amount,
				Hours:    resource.Hours,
				Currency: summary.Currency,
			})
		}
	}

	// Process network gateways
	if summary.NetworkGateways != nil && summary.NetworkGateways.NetworkGateway != nil {
		for _, resource := range summary.NetworkGateways.NetworkGateway.Resources {
			resources = append(resources, resourceInfo{
				Type:     "Network Gateway",
				UUID:     resource.ResourceID,
				Name:     resource.ResourceID,
				Amount:   resource.Amount,
				Hours:    resource.Hours,
				Currency: summary.Currency,
			})
		}
	}

	// Process networks (IPv4 addresses)
	if summary.Networks != nil && summary.Networks.IPv4Address != nil {
		for _, resource := range summary.Networks.IPv4Address.Resources {
			resources = append(resources, resourceInfo{
				Type:     "IPv4 Address",
				UUID:     resource.ResourceID,
				Name:     resource.ResourceID,
				Amount:   resource.Amount,
				Hours:    resource.Hours,
				Currency: summary.Currency,
			})
		}
	}

	return resources
}

func (c *listCommand) buildOutput(resources []resourceInfo, totalAmount float64, currency, month string) output.Output {
	rows := []output.TableRow{}

	for _, r := range resources {
		rows = append(rows, output.TableRow{
			r.Type,
			r.Name,
			r.UUID,
			fmt.Sprintf("%.2f", r.Amount),
			fmt.Sprintf("%d", r.Hours),
			r.Currency,
		})
	}

	// Calculate subtotal if filtered
	var subtotal float64
	for _, r := range resources {
		subtotal += r.Amount
	}

	// Add separator and total
	if len(rows) > 0 {
		rows = append(rows, output.TableRow{
			"---",
			"---",
			"---",
			"---",
			"---",
			"---",
		})

		if c.match != "" || c.category != "" {
			// Show filtered subtotal
			rows = append(rows, output.TableRow{
				"SUBTOTAL (filtered)",
				"",
				"",
				fmt.Sprintf("%.2f", subtotal),
				"",
				currency,
			})
		}

		rows = append(rows, output.TableRow{
			"TOTAL (all resources)",
			"",
			"",
			fmt.Sprintf("%.2f", totalAmount),
			"",
			currency,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: struct {
			Resources []resourceInfo `json:"resources"`
			Total     float64        `json:"total"`
			Currency  string         `json:"currency"`
		}{
			Resources: resources,
			Total:     totalAmount,
			Currency:  currency,
		},
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "type", Header: "Type"},
				{Key: "name", Header: "Name"},
				{Key: "uuid", Header: "UUID"},
				{Key: "amount", Header: "Amount"},
				{Key: "hours", Header: "Hours"},
				{Key: "currency", Header: "Currency"},
			},
			Rows: rows,
		},
	}
}

// MaximumExecutions implements commands.Command
func (c *listCommand) MaximumExecutions() int {
	return 1
}
