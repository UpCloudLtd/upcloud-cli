package server

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/spf13/pflag"
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
	flagSet.StringVar(&s.pricingDuration, "pricing-duration", "1m", "Duration for pricing calculation (e.g., 1h, 24h, 1m, 3m, 12m)")

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
	if showPricing {
		// Check if it's a month-based duration (e.g., 1m, 3m, 12m)
		if strings.HasSuffix(s.pricingDuration, "m") && !strings.Contains(s.pricingDuration, "h") && !strings.Contains(s.pricingDuration, "s") {
			// Try to parse as months
			monthStr := strings.TrimSuffix(s.pricingDuration, "m")
			months, err := strconv.Atoi(monthStr)
			if err == nil && months > 0 {
				// Use 28 days per month for pricing calculations (UpCloud bills max 28 days per month)
				duration = time.Duration(months) * 28 * 24 * time.Hour
			} else {
				// Not a valid month duration, try as standard duration
				duration, err = time.ParseDuration(s.pricingDuration)
				if err != nil {
					return nil, fmt.Errorf("invalid pricing-duration: %s (use formats like '1h', '24h', '1m', '3m', '12m')", s.pricingDuration)
				}
			}
		} else {
			// Parse as standard duration
			var err error
			duration, err = time.ParseDuration(s.pricingDuration)
			if err != nil {
				return nil, fmt.Errorf("invalid pricing-duration: %s (use formats like '1h', '24h', '1m', '3m', '12m')", s.pricingDuration)
			}
		}
	}

	// Fetch pricing information if requested
	var priceZone *upcloud.PriceZone
	if showPricing {
		priceZones, err := exec.All().GetPriceZones(exec.Context())
		if err != nil {
			// Continue without pricing - just show plans
			showPricing = false
		} else if priceZones != nil {
			// Find the requested zone
			for _, zone := range priceZones.PriceZones {
				if zone.Name == s.pricingZone {
					priceZone = &zone
					break
				}
			}
			if priceZone == nil {
				return nil, fmt.Errorf("pricing zone %s not found", s.pricingZone)
			}
		} else {
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
		if showPricing && priceZone != nil {
			cost := getPlanCost(p, priceZone, duration)
			row = append(row, fmt.Sprintf("%.4f", cost))
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
		columns = append(columns, output.TableColumn{
			Key:    "cost",
			Header: formatPricingHeader(pricingDuration),
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
func getPlanCost(plan upcloud.Plan, priceZone *upcloud.PriceZone, duration time.Duration) float64 {
	if priceZone == nil {
		return 0
	}

	// Try to find specific plan pricing first using reflection
	// Field naming convention varies, e.g., "1xCPU-1GB" → "ServerPlan1xCPU1GB"
	fieldName := "ServerPlan" + strings.ReplaceAll(plan.Name, "-", "")
	fieldName = strings.ReplaceAll(fieldName, ".", "")

	v := reflect.ValueOf(*priceZone)
	field := v.FieldByName(fieldName)

	var hourlyPrice float64

	if field.IsValid() && !field.IsNil() {
		// Found specific plan pricing
		priceField := field.Elem().FieldByName("Amount")
		if priceField.IsValid() && priceField.Kind() == reflect.Int {
			// Amount is in 1/100000 of currency unit
			hourlyPrice = float64(priceField.Int()) / 100000.0
		}
	}

	// If no specific plan price, calculate from components
	if hourlyPrice == 0 && priceZone.ServerCore != nil && priceZone.ServerMemory != nil {
		corePrice := float64(priceZone.ServerCore.Amount) / 100000.0 // Amount is in 1/100000 of currency unit
		memPrice := float64(priceZone.ServerMemory.Amount) / 100000.0

		// Price per core * number of cores + price per GB * memory in GB
		hourlyPrice = (corePrice * float64(plan.CoreNumber)) + (memPrice * float64(plan.MemoryAmount) / 1024.0)
	}

	// Calculate cost for the requested duration
	// UpCloud bills per (starting) hour, so round up to next full hour
	return hourlyPrice * math.Ceil(duration.Hours())
}

// formatPricingHeader creates a human-readable header for the cost column
func formatPricingHeader(pricingDuration string) string {
	// Check if it's a month-based duration
	if strings.HasSuffix(pricingDuration, "m") && !strings.Contains(pricingDuration, "h") && !strings.Contains(pricingDuration, "s") {
		monthStr := strings.TrimSuffix(pricingDuration, "m")
		months, err := strconv.Atoi(monthStr)
		if err == nil && months > 0 {
			if months == 1 {
				return "Price (per month)"
			}
			return fmt.Sprintf("Price (per %d months)", months)
		}
	}

	// Handle hour-based durations
	switch pricingDuration {
	case "1h":
		return "Price (per hour)"
	case "24h":
		return "Price (per day)"
	default:
		// For other durations, just display the duration string
		return fmt.Sprintf("Price (%s)", pricingDuration)
	}
}
