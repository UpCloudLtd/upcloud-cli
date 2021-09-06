package ipaddress_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/ipaddress"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
)

func TestModifyCommand(t *testing.T) {
	t.Parallel()
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
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mService := smock.Service{}
			mService.On(targetMethod, &test.expected).Return(&ip, nil)
			mService.On("GetIPAddresses").Return(&upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{ip}}, nil)
			conf := config.New()

			c := commands.BuildCommand(ipaddress.ModifyCommand(), nil, conf)
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
