package account

import (
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestShowCommand(t *testing.T) {
	text.DisableColors()
	account := upcloud.Account{
		Credits:  12345,
		UserName: "upctl_test",
		ResourceLimits: upcloud.ResourceLimits{
			Cores:                 100,
			DetachedFloatingIps:   10,
			ManagedObjectStorages: 20,
			Memory:                307200,
			NetworkPeerings:       100,
			Networks:              100,
			NTPExcessGiB:          20000,
			PublicIPv4:            0,
			PublicIPv6:            100,
			StorageHDD:            10240,
			StorageMaxIOPS:        10240,
			StorageSSD:            10240,
			LoadBalancers:         50,
			GPUs:                  10,
		},
	}

	expected := `  
  Username: upctl_test 
  Credits:  €123.45    
  
  Resource Limits:
    Cores:                      100 
    Detached Floating IPs:       10 
    Load balancers:              50 
    Managed object storages:     20 
    Memory:                  307200 
    Network peerings:           100 
    Networks:                   100 
    NTP excess GiB:           20000 
    Public IPv4:                  0 
    Public IPv6:                100 
    Storage HDD:              10240 
    Storage MaxIOPS:          10240 
    Storage SSD:              10240 
    GPUs:                        10 

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
