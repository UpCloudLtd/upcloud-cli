package ipaddress

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/gemalto/flume"
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
		arg      string
		flags    []string
		error    string
		expected request.ModifyIPAddressRequest
	}{
		{
			name: "set optional fields, ip identified by address",
			arg:  ip.Address,
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
			mService := smock.Service{}
			expected := test.expected
			mService.On(targetMethod, &expected).Return(&ip, nil)
			mService.On("GetIPAddresses").Return(&upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{ip}}, nil)
			conf := config.New()

			c := commands.BuildCommand(ModifyCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, &mService, flume.New("test")), test.arg)
			if err != nil {
				assert.Equal(t, test.error, err.Error())
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
