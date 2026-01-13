package server

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/pflag"
)

const (
	durationHour  = "hour"
	durationMonth = "month"
)

// PlanListCommand creates the "server plans" command
func PlanListCommand() commands.Command {
	return &planListCommand{
		BaseCommand: commands.New("plans", "List server plans", "upctl server plans"),
	}
}

type planListCommand struct {
	*commands.BaseCommand
	pricingZone     string
	pricingDuration string
}

// InitCommand initializes the command flags
func (s *planListCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.pricingZone, "pricing", "", "Show pricing for the specified zone (e.g., de-fra1)")
	flagSet.StringVar(&s.pricingDuration, "pricing-duration", durationMonth, "Duration for pricing calculation (e.g., 'hour', 'month', '1h', '24h')")

	s.BaseCommand.Cobra().Flags().AddFlagSet(flagSet)
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *planListCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	plansObj, err := exec.All().GetPlans(exec.Context())
	if err != nil {
		return nil, err
	}

	plans := plansObj.Plans
	sort.Slice(plans, func(i, j int) bool {
		if plans[i].CoreNumber != plans[j].CoreNumber {
			return plans[i].CoreNumber < plans[j].CoreNumber
		}

		if plans[i].MemoryAmount != plans[j].MemoryAmount {
			return plans[i].MemoryAmount < plans[j].MemoryAmount
		}

		return plans[i].StorageSize < plans[j].StorageSize
	})

	// Check if pricing should be shown
	showPricing := s.pricingZone != ""

	// Validate that pricing-duration is only used with --pricing
	if !showPricing && s.pricingDuration != "1m" {
		// User specified pricing-duration without specifying a pricing zone
		return nil, fmt.Errorf("--pricing-duration requires --pricing zone to be specified")
	}

	// Parse pricing duration only if showing pricing
	var duration time.Duration

	// Use 28 days per month for pricing calculations (UpCloud bills max 28 days per month)
	month := 28 * 24 * time.Hour

	if showPricing {
		// Handle special keywords first
		switch strings.ToLower(s.pricingDuration) {
		case durationHour:
			duration = 1 * time.Hour
		case durationMonth:
			// Use 28 days per month for pricing calculations (UpCloud bills max 28 days per month)
			duration = month
		default:
			// Parse as standard duration
			var err error
			duration, err = time.ParseDuration(s.pricingDuration)
			if err != nil {
				return nil, fmt.Errorf("invalid pricing-duration: %s (use formats like 'hourly', 'monthly', '1h', '24h')", s.pricingDuration)
			}
		}
	}

	// Fetch pricing information if requested
	var pricing map[string]upcloud.Price
	if showPricing {
		pricingByZone, err := exec.All().GetPricingByZone(exec.Context())
		switch {
		case err != nil:
			// Continue without pricing - just show plans
			showPricing = false
		case pricingByZone != nil:
			// Find the requested zone
			var ok bool
			pricing, ok = (*pricingByZone)[s.pricingZone]
			if !ok {
				return nil, fmt.Errorf("pricing zone %s not found", s.pricingZone)
			}
		default:
			// priceZones is nil, disable pricing
			showPricing = false
		}
	}

	rows := make(map[string][]output.TableRow)
	for _, p := range plans {
		key := planType(p)
		row := output.TableRow{
			p.Name,
			p.CoreNumber,
			p.MemoryAmount,
			p.StorageSize,
			p.StorageTier,
			p.PublicTrafficOut,
		}

		// Add GPU fields only for GPU plans
		if key == "gpu" {
			row = append(row, p.GPUModel, p.GPUAmount)
		}

		// Add cost if requested
		if showPricing && pricing != nil {
			cost := getPlanCost(p, pricing, duration)
			row = append(row, cost)
		}

		rows[key] = append(rows[key], row)
	}

	return output.MarshaledWithHumanOutput{
		Value: plans,
		Output: output.Combined{
			planSection("general_purpose", "General purpose", rows["general_purpose"], showPricing, s.pricingDuration),
			planSection("gpu", "GPU", rows["gpu"], showPricing, s.pricingDuration),
			planSection("cloud_native", "Cloud native", rows["cloud_native"], showPricing, s.pricingDuration),
			planSection("high_cpu", "High CPU", rows["high_cpu"], showPricing, s.pricingDuration),
			planSection("high_memory", "High memory", rows["high_memory"], showPricing, s.pricingDuration),
			planSection("developer", "Developer", rows["developer"], showPricing, s.pricingDuration),
		},
	}, nil
}

