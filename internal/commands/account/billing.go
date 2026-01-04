package account

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// BillingCommand creates the enhanced 'account billing' command with backward compatibility
func BillingCommand() commands.Command {
	return &billingCommand{
		BaseCommand: commands.New(
			"billing",
			"Show billing information",
			"upctl account billing", // defaults to current month
			"upctl account billing --period 'last month'", // flexible period
			"upctl account billing --year 2025 --month 7", // backward compatible
		),
	}
}

type billingCommand struct {
	*commands.BaseCommand
	// Legacy flags (kept for backward compatibility)
	year  int
	month int

	// Enhanced flags
	period     string
	resourceID string
	username   string
	match      string
	category   string
	detailed   bool
}

// InitCommand implements Command.InitCommand
func (s *billingCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}

	// Legacy flags - keep exact same names and descriptions for backward compatibility
	flagSet.IntVar(&s.year, "year", 0, "Year for billing information.")
	flagSet.IntVar(&s.month, "month", 0, "Month for billing information.")
	flagSet.StringVar(&s.resourceID, "resource-id", "", "For IP addresses: the address itself, others, resource UUID")
	flagSet.StringVar(&s.username, "username", "", "Valid username")

	// New enhanced flags
	flagSet.StringVar(&s.period, "period", "", "Billing period: 'month', 'last month', '3months', 'YYYY-MM', or '2months from 2024-06'")
	flagSet.StringVar(&s.match, "match", "", "Filter resources by name (case-insensitive substring)")
	flagSet.StringVar(&s.category, "category", "", "Filter by category: server, storage, database, load-balancer, kubernetes, gateway")
	flagSet.BoolVar(&s.detailed, "detailed", false, "Show detailed breakdown with resource names")

	s.AddFlags(flagSet)

	// Only mark as required if no period flag is provided
	// This maintains backward compatibility while allowing new usage
	if s.period == "" {
		// Note: We'll handle this logic in ExecuteWithoutArguments instead
	}
}

