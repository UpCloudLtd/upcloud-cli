package server

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
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

func TestPlanListMachineOutputIncludesPrices(t *testing.T) {
	plans := upcloud.Plans{Plans: []upcloud.Plan{
		{Name: "1xCPU-1GB", CoreNumber: 1, MemoryAmount: 1024, StorageSize: 25, StorageTier: "maxiops", PublicTrafficOut: 1024},
		{Name: "2xCPU-4GB", CoreNumber: 2, MemoryAmount: 4096, StorageSize: 80, StorageTier: "maxiops", PublicTrafficOut: 4096},
	}}
	prices := upcloud.PricesByZone{
		"de-fra1": {
			"server_plan_1xCPU-1GB": {Amount: 1, Price: 1},
		},
	}

	for _, test := range []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name: "without prices",
			args: []string{},
			expected: `[
				{"core_number": 1, "memory_amount": 1024, "name": "1xCPU-1GB", "public_traffic_out": 1024, "storage_size": 25, "storage_tier": "maxiops", "gpu_amount": 0, "gpu_model": ""},
				{"core_number": 2, "memory_amount": 4096, "name": "2xCPU-4GB", "public_traffic_out": 4096, "storage_size": 80, "storage_tier": "maxiops", "gpu_amount": 0, "gpu_model": ""}
			]`,
		},
		{
			name: "with hourly prices",
			args: []string{"--prices", "de-fra1", "--prices-duration", "hour"},
			expected: `[
				{"core_number": 1, "memory_amount": 1024, "name": "1xCPU-1GB", "public_traffic_out": 1024, "storage_size": 25, "storage_tier": "maxiops", "gpu_amount": 0, "gpu_model": "", "cost": 0.01},
				{"core_number": 2, "memory_amount": 4096, "name": "2xCPU-4GB", "public_traffic_out": 4096, "storage_size": 80, "storage_tier": "maxiops", "gpu_amount": 0, "gpu_model": "", "cost": null}
			]`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			conf.Viper().Set(config.KeyOutput, config.ValueOutputJSON)

			mService := new(smock.Service)
			mService.On("GetPlans").Return(&plans, nil)
			mService.On("GetPricesByZone").Return(&prices, nil)

			c := commands.BuildCommand(PlanListCommand(), nil, conf)
			c.Cobra().SetArgs(test.args)

			output, err := mockexecute.MockExecute(c, mService, conf)
			assert.NoError(t, err)
			assert.JSONEq(t, test.expected, output)
		})
	}
}
