package databaseproperties

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var propertiesTestdata = upcloud.ManagedDatabaseType{
	Properties: map[string]upcloud.ManagedDatabaseServiceProperty{
		"timescaledb": {
			Title:       "TimescaleDB extension configuration values",
			Description: "System-wide settings for the timescaledb extension",
			Type:        "object",
			Properties: map[string]upcloud.ManagedDatabaseServiceProperty{
				"max_background_workers": {
					Default:     16.0,
					Example:     8.0,
					Title:       "timescaledb.max_background_workers",
					Type:        "integer",
					Description: "The number of background workers for timescaledb operations. You should configure this setting to the sum of your number of databases and the total number of concurrent background workers you want running at any given point in time.",
					Minimum:     upcloud.Float64Ptr(1),
					Maximum:     upcloud.Float64Ptr(4096),
				},
			},
		},
	},
}

func TestDatabasePropertiesByType(t *testing.T) {
	text.DisableColors()

	mService := smock.Service{}
	mService.On("GetManagedDatabaseServiceType", mock.Anything).Return(&propertiesTestdata, nil)

	conf := config.New()
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(DBTypeCommand("pg", "PostgreSQL"), nil, conf)

	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)

	expected := `
 Property                             Create only   Type      Example 
──────────────────────────────────── ───────────── ───────── ─────────
 timescaledb                          no            object            
 timescaledb.max_background_workers   no            integer   8       

`
	assert.Equal(t, expected, output)
}