func planType(p upcloud.Plan) string {
	if strings.HasPrefix(p.Name, "DEV-") {
		return "developer"
	}
	if strings.HasPrefix(p.Name, "HICPU-") {
		return "high_cpu"
	}
	if strings.HasPrefix(p.Name, "HIMEM-") {
		return "high_memory"
	}
	if strings.HasPrefix(p.Name, "CLOUDNATIVE-") {
		return "cloud_native"
	}
	if strings.HasPrefix(p.Name, "GPU-") {
		return "gpu"
	}
	return "general_purpose"
}

func planSection(key, title string, rows []output.TableRow, showPricing bool, pricingDuration string) output.CombinedSection {
	columns := []output.TableColumn{
		{Key: "name", Header: "Name"},
		{Key: "cores", Header: "Cores"},
		{Key: "memory", Header: "Memory"},
		{Key: "storage", Header: "Storage size"},
		{Key: "storage_tier", Header: "Storage tier"},
		{Key: "egress_transfer", Header: "Transfer out (GiB/month)"},
	}

	if key == "gpu" {
		columns = append(columns,
			output.TableColumn{Key: "gpu_model", Header: "GPU model"},
			output.TableColumn{Key: "gpu_amount", Header: "GPU amount"},
		)
	}

	if showPricing {
		decimals := 3
		if pricingDuration == durationMonth {
			decimals = 2
		}

		columns = append(columns, output.TableColumn{
			Key:    "cost",
			Header: formatPricingHeader(pricingDuration),
			Format: getFormatPrice(decimals),
		})
	}

	return output.CombinedSection{
		Key:   key,
		Title: title,
		Contents: output.Table{
			Columns: columns,
			Rows:    rows,
		},
	}
}

// getPlanCost calculates the cost for a given plan
func getPlanCost(plan upcloud.Plan, pricing map[string]upcloud.Price, duration time.Duration) float64 {
	if pricing == nil {
		return math.NaN()
	}

	fieldName := "server_plan_" + plan.Name

	price, ok := pricing[fieldName]
	if !ok {
		return math.NaN()
	}

	hourlyPrice := price.Price / 100

	// Calculate cost for the requested duration
	// UpCloud bills per (starting) hour, so round up to next full hour
	return hourlyPrice * math.Ceil(duration.Hours())
}

// formatPricingHeader creates a human-readable header for the cost column
func formatPricingHeader(pricingDuration string) string {
	switch strings.ToLower(pricingDuration) {
	case durationHour:
		return "Price (per hour)"
	case durationMonth:
		return "Price (per month)"
	case "1h":
		return "Price (per hour)"
	case "24h":
		return "Price (per day)"
	default:
		// For other durations, just display the duration string
		return fmt.Sprintf("Price (per %s)", pricingDuration)
	}
}

func getFormatPrice(decimals int) func(any) (text.Colors, string, error) {
	format := fmt.Sprintf("%%.%df", decimals)
	return func(val any) (text.Colors, string, error) {
		price, ok := val.(float64)
		if !ok {
			return nil, "", fmt.Errorf("cannot parse price from %T, expected string", val)
		}
		if math.IsNaN(price) {
			return text.Colors{text.FgHiBlack}, "unknown", nil
		}

		return nil, fmt.Sprintf(format, price), nil
	}
}
