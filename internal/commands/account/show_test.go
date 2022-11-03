package account

import (
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestShowCommand(t *testing.T) {
	text.DisableColors()
	account := upcloud.Account{
		Credits:  42,
		UserName: "opencredo",
		ResourceLimits: upcloud.ResourceLimits{
			Cores:               100,
			DetachedFloatingIps: 0,
			Memory:              307200,
			Networks:            100,
			PublicIPv4:          0,
			PublicIPv6:          100,
			StorageHDD:          10240,
			StorageSSD:          10240,
		},
	}

	expected := `  
  Username: opencredo 
  Credits:  0.42$     
  
  Resource Limits:
    Cores:                    100 
    Detached Floating IPs:      0 
    Memory:                307200 
    Networks:                 100 
    Public IPv4:                0 
    Public IPv6:              100 
    Storage HDD:            10240 
    Storage SSD:            10240 

`

	conf := config.New()
	testCmd := ShowCommand()
	mService := new(smock.Service)

	mService.On("GetAccount").Return(&account, nil)
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(testCmd, nil, conf)
	output, err := mockexecute.MockExecute(command, mService, conf)

	assert.NoError(t, err)
	assert.Equal(t, expected, output)
	mService.AssertNumberOfCalls(t, "GetAccount", 1)
}
