package account

import (
	"bytes"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockAccountService struct{}

func (m *MockAccountService) GetAccount() (*upcloud.Account, error) {
	return nil, nil
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

	buf := new(bytes.Buffer)
	command := ShowCommand(&MockAccountService{})
	err := command.HandleOutput(buf, &account)

	assert.Nil(t, err)
	assert.Equal(t, expected, buf.String())
}
