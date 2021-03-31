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

func TestModifyCommand(t *testing.T) {
	targetMethod := "ModifyIPAddress"

	ip := upcloud.IPAddress{
		Address:   "127.0.0.1",
		PTRRecord: "old.ptr.com",
	}

	for _, test := range []struct {
		name     string
		args     []string
		flags    []string
		error    string
		expected request.ModifyIPAddressRequest
	}{
		{
			name: "set optional fields, ip identified by address",
			args: []string{ip.Address},
			flags: []string{
				"--ptr-record", "example.com",
				"--mac", "AA-00-04-00-XX-YY",
			},
			expected: request.ModifyIPAddressRequest{
				IPAddress: ip.Address,
				PTRRecord: "example.com",
				MAC:       "AA-00-04-00-XX-YY",
			},
		},
		{
			name: "set optional fields, ip identified by PTR Record",
			args: []string{ip.PTRRecord},
			flags: []string{
				"--ptr-record", "example.com",
				"--mac", "AA-00-04-00-XX-YY",
			},
			expected: request.ModifyIPAddressRequest{
				IPAddress: ip.Address,
				PTRRecord: "example.com",
				MAC:       "AA-00-04-00-XX-YY",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cachedIPs = nil
			mService := smock.Service{}
			mService.On(targetMethod, &test.expected).Return(&ip, nil)
			mService.On("GetIPAddresses").Return(&upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{ip}}, nil)

			c := commands.BuildCommand(ModifyCommand(&mService), nil, config.New())
			err := c.SetFlags(test.flags)
			assert.NoError(t, err)

			_, err = c.MakeExecuteCommand()(test.args)

			if err != nil {
				assert.Equal(t, test.error, err.Error())
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
