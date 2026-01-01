package account

import (
	"fmt"
	"sort"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// BillingCommand creates the 'account billing' command
func BillingCommand() commands.Command {
	return &billingCommand{
		BaseCommand: commands.New(
			"billing",
			"Show billing information",
			"upctl account billing --year 2025 --month 7",
		),
	}
}

type billingCommand struct {
	*commands.BaseCommand
	year       int
	month      int
	resourceID string
	username   string
}

// InitCommand implements Command.InitCommand
func (s *billingCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}

	flagSet.IntVar(&s.year, "year", 0, "Year for billing information.")
	flagSet.IntVar(&s.month, "month", 0, "Month for billing information.")
	flagSet.StringVar(&s.resourceID, "resource-id", "", "For IP addresses: the address itself, others, resource UUID")
	flagSet.StringVar(&s.username, "username", "", "Valid username")

	s.AddFlags(flagSet)

	commands.Must(s.Cobra().MarkFlagRequired("year"))
	commands.Must(s.Cobra().MarkFlagRequired("month"))
}

func firstElementAsString(row output.TableRow) string {
	s, ok := row[0].(string)
	if !ok {
		return ""
	}
	return s
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *billingCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if s.year < 1900 || s.year > 9999 {
		return nil, fmt.Errorf("invalid year: %d", s.year)
	}
	if s.month < 1 || s.month > 12 {
		return nil, fmt.Errorf("invalid month: %d", s.month)
	}

	svc := exec.Account()
	summary, err := svc.GetBillingSummary(exec.Context(), &request.GetBillingSummaryRequest{
		YearMonth:  fmt.Sprintf("%d-%02d", s.year, s.month),
		ResourceID: s.resourceID,
		Username:   s.username,
	})
	if err != nil {
		return nil, err
	}

	createCategorySections := func() []output.CombinedSection {
		var sections []output.CombinedSection
		var summaryRows []output.TableRow

		categories := map[string]*upcloud.BillingCategory{
			"Servers":                 summary.Servers,
			"Managed Databases":       summary.ManagedDatabases,
			"Managed Object Storages": summary.ManagedObjectStorages,
			"Managed Load Balancers":  summary.ManagedLoadbalancers,
			"Managed Kubernetes":      summary.ManagedKubernetes,
			"Network Gateways":        summary.NetworkGateways,
			"Networks":                summary.Networks,
			"Storages":                summary.Storages,
		}

		for categoryName, category := range categories {
			if category != nil {
				summaryRows = append(summaryRows, output.TableRow{categoryName, category.TotalAmount})
				resourceGroups := map[string]*upcloud.BillingResourceGroup{
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

	combined := output.Combined(createCategorySections())

	return output.MarshaledWithHumanOutput{
		Value:  summary,
		Output: combined,
	}, nil
}
