package database

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	for _, test := range []struct {
		name     string
		expected string
		db       *upcloud.ManagedDatabase
	}{
		{
			name:     "nil database",
			expected: "",
			db:       nil,
		}, {
			name:     "nil metadata",
			expected: "",
			db:       &upcloud.ManagedDatabase{Metadata: nil},
		}, {
			name:     "pg",
			expected: "15",
			db: &upcloud.ManagedDatabase{
				Type:     upcloud.ManagedDatabaseServiceTypePostgreSQL,
				Metadata: &upcloud.ManagedDatabaseMetadata{PGVersion: "15"},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			version := getVersion(test.db)
			assert.Equal(t, test.expected, version)
		})
	}
}
