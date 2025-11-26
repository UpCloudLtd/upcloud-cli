package account

import (
	"fmt"

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
	flagSet.StringVar(&s.resourceID, "resource_id", "", "For IP addresses: the address itself, others, resource UUID")
	flagSet.StringVar(&s.username, "username", "", "Valid username")

	s.AddFlags(flagSet)

	commands.Must(s.Cobra().MarkFlagRequired("year"))
	commands.Must(s.Cobra().MarkFlagRequired("month"))
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

	createCategorySections := func() []output.DetailSection {
		sections := []output.DetailSection{
			{
				Rows: []output.DetailRow{
					{Title: "Currency:", Key: "currency", Value: summary.Currency},
					{Title: "Total Amount:", Key: "total_amount", Value: summary.TotalAmount},
				},
			},
		}

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
				section := output.DetailSection{
					Title: categoryName,
					Rows: []output.DetailRow{
						{Title: "Total Amount:", Key: "total_amount", Value: category.TotalAmount},
					},
				}

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
					if group != nil {
						section.Rows = append(section.Rows, output.DetailRow{
							Title: fmt.Sprintf("%s Total:", groupName),
							Key:   fmt.Sprintf("%s_total", groupName),
							Value: group.TotalAmount,
						})

						// Add individual resources
						for i, resource := range group.Resources {
							section.Rows = append(section.Rows, output.DetailRow{
								Title: fmt.Sprintf("Resource %d ID:", i+1),
								Key:   fmt.Sprintf("%s_resource_%d_id", groupName, i+1),
								Value: resource.ResourceID,
							})
							section.Rows = append(section.Rows, output.DetailRow{
								Title: fmt.Sprintf("Resource %d Amount:", i+1),
								Key:   fmt.Sprintf("%s_resource_%d_amount", groupName, i+1),
								Value: resource.Amount,
							})
							section.Rows = append(section.Rows, output.DetailRow{
								Title: fmt.Sprintf("Resource %d Hours:", i+1),
								Key:   fmt.Sprintf("%s_resource_%d_hours", groupName, i+1),
								Value: resource.Hours,
							})
						}
					}
				}

				sections = append(sections, section)
			}
		}

		return sections
	}

	details := output.Details{
		Sections: createCategorySections(),
	}

	return output.MarshaledWithHumanOutput{
		Value:  summary,
		Output: details,
	}, nil
}
