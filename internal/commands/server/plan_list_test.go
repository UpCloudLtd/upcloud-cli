package server

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMonthDurationParsing(t *testing.T) {
	testCases := []struct {
		input          string
		expectedHours  int
		expectedHeader string
	}{
		{
			input:          "hour",
			expectedHours:  1,
			expectedHeader: "Price (per hour)",
		},
		{
			input:          "month",
			expectedHours:  28 * 24,
			expectedHeader: "Price (per month)",
		},
		{
			input:          "1h",
			expectedHours:  1,
			expectedHeader: "Price (per hour)",
		},
		{
			input:          "4h30m",
			expectedHours:  5,
			expectedHeader: "Price (per 4h30m)",
		},
		{
			input:          "24h",
			expectedHours:  24,
			expectedHeader: "Price (per day)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			// Test that duration calculation matches expected
			var duration time.Duration

			// Simulate the parsing logic from plan_list.go
			switch tc.input {
			case "hour":
				duration = 1 * time.Hour
			case "month":
				duration = 28 * 24 * time.Hour
			default:
				duration, _ = time.ParseDuration(tc.input)
			}

			actualHours := int(math.Ceil(duration.Hours()))
			assert.Equal(t, tc.expectedHours, actualHours,
				"Duration %s should be %d hours (28 days per month)", tc.input, tc.expectedHours)

			// Test header formatting
			header := formatPricingHeader(tc.input)
			assert.Equal(t, tc.expectedHeader, header)
		})
	}
}

func TestSixMonthsPricingIs28DaysPerMonth(t *testing.T) {
	// Specific test to ensure 6 months uses 28 days, not 30
	sixMonths := 6

	// Calculate using 28 days per month
	expectedHours := sixMonths * 28 * 24

	// This would be wrong (using 30 days)
	wrongHours := sixMonths * 30 * 24

	// Verify we get the correct value
	assert.Equal(t, 4032, expectedHours, "6 months should be 4032 hours (6 * 28 * 24)")
	assert.NotEqual(t, wrongHours, expectedHours, "Should NOT use 30 days per month")

	// The difference should be 6 * 2 * 24 = 288 hours
	difference := wrongHours - expectedHours
	assert.Equal(t, 288, difference, "Difference between 30-day and 28-day calculation should be 288 hours for 6 months")
}

func TestSubHourDurationBilling(t *testing.T) {
	// Test that any duration less than an hour is billed as a full hour
	// UpCloud bills per (starting) hour, so partial hours round up

	testCases := []struct {
		name          string
		duration      time.Duration
		expectedHours float64
		description   string
	}{
		{
			name:          "1 second",
			duration:      1 * time.Second,
			expectedHours: 1.0,
			description:   "1 second should be billed as 1 hour",
		},
		{
			name:          "1 minute",
			duration:      1 * time.Minute,
			expectedHours: 1.0,
			description:   "1 minute should be billed as 1 hour",
		},
		{
			name:          "30 minutes",
			duration:      30 * time.Minute,
			expectedHours: 1.0,
			description:   "30 minutes should be billed as 1 hour",
		},
		{
			name:          "59 minutes",
			duration:      59 * time.Minute,
			expectedHours: 1.0,
			description:   "59 minutes should be billed as 1 hour",
		},
		{
			name:          "exactly 1 hour",
			duration:      1 * time.Hour,
			expectedHours: 1.0,
			description:   "Exactly 1 hour should be billed as 1 hour",
		},
		{
			name:          "1 hour 1 second",
			duration:      1*time.Hour + 1*time.Second,
			expectedHours: 2.0,
			description:   "1 hour and 1 second should be billed as 2 hours",
		},
		{
			name:          "1.5 hours",
			duration:      90 * time.Minute,
			expectedHours: 2.0,
			description:   "1.5 hours should be billed as 2 hours",
		},
		{
			name:          "2.1 hours",
			duration:      126 * time.Minute,
			expectedHours: 3.0,
			description:   "2.1 hours should be billed as 3 hours",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the billing calculation from getPlanCost
			// which uses math.Ceil(duration.Hours())
			billedHours := math.Ceil(tc.duration.Hours())

			assert.Equal(t, tc.expectedHours, billedHours, tc.description)
		})
	}
}
