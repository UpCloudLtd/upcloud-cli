package billing

import (
	"fmt"
	"strings"
	"time"
)

// parsePeriod parses various period formats into YYYY-MM format for API
func parsePeriod(period string) (string, string, error) {
	now := time.Now()

	// Handle YYYY-MM format directly
	if matched, _ := fmt.Sscanf(period, "%d-%d", new(int), new(int)); matched == 2 {
		return period, period, nil
	}

	// Handle named periods
	switch period {
	case "month", "current", "":
		yearMonth := now.Format("2006-01")
		return yearMonth, fmt.Sprintf("current month (%s)", yearMonth), nil
	case "day", "today":
		yearMonth := now.Format("2006-01")
		return yearMonth, fmt.Sprintf("today (%s)", now.Format("2006-01-02")), nil
	case "quarter":
		// Get current quarter
		quarter := (now.Month()-1)/3 + 1
		yearMonth := now.Format("2006-01")
		return yearMonth, fmt.Sprintf("Q%d %d (current month: %s)", quarter, now.Year(), yearMonth), nil
	case "year":
		yearMonth := now.Format("2006-01")
		return yearMonth, fmt.Sprintf("year %d (current month: %s)", now.Year(), yearMonth), nil
	}

	// Handle relative periods from a base date (e.g., "2months from 2024-06", "+3months from 2024-01")
	if strings.Contains(period, " from ") {
		parts := strings.Split(period, " from ")
		if len(parts) == 2 {
			relPeriod := parts[0]
			baseDate := parts[1]

			// Parse base date
			var baseTime time.Time
			if matched, _ := fmt.Sscanf(baseDate, "%d-%d", new(int), new(int)); matched == 2 {
				// Parse as YYYY-MM
				baseTime, _ = time.Parse("2006-01", baseDate)
			} else {
				return "", "", fmt.Errorf("invalid base date format: %s (use YYYY-MM)", baseDate)
			}

			// Parse relative period
			var amount int
			var unit string
			forward := false

			// Check for + prefix (forward in time)
			if strings.HasPrefix(relPeriod, "+") {
				forward = true
				relPeriod = strings.TrimPrefix(relPeriod, "+")
			}

			if matched, _ := fmt.Sscanf(relPeriod, "%d%s", &amount, &unit); matched == 2 {
				var targetTime time.Time
				if forward {
					switch unit {
					case "day", "days":
						targetTime = baseTime.AddDate(0, 0, amount)
					case "week", "weeks":
						targetTime = baseTime.AddDate(0, 0, amount*7)
					case "month", "months":
						targetTime = baseTime.AddDate(0, amount, 0)
					case "year", "years":
						targetTime = baseTime.AddDate(amount, 0, 0)
					default:
						return "", "", fmt.Errorf("unknown period unit: %s", unit)
					}
				} else {
					switch unit {
					case "day", "days":
						targetTime = baseTime.AddDate(0, 0, -amount)
					case "week", "weeks":
						targetTime = baseTime.AddDate(0, 0, -amount*7)
					case "month", "months":
						targetTime = baseTime.AddDate(0, -amount, 0)
					case "year", "years":
						targetTime = baseTime.AddDate(-amount, 0, 0)
					default:
						return "", "", fmt.Errorf("unknown period unit: %s", unit)
					}
				}
				yearMonth := targetTime.Format("2006-01")
				direction := "before"
				if forward {
					direction = "after"
				}
				return yearMonth, fmt.Sprintf("%s %s %s (%s)", relPeriod, direction, baseDate, yearMonth), nil
			}
		}
	}

	// Handle simple relative periods from now (e.g., "3days", "2months", "2weeks")
	var amount int
	var unit string
	if matched, _ := fmt.Sscanf(period, "%d%s", &amount, &unit); matched == 2 {
		var targetTime time.Time
		switch unit {
		case "day", "days":
			targetTime = now.AddDate(0, 0, -amount)
		case "week", "weeks":
			targetTime = now.AddDate(0, 0, -amount*7)
		case "month", "months":
			targetTime = now.AddDate(0, -amount, 0)
		case "year", "years":
			targetTime = now.AddDate(-amount, 0, 0)
		default:
			return "", "", fmt.Errorf("unknown period unit: %s", unit)
		}
		yearMonth := targetTime.Format("2006-01")
		return yearMonth, fmt.Sprintf("%s ago (%s)", period, yearMonth), nil
	}

	// Handle "last" periods (with spaces and without)
	if strings.HasPrefix(period, "last") {
		// Handle with space: "last month", "last quarter", "last year"
		parts := strings.Fields(period)
		unit := ""
		if len(parts) == 2 {
			unit = parts[1]
		} else {
			// Handle without space: "lastmonth", "lastquarter", "lastyear"
			unit = strings.TrimPrefix(period, "last")
		}

		switch unit {
		case "month":
			targetTime := now.AddDate(0, -1, 0)
			yearMonth := targetTime.Format("2006-01")
			return yearMonth, fmt.Sprintf("last month (%s)", yearMonth), nil
		case "quarter":
			// Go back 3 months for last quarter
			targetTime := now.AddDate(0, -3, 0)
			yearMonth := targetTime.Format("2006-01")
			quarter := (targetTime.Month()-1)/3 + 1
			return yearMonth, fmt.Sprintf("last quarter Q%d %d (showing %s)", quarter, targetTime.Year(), yearMonth), nil
		case "year":
			// Go back 12 months for last year
			targetTime := now.AddDate(-1, 0, 0)
			yearMonth := targetTime.Format("2006-01")
			return yearMonth, fmt.Sprintf("last year %d (showing %s)", targetTime.Year(), yearMonth), nil
		}
	}

	return "", "", fmt.Errorf("unknown period format: %s. Use formats like 'month', 'day', 'quarter', 'year', '3days', '2months', 'last month', or 'YYYY-MM'", period)
}
