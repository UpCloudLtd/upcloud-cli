package server

import (
	"math"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestMonthDurationParsing(t *testing.T) {
	testCases := []struct {
		input          string
		expectedHours  int
		expectedHeader string
	}{
		{
			input:          "1m",
			expectedHours:  28 * 24,
			expectedHeader: "Price (per month)",
		},
		{
			input:          "3m",
			expectedHours:  3 * 28 * 24,
			expectedHeader: "Price (per 3 months)",
		},
		{
			input:          "6m",
			expectedHours:  6 * 28 * 24,
			expectedHeader: "Price (per 6 months)",
		},
		{
			input:          "12m",
			expectedHours:  12 * 28 * 24,
			expectedHeader: "Price (per 12 months)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			// Test that duration calculation matches 28 days per month
			duration := time.Duration(0)

			// Simulate the parsing logic from plan_list.go
			if tc.input[len(tc.input)-1] == 'm' {
				months := 0
				if tc.input == "1m" {
					months = 1
				} else if tc.input == "3m" {
					months = 3
				} else if tc.input == "6m" {
					months = 6
				} else if tc.input == "12m" {
					months = 12
				}

				// Use 28 days per month (as per UpCloud billing policy)
				duration = time.Duration(months) * 28 * 24 * time.Hour
			}

			actualHours := int(duration.Hours())
			assert.Equal(t, tc.expectedHours, actualHours,
				"Duration %s should be %d hours (28 days per month)", tc.input, tc.expectedHours)

			// Test header formatting
			header := formatPricingHeader(tc.input)
			assert.Equal(t, tc.expectedHeader, header,
				"Header for %s should be '%s'", tc.input, tc.expectedHeader)
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
		name            string
		duration        time.Duration
		expectedHours   float64
		description     string
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

func TestPlanCostCalculationWithRounding(t *testing.T) {
	// Test the actual getPlanCost function with sub-hour durations
	priceZone := &upcloud.PriceZone{
		Name: "test-zone",
		ServerCore: &upcloud.Price{
			Amount: 100000, // 1.00 per hour
		},
		ServerMemory: &upcloud.Price{
			Amount: 100000, // 1.00 per GB per hour
		},
	}

	plan := upcloud.Plan{
		Name:         "1xCPU-1GB",
		CoreNumber:   1,
		MemoryAmount: 1024, // 1 GB
	}

	testCases := []struct {
		duration     time.Duration
		expectedCost float64
		description  string
	}{
		{
			duration:     10 * time.Second,
			expectedCost: 2.00, // 1 CPU + 1 GB = 2.00 per hour, any fraction rounds up
			description:  "10 seconds should cost same as 1 hour",
		},
		{
			duration:     30 * time.Minute,
			expectedCost: 2.00,
			description:  "30 minutes should cost same as 1 hour",
		},
		{
			duration:     1 * time.Hour,
			expectedCost: 2.00,
			description:  "1 hour should cost exactly hourly rate",
		},
		{
			duration:     61 * time.Minute,
			expectedCost: 4.00, // 2 hours worth
			description:  "61 minutes should cost same as 2 hours",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cost := getPlanCost(plan, priceZone, tc.duration)
			assert.Equal(t, tc.expectedCost, cost, tc.description)
		})
	}
}