package server

import (
	"testing"

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
