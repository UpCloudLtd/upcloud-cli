package databaseproperties

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDatabasePropertiesShow(t *testing.T) {
	text.DisableColors()

	mService := smock.Service{}
	mService.On("GetManagedDatabaseServiceType", mock.Anything).Return(&propertiesTestdata, nil)

	conf := config.New()
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(ShowCommand("pg", "PostgreSQL"), nil, conf)

	command.Cobra().SetArgs([]string{"timescaledb.max_background_workers"})
	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.NoError(t, err)

	expected := `  
  Key:         timescaledb.max_background_workers                                                                                                                                                                                                       
  Title:       timescaledb.max_background_workers                                                                                                                                                                                                       
  Description: The number of background workers for timescaledb operations. You should configure this setting to the sum of your number of databases and the total number of concurrent background workers you want running at any given point in time. 
  Create only: no                                                                                                                                                                                                                                       
  Type:        integer                                                                                                                                                                                                                                  
  Default:     16                                                                                                                                                                                                                                       
  Minimum:     1                                                                                                                                                                                                                                        
  Maximum:     4096                                                                                                                                                                                                                                     

`
	assert.Equal(t, expected, output)
}
