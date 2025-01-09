package testutils

import (
	"github.com/stretchr/testify/require"
	"testing"

	"time"
)

func MustParseRFC3339(t *testing.T, timeStr string) *time.Time {
	t.Helper()

	p, err := time.Parse(time.RFC3339, timeStr)
	require.NoError(t, err)

	p = p.UTC()

	return &p
}
