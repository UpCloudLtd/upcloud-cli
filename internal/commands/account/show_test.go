package account

import (
	"bytes"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockAccountService struct{}

func (m *MockAccountService) GetAccount() (*upcloud.Account, error) {
	return &upcloud.Account{
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
	}, nil
}

func TestShowCommand(t *testing.T) {

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

	cfg := config.New(viper.New())
	cfg.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(ShowCommand(&MockAccountService{}), nil, cfg)
	out, err := command.(commands.NewCommand).Execute(commands.NewExecutor(cfg), []string{})
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, cfg, out)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
}
