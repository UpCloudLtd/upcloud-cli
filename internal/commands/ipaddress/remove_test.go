package ipaddress

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestRemoveCommand(t *testing.T) {
	targetMethod := "ReleaseIPAddress"

	ip := upcloud.IPAddress{
		Address:   "127.0.0.1",
		PTRRecord: "old.ptr.com",
	}

	for _, test := range []struct {
		name     string
		args     []string
		error    string
		expected request.ReleaseIPAddressRequest
	}{
		{
			name: "set optional fields, ip identified by address",
			args: []string{ip.Address},
			expected: request.ReleaseIPAddressRequest{
				IPAddress: ip.Address,
			},
		},
		{
			name: "set optional fields, ip identified by PTR Record",
			args: []string{ip.PTRRecord},
			expected: request.ReleaseIPAddressRequest{
				IPAddress: ip.Address,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cachedIPs = nil
			mService := smock.Service{}
			mService.On(targetMethod, &test.expected).Return(nil)
			mService.On("GetIPAddresses").Return(&upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{ip}}, nil)

			c := commands.BuildCommand(RemoveCommand(&mService), nil, config.New())

			_, err := c.MakeExecuteCommand()(test.args)

			if err != nil {
				assert.Equal(t, test.error, err.Error())
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