// parsePeriod converts various period formats into YYYY-MM for the API
// Supports formats like: "month", "last month", "3months", "2024-07", "2months from 2024-06"
func parsePeriod(period string) (string, string, error) {
	now := time.Now()

	// Handle YYYY-MM format directly
	if matched, _ := fmt.Sscanf(period, "%d-%d", new(int), new(int)); matched == 2 {
		return period, period, nil
	}

	// Handle named periods
	switch strings.ToLower(period) {
	case "month", "current", "":
		yearMonth := now.Format("2006-01")
		return yearMonth, fmt.Sprintf("current month (%s)", yearMonth), nil
	case "day", "today":
		yearMonth := now.Format("2006-01")
		return yearMonth, fmt.Sprintf("today (%s)", now.Format("2006-01-02")), nil
	case "quarter":
		quarter := (now.Month()-1)/3 + 1
		yearMonth := now.Format("2006-01")
		return yearMonth, fmt.Sprintf("Q%d %d (current month: %s)", quarter, now.Year(), yearMonth), nil
	case "year":
		yearMonth := now.Format("2006-01")
		return yearMonth, fmt.Sprintf("year %d (current month: %s)", now.Year(), yearMonth), nil
	case "last month":
		lastMonth := now.AddDate(0, -1, 0)
		yearMonth := lastMonth.Format("2006-01")
		return yearMonth, fmt.Sprintf("last month (%s)", yearMonth), nil
	case "last quarter":
		lastQuarter := now.AddDate(0, -3, 0)
		quarter := (lastQuarter.Month()-1)/3 + 1
		yearMonth := now.Format("2006-01")
		return yearMonth, fmt.Sprintf("Q%d %d (current month: %s)", quarter, lastQuarter.Year(), yearMonth), nil
	case "last year":
		lastYear := now.AddDate(-1, 0, 0)
		yearMonth := lastYear.Format("2006-01")
		return yearMonth, fmt.Sprintf("last year %d (current month: %s)", lastYear.Year(), yearMonth), nil
	}

	// Handle relative from base date (e.g., "2months from 2024-06")
	if strings.Contains(period, " from ") {
		parts := strings.Split(period, " from ")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid relative period format: %s", period)
		}

		relPeriod := parts[0]
		baseDate := parts[1]

		baseTime, err := time.Parse("2006-01", baseDate)
		if err != nil {
			return "", "", fmt.Errorf("invalid base date format: %s (use YYYY-MM)", baseDate)
		}

		forward := strings.HasPrefix(relPeriod, "+")
		relPeriod = strings.TrimPrefix(relPeriod, "+")
		relPeriod = strings.TrimPrefix(relPeriod, "-")

		var amount int
		var unit string
		if matched, _ := fmt.Sscanf(relPeriod, "%d%s", &amount, &unit); matched == 2 {
			multiplier := 1
			if !forward {
				multiplier = -1
			}

			var targetTime time.Time
			switch strings.ToLower(unit) {
			case "month", "months":
				targetTime = baseTime.AddDate(0, amount*multiplier, 0)
			case "year", "years":
				targetTime = baseTime.AddDate(amount*multiplier, 0, 0)
			default:
				return "", "", fmt.Errorf("unsupported unit for relative period: %s", unit)
			}

			yearMonth := targetTime.Format("2006-01")
			direction := "before"
			if forward {
				direction = "after"
			}
			return yearMonth, fmt.Sprintf("%d %s %s %s (%s)", amount, unit, direction, baseDate, yearMonth), nil
		}
	}

	// Handle simple relative periods (e.g., "3months", "2weeks")
	var amount int
	var unit string
	if matched, _ := fmt.Sscanf(period, "%d%s", &amount, &unit); matched == 2 {
		var targetTime time.Time
		switch strings.ToLower(unit) {
		case "day", "days":
			targetTime = now.AddDate(0, 0, -amount)
		case "week", "weeks":
			targetTime = now.AddDate(0, 0, -amount*7)
		case "month", "months":
			targetTime = now.AddDate(0, -amount, 0)
		case "year", "years":
			targetTime = now.AddDate(-amount, 0, 0)
		default:
			return "", "", fmt.Errorf("unknown period unit: %s (use day/week/month/year)", unit)
		}
		yearMonth := targetTime.Format("2006-01")
		return yearMonth, fmt.Sprintf("%d %s ago (%s)", amount, unit, yearMonth), nil
	}

	return "", "", fmt.Errorf("unrecognized period format: %s", period)
}

func firstElementAsString(row output.TableRow) string {
	if len(row) == 0 {
		return ""
	}
	s, ok := row[0].(string)
	if !ok {
		return ""
	}
	return s
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *billingCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	var yearMonth string
	var err error

	// Determine the period to query - three-way priority:
	// 1. If --period is specified, use it
	// 2. If --year and --month are specified, use them (backward compatibility)
	// 3. Default to current month
	if s.period != "" {
		yearMonth, _, err = parsePeriod(s.period)
		if err != nil {
			return nil, err
		}
	} else if s.year > 0 && s.month > 0 {
		// Legacy behavior - exact same validation as original
		if s.year < 1900 || s.year > 9999 {
			return nil, fmt.Errorf("invalid year: %d", s.year)
		}
		if s.month < 1 || s.month > 12 {
			return nil, fmt.Errorf("invalid month: %d", s.month)
		}
		yearMonth = fmt.Sprintf("%d-%02d", s.year, s.month)
	} else if s.year > 0 || s.month > 0 {
		// Maintain original behavior - both must be set if either is
		return nil, fmt.Errorf("both --year and --month must be specified together")
	} else {
		// New default behavior - current month
		yearMonth, _, err = parsePeriod("")
		if err != nil {
			return nil, err
		}
	}

	svc := exec.Account()
	summary, err := svc.GetBillingSummary(exec.Context(), &request.GetBillingSummaryRequest{
		YearMonth:  yearMonth,
		ResourceID: s.resourceID,
		Username:   s.username,
	})
	if err != nil {
		return nil, err
	}

	// Fetch resource names if detailed view is requested
	var resourceNames map[string]string
	if s.detailed {
		resourceNames = s.fetchResourceNames(exec, summary)
	}

	// Build output sections (enhanced or original based on flags)
	var sections []output.CombinedSection
	if s.detailed || s.match != "" || s.category != "" {
		sections = s.buildEnhancedSections(summary, resourceNames)
	} else {
		sections = s.buildOriginalSections(summary)
	}

	return output.MarshaledWithHumanOutput{
		Value:  summary,
		Output: output.Combined(sections),
	}, nil
}

