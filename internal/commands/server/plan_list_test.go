package server

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestDurationHeader(t *testing.T) {
	testCases := []struct {
		input          string
		expectedHeader string
	}{
		{
			input:          "hour",
			expectedHeader: "Price (per hour)",
		},
		{
			input:          "month",
			expectedHeader: "Price (per month)",
		},
		{
			input:          "1h",
			expectedHeader: "Price (per hour)",
		},
		{
			input:          "4h30m",
			expectedHeader: "Price (per 4h30m)",
		},
		{
			input:          "24h",
			expectedHeader: "Price (per day)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			header := formatPricingHeader(tc.input)
			assert.Equal(t, tc.expectedHeader, header)
		})
	}
}

func TestGetPlanCost(t *testing.T) {
	prices := map[string]upcloud.Price{
		"server_plan_1xCPU-1GB": {
			Amount: 1,
			Price:  1.0416,
		},
	}

	testcases := []struct {
		plan     upcloud.Plan
		duration string
		expected float64
	}{
		{
			plan:     upcloud.Plan{Name: "1xCPU-1GB"},
			duration: "1h",
			expected: 0.010416,
		},
		{
			plan:     upcloud.Plan{Name: "1xCPU-1GB"},
			duration: "2h30m",
			expected: 0.031248,
		},
		{
			plan:     upcloud.Plan{Name: "1xCPU-1GB"},
			duration: "hour",
			expected: 0.010416,
		},
		{
			plan:     upcloud.Plan{Name: "1xCPU-1GB"},
			duration: "month",
			expected: 7,
		},
	}

	for _, tc := range testcases {
		name := tc.plan.Name + "-" + tc.duration
		t.Run(name, func(t *testing.T) {
			duration, err := getDuration(tc.duration)
			assert.NoError(t, err)

			cost := getPlanCost(tc.plan, prices, duration)
			assert.InDelta(t, tc.expected, cost, 0.001)
		})
	}
}
