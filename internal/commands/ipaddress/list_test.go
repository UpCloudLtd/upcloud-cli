package ipaddress

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

func TestList_MarshaledOutput(t *testing.T) {
	ipAddress := upcloud.IPAddress{
		Address:    "94.237.117.150",
		Access:     "public",
		Family:     "IPv4",
		PartOfPlan: upcloud.FromBool(true),
		PTRRecord:  "94-237-117-150.fi-hel1.upcloud.host",
		ServerUUID: "005ab220-7ff6-42c9-8615-e4c02eb4104e",
		MAC:        "ee:1b:db:ca:6b:80",
		Floating:   upcloud.FromBool(false),
		Zone:       "fi-hel1",
	}
	ipAddresses := upcloud.IPAddresses{
		IPAddresses: []upcloud.IPAddress{
			ipAddress,
		},
	}

	mService := smock.Service{}
	mService.On("GetIPAddresses").Return(&ipAddresses, nil)

	conf := config.New()
	conf.Viper().Set(config.KeyOutput, config.ValueOutputJSON)

	command := commands.BuildCommand(ListCommand(), nil, conf)

	output, err := mockexecute.MockExecute(command, &mService, conf)

	assert.Nil(t, err)
	testutils.AssertOutputHasType(t, output, &upcloud.IPAddressSlice{})
}