// fetchResourceNames retrieves names for servers and storage resources
func (s *billingCommand) fetchResourceNames(exec commands.Executor, summary *upcloud.BillingSummary) map[string]string {
	names := make(map[string]string)

	// Fetch server names
	if summary.Servers != nil && summary.Servers.Server != nil {
		servers, _ := exec.Server().GetServers(exec.Context())
		if servers != nil {
			for _, server := range servers.Servers {
				names[server.UUID] = server.Title
			}
		}
	}

	// Fetch storage names
	if summary.Storages != nil && summary.Storages.Storage != nil {
		storages, _ := exec.Storage().GetStorages(exec.Context(), &request.GetStoragesRequest{})
		if storages != nil {
			for _, storage := range storages.Storages {
				names[storage.UUID] = storage.Title
			}
		}
	}

	return names
}

// getCategories returns all billing categories from the summary
func getCategories(summary *upcloud.BillingSummary) map[string]*upcloud.BillingCategory {
	return map[string]*upcloud.BillingCategory{
		"Servers":                 summary.Servers,
		"Managed Databases":       summary.ManagedDatabases,
		"Managed Object Storages": summary.ManagedObjectStorages,
		"Managed Load Balancers":  summary.ManagedLoadbalancers,
		"Managed Kubernetes":      summary.ManagedKubernetes,
		"Network Gateways":        summary.NetworkGateways,
		"Networks":                summary.Networks,
		"Storages":                summary.Storages,
	}
}

// getResourceGroups returns all resource groups from a billing category
func getResourceGroups(category *upcloud.BillingCategory) map[string]*upcloud.BillingResourceGroup {
	return map[string]*upcloud.BillingResourceGroup{
		"Server":                 category.Server,
		"Managed Database":       category.ManagedDatabase,
		"Managed Object Storage": category.ManagedObjectStorage,
		"Managed Load Balancer":  category.ManagedLoadbalancer,
		"Managed Kubernetes":     category.ManagedKubernetes,
		"Network Gateway":        category.NetworkGateway,
		"IPv4 Address":           category.IPv4Address,
		"Backup":                 category.Backup,
		"Storage":                category.Storage,
		"Template":               category.Template,
	}
}

// buildOriginalSections maintains exact original output format for backward compatibility
func (s *billingCommand) buildOriginalSections(summary *upcloud.BillingSummary) []output.CombinedSection {
	var sections []output.CombinedSection
	var summaryRows []output.TableRow

	categories := getCategories(summary)

	for categoryName, category := range categories {
		if category != nil {
			summaryRows = append(summaryRows, output.TableRow{categoryName, category.TotalAmount})
			resourceGroups := getResourceGroups(category)

			for groupName, group := range resourceGroups {
				if group != nil && len(group.Resources) > 0 {
					var resourceRows []output.TableRow
					for _, resource := range group.Resources {
						resourceRows = append(resourceRows, output.TableRow{
							resource.ResourceID,
							resource.Amount,
							resource.Hours,
						})
					}

					sections = append(sections, output.CombinedSection{
						Key:   fmt.Sprintf("%s_%s_resources", categoryName, groupName),
						Title: fmt.Sprintf("%s - %s Resources:", categoryName, groupName),
						Contents: output.Table{
							Columns: []output.TableColumn{
								{Key: "resource_id", Header: "Resource ID"},
								{Key: "amount", Header: "Amount"},
								{Key: "hours", Header: "Hours"},
							},
							Rows:         resourceRows,
							EmptyMessage: fmt.Sprintf("No resources for %s.", groupName),
						},
					})
				}
			}
		}
	}

	sort.Slice(summaryRows, func(i, j int) bool {
		return firstElementAsString(summaryRows[i]) < firstElementAsString(summaryRows[j])
	})
	summaryRows = append(summaryRows, output.TableRow{"Total", summary.TotalAmount})

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Title < sections[j].Title
	})
	sections = append([]output.CombinedSection{{
		Key:   "summary",
		Title: "Summary:",
		Contents: output.Table{
			Columns: []output.TableColumn{
				{Key: "resource", Header: "Resource"},
				{Key: "total_amount", Header: "Amount"},
			},
			Rows: summaryRows,
		},
	}}, sections...)
	return sections
}

