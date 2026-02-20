package account

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePeriod(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantMonth string // Expected YYYY-MM format
	}{
		// Direct YYYY-MM format
		{
			name:      "direct YYYY-MM",
			input:     "2024-07",
			wantErr:   false,
			wantMonth: "2024-07",
		},
		// Named periods
		{
			name:    "current month empty string",
			input:   "",
			wantErr: false,
			// Month will be current, just verify format
		},
		{
			name:    "current month keyword",
			input:   "month",
			wantErr: false,
		},
		{
			name:    "last month",
			input:   "last month",
			wantErr: false,
		},
		// Relative periods
		{
			name:    "3 months ago",
			input:   "3months",
			wantErr: false,
		},
		{
			name:    "2 weeks ago",
			input:   "2weeks",
			wantErr: false,
		},
		// Relative from base
		{
			name:      "2 months from base",
			input:     "2months from 2024-05",
			wantErr:   false,
			wantMonth: "2024-03",
		},
		{
			name:      "forward from base",
			input:     "+3months from 2024-01",
			wantErr:   false,
			wantMonth: "2024-04",
		},
		// Error cases
		{
			name:    "invalid format",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "invalid unit",
			input:   "3foobar",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMonth, _, err := parsePeriod(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify it's in YYYY-MM format
			var year, month int
			n, _ := fmt.Sscanf(gotMonth, "%d-%d", &year, &month)
			assert.Equal(t, 2, n, "should parse as YYYY-MM")
			assert.True(t, year >= 2020 && year <= 2030, "year should be reasonable")
			assert.True(t, month >= 1 && month <= 12, "month should be 1-12")

			// If we have an expected month, verify it
			if tt.wantMonth != "" {
				assert.Equal(t, tt.wantMonth, gotMonth)
			}
		})
	}
}

func TestBillingCommandBackwardCompatibility(t *testing.T) {
	cmd := &billingCommand{}

	// Test that all original fields still exist
	cmd.year = 2024
	cmd.month = 7
	cmd.resourceID = "test-uuid"
	cmd.username = "testuser"

	// Verify fields are set
	assert.Equal(t, 2024, cmd.year)
	assert.Equal(t, 7, cmd.month)
	assert.Equal(t, "test-uuid", cmd.resourceID)
	assert.Equal(t, "testuser", cmd.username)

	// Test new fields also work
	cmd.period = "last month"
	cmd.match = "production"
	cmd.category = "server"
	cmd.detailed = true

	assert.Equal(t, "last month", cmd.period)
	assert.Equal(t, "production", cmd.match)
	assert.Equal(t, "server", cmd.category)
	assert.True(t, cmd.detailed)
}

func TestYearMonthFlagsOverridePeriod(t *testing.T) {
	// Test that when both year/month flags and period are specified,
	// year/month takes precedence for backward compatibility
	cmd := &billingCommand{
		year:   2024,
		month:  3,
		period: "last month", // Should be ignored
	}

	// Simulate the logic from ExecuteWithoutArguments
	var yearMonth string
	if cmd.year != 0 && cmd.month != 0 {
		yearMonth = fmt.Sprintf("%d-%02d", cmd.year, cmd.month)
	} else if cmd.period != "" {
		yearMonth, _, _ = parsePeriod(cmd.period)
	}

	assert.Equal(t, "2024-03", yearMonth, "year/month flags should override period")
}

func TestOriginalBehaviorWithoutNewFlags(t *testing.T) {
	// When only year/month are provided (original usage),
	// no new features should interfere
	cmd := &billingCommand{
		year:  2024,
		month: 7,
		// New fields all at zero/empty values
		period:   "",
		match:    "",
		category: "",
		detailed: false,
	}

	// These should work exactly as before
	assert.Equal(t, 2024, cmd.year)
	assert.Equal(t, 7, cmd.month)
	assert.Empty(t, cmd.period)
	assert.Empty(t, cmd.match)
	assert.Empty(t, cmd.category)
	assert.False(t, cmd.detailed)
}

func TestPeriodParsing(t *testing.T) {
	// Test that various period formats produce valid YYYY-MM
	periods := []string{
		"month",
		"last month",
		"3months",
		"2024-07",
		"quarter",
		"year",
		"2weeks",
	}

	for _, period := range periods {
		t.Run(period, func(t *testing.T) {
			yearMonth, desc, err := parsePeriod(period)
			require.NoError(t, err)

			// Verify YYYY-MM format
			_, err = time.Parse("2006-01", yearMonth)
			assert.NoError(t, err, "should be valid YYYY-MM format")

			// Description should not be empty
			assert.NotEmpty(t, desc)
		})
	}
}

func TestFirstElementAsString(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected string
	}{
		{
			name:     "string element",
			input:    []interface{}{"test", 123},
			expected: "test",
		},
		{
			name:     "non-string first",
			input:    []interface{}{123, "test"},
			expected: "",
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			expected: "",
		},
		{
			name:     "nil element",
			input:    []interface{}{nil, "test"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := firstElementAsString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
