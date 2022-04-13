package account

import (
	"bytes"
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/gemalto/flume"
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
	out, err := command.(commands.NoArgumentCommand).ExecuteWithoutArguments(commands.NewExecutor(conf, mService, flume.New("test")))
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf, out)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())

	mService.AssertNumberOfCalls(t, "GetAccount", 1)
}