// buildEnhancedSections provides enhanced output with names and filtering
func (s *billingCommand) buildEnhancedSections(summary *upcloud.BillingSummary, resourceNames map[string]string) []output.CombinedSection {
	var sections []output.CombinedSection
	var summaryRows []output.TableRow

	categories := getCategories(summary)

	// Apply category filter if specified
	if s.category != "" {
		filtered := make(map[string]*upcloud.BillingCategory)
		categoryLower := strings.ToLower(s.category)
		for name, cat := range categories {
			if strings.Contains(strings.ToLower(name), categoryLower) {
				filtered[name] = cat
			}
		}
		categories = filtered
	}

	for categoryName, category := range categories {
		if category != nil {
			summaryRows = append(summaryRows, output.TableRow{categoryName, category.TotalAmount})

			if s.detailed {
				resourceRows := s.buildResourceRows(category, resourceNames)
				if len(resourceRows) > 0 {
					sections = append(sections, output.CombinedSection{
						Key:   fmt.Sprintf("%s_resources", strings.ReplaceAll(strings.ToLower(categoryName), " ", "_")),
						Title: fmt.Sprintf("%s Resources:", categoryName),
						Contents: output.Table{
							Columns: []output.TableColumn{
								{Key: "resource_id", Header: "Resource ID"},
								{Key: "name", Header: "Name"},
								{Key: "amount", Header: "Amount"},
								{Key: "hours", Header: "Hours"},
							},
							Rows:         resourceRows,
							EmptyMessage: fmt.Sprintf("No resources for %s", categoryName),
						},
					})
				}
			}
		}
	}

	sort.Slice(summaryRows, func(i, j int) bool {
		return firstElementAsString(summaryRows[i]) < firstElementAsString(summaryRows[j])
	})
	summaryRows = append(summaryRows, output.TableRow{"Total", summary.TotalAmount})

	sections = append([]output.CombinedSection{{
		Key:   "summary",
		Title: "Summary:",
		Contents: output.Table{
			Columns: []output.TableColumn{
				{Key: "resource", Header: "Resource"},
				{Key: "total_amount", Header: "Amount"},
			},
			Rows: summaryRows,
		},
	}}, sections...)

	return sections
}

func (s *billingCommand) buildResourceRows(category *upcloud.BillingCategory, resourceNames map[string]string) []output.TableRow {
	var rows []output.TableRow

	resourceGroups := getResourceGroups(category)

	for _, group := range resourceGroups {
		if group != nil {
			for _, resource := range group.Resources {
				name := resourceNames[resource.ResourceID]
				if name == "" {
					name = "-"
				}

				// Apply name filter if specified
				if s.match != "" {
					if !strings.Contains(strings.ToLower(name), strings.ToLower(s.match)) &&
						!strings.Contains(strings.ToLower(resource.ResourceID), strings.ToLower(s.match)) {
						continue
					}
				}

				rows = append(rows, output.TableRow{
					resource.ResourceID,
					name,
					resource.Amount,
					resource.Hours,
				})
			}
		}
	}

	return rows
}
