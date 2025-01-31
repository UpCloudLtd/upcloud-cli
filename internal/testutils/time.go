package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func MustParseRFC3339(t *testing.T, timeStr string) *time.Time {
	t.Helper()

	p, err := time.Parse(time.RFC3339, timeStr)
	require.NoError(t, err)

	p = p.UTC()

	return &p
}
