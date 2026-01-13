package server

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/UpCloudLtd/progress/messages"
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
	pricesZone     string
	pricesDuration string
}

// InitCommand initializes the command flags
func (s *planListCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.pricesZone, "prices", "", "Show prices for the specified zone (e.g., de-fra1)")
	flagSet.StringVar(&s.pricesDuration, "prices-duration", durationMonth, "Duration for prices calculation (e.g., 'hour', 'month', '1h', '24h')")

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

	// Check if prices should be shown
	showPrices := s.pricesZone != ""

	// Validate that prices-duration is only used with --prices
	if !showPrices && s.pricesDuration != "month" {
		// User specified prices-duration without specifying a prices zone
		return nil, fmt.Errorf("--prices-duration requires --prices zone to be specified")
	}

	// Parse prices duration only if showing prices
	var duration time.Duration

	// Use 28 days per month for prices calculations (UpCloud bills max 28 days per month)
	month := 28 * 24 * time.Hour

	if showPrices {
		// Handle special keywords first
		switch strings.ToLower(s.pricesDuration) {
		case durationHour:
			duration = 1 * time.Hour
		case durationMonth:
			// Use 28 days per month for pricing calculations (UpCloud bills max 28 days per month)
			duration = month
		default:
			// Parse as standard duration
			var err error
			duration, err = time.ParseDuration(s.pricesDuration)
			if err != nil {
				return nil, fmt.Errorf("invalid prices-duration: %s (use formats like 'hour', 'month', '1h', '24h')", s.pricesDuration)
			}
		}
	}

	// Fetch pricing information if requested
	var prices map[string]upcloud.Price
	if showPrices {
		pricesByZone, err := exec.All().GetPricesByZone(exec.Context())
		switch {
		case err != nil:
			exec.PushProgressUpdate(messages.Update{
				Message: "Getting prices information failed. Plans are displayed without price details",
				Status:  messages.MessageStatusWarning,
				Details: "Error: " + err.Error(),
			})
			// Continue without pricing - just show plans
			showPrices = false
		case pricesByZone != nil:
			// Find the requested zone
			var ok bool
			prices, ok = (*pricesByZone)[s.pricesZone]
			if !ok {
				return nil, fmt.Errorf("pricing zone %s not found", s.pricesZone)
			}
		default:
			// priceZones is nil, disable pricing
			showPrices = false
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
		if showPrices && prices != nil {
			cost := getPlanCost(p, prices, duration)
			row = append(row, cost)
		}

		rows[key] = append(rows[key], row)
	}

	return output.MarshaledWithHumanOutput{
		Value: plans,
		Output: output.Combined{
			planSection("general_purpose", "General purpose", rows["general_purpose"], showPrices, s.pricesDuration),
			planSection("gpu", "GPU", rows["gpu"], showPrices, s.pricesDuration),
			planSection("cloud_native", "Cloud native", rows["cloud_native"], showPrices, s.pricesDuration),
			planSection("high_cpu", "High CPU", rows["high_cpu"], showPrices, s.pricesDuration),
			planSection("high_memory", "High memory", rows["high_memory"], showPrices, s.pricesDuration),
			planSection("developer", "Developer", rows["developer"], showPrices, s.pricesDuration),
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
