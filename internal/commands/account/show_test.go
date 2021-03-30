package account

import (
	"bytes"
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) GetAccount() (*upcloud.Account, error) {
	args := m.Called()
	return args[0].(*upcloud.Account), args.Error(1)
}

func TestShowCommand(t *testing.T) {

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
	mService := new(MockAccountService)

	mService.On("GetAccount").Return(&account, nil)
	conf.Service = mService
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	cmd := commands.BuildCommand(testCmd, nil, conf)

	out, err := cmd.(commands.NewCommand).MakeExecutor()([]string{})
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf, out)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())

	mService.AssertNumberOfCalls(t, "GetAccount", 1)
}
