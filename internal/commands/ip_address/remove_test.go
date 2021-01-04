package ip_address

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveCommand(t *testing.T) {
	methodName := "ReleaseIPAddress"

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
			mips := MockIpAddressService{}
			mips.On(methodName, &test.expected).Return(nil)
			mips.On("GetIPAddresses").Return(&upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{ip}}, nil)

			c := commands.BuildCommand(RemoveCommand(&mips), nil, config.New(viper.New()))

			_, err := c.MakeExecuteCommand()(test.args)

			if err != nil {
				assert.Equal(t, test.error, err.Error())
			} else {
				mips.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
