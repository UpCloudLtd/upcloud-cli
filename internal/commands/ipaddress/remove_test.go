package ipaddress

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
)

func TestRemoveCommand(t *testing.T) {
	ip := upcloud.IPAddress{
		Address:   "127.0.0.1",
		PTRRecord: "old.ptr.com",
	}

	for _, test := range []struct {
		name     string
		arg      string
		error    string
		expected request.ReleaseIPAddressRequest
	}{
		{
			name: "remove ip identified by address",
			arg:  ip.Address,
			expected: request.ReleaseIPAddressRequest{
				IPAddress: ip.Address,
			},
		},
	} {
		targetMethod := "ReleaseIPAddress"
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On(targetMethod, &test.expected).Return(nil)
			mService.On("GetIPAddresses").Return(&upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{ip}}, nil)
			conf := config.New()

			c := commands.BuildCommand(RemoveCommand(), nil, conf)
			_, err := c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, &mService, flume.New("test")), test.arg)

			if err != nil {
				assert.Equal(t, test.error, err.Error())
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
