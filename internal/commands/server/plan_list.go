package server

import (
	"fmt"
	"reflect"
	"sort"
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
	showCost bool
	duration time.Duration
}

// InitCommand initializes the command flags
func (s *planListCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.BoolVar(&s.showCost, "cost", false, "Show cost information for plans")
	flagSet.DurationVar(&s.duration, "duration", time.Hour, "Duration for cost calculation (e.g., 1h, 24h, 168h for week, 720h for month)")

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

	// Fetch pricing information if requested
	var priceZones *upcloud.PriceZones
	if s.showCost {
		priceZones, err = exec.All().GetPriceZones(exec.Context())
		if err != nil {
			// Continue without pricing - just show plans
			s.showCost = false
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
		if s.showCost && priceZones != nil {
			cost := getPlanCost(p, priceZones, s.duration)
			row = append(row, fmt.Sprintf("%.4f", cost))
		}

		rows[key] = append(rows[key], row)
	}

	return output.MarshaledWithHumanOutput{
		Value: plans,
		Output: output.Combined{
			planSection("general_purpose", "General purpose", rows["general_purpose"], s.showCost, s.duration),
			planSection("gpu", "GPU", rows["gpu"], s.showCost, s.duration),
			planSection("cloud_native", "Cloud native", rows["cloud_native"], s.showCost, s.duration),
			planSection("high_cpu", "High CPU", rows["high_cpu"], s.showCost, s.duration),
			planSection("high_memory", "High memory", rows["high_memory"], s.showCost, s.duration),
			planSection("developer", "Developer", rows["developer"], s.showCost, s.duration),
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

func planSection(key, title string, rows []output.TableRow, showCost bool, duration time.Duration) output.CombinedSection {
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

	if showCost {
		columns = append(columns, output.TableColumn{
			Key:    "cost",
			Header: formatDurationHeader(duration),
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
func getPlanCost(plan upcloud.Plan, priceZones *upcloud.PriceZones, duration time.Duration) float64 {
	if priceZones == nil || len(priceZones.PriceZones) == 0 {
		return 0
	}

	// Get the first price zone
	priceZone := priceZones.PriceZones[0]

	// Try to find specific plan pricing first using reflection
	// Field naming convention varies, e.g., "1xCPU-1GB" → "ServerPlan1xCPU1GB"
	fieldName := "ServerPlan" + strings.ReplaceAll(plan.Name, "-", "")
	fieldName = strings.ReplaceAll(fieldName, ".", "")

	v := reflect.ValueOf(priceZone)
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
	return hourlyPrice * duration.Hours()
}

// formatDurationHeader creates a human-readable header for the cost column
func formatDurationHeader(duration time.Duration) string {
	hours := duration.Hours()

	// Common duration labels
	switch {
	case hours == 1:
		return "Cost (per hour)"
	case hours == 24:
		return "Cost (per day)"
	case hours == 24*7:
		return "Cost (per week)"
	case hours >= 24*30 && hours <= 24*31:
		return "Cost (per month)"
	case hours >= 24*365 && hours <= 24*366:
		return "Cost (per year)"
	case hours < 1:
		minutes := int(duration.Minutes())
		return fmt.Sprintf("Cost (per %d min)", minutes)
	case hours < 24:
		if hours == float64(int(hours)) {
			return fmt.Sprintf("Cost (per %d hours)", int(hours))
		}
		return fmt.Sprintf("Cost (per %.1f hours)", hours)
	case hours < 24*7:
		days := hours / 24
		if days == float64(int(days)) {
			return fmt.Sprintf("Cost (per %d days)", int(days))
		}
		return fmt.Sprintf("Cost (per %.1f days)", days)
	default:
		days := hours / 24
		return fmt.Sprintf("Cost (per %.0f days)", days)
	}
}
