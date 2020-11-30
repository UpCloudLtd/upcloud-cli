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

func TestAssignCommand(t *testing.T) {
	methodName := "AssignIPAddress"

	ip := upcloud.IPAddress{}

	for _, test := range []struct {
		name     string
		flags    []string
		error    string
		expected request.AssignIPAddressRequest
	}{
		{
			name:  "using default value with zone",
			flags: []string{"--zone", "uk-lon1"},
			expected: request.AssignIPAddressRequest{
				Access:   upcloud.IPAddressAccessPublic,
				Family:   upcloud.IPAddressFamilyIPv4,
				Floating: upcloud.FromBool(false),
				Zone:     "uk-lon1",
			},
		},
		{
			name: "set optional fields",
			flags: []string{
				"--floating",
				"--family", upcloud.IPAddressFamilyIPv6,
				"--access", upcloud.IPAddressAccessPrivate,
				"--zone", "uk-lon1",
				"--mac", "AA-00-04-00-XX-YY",
			},
			expected: request.AssignIPAddressRequest{
				Floating: upcloud.FromBool(true),
				Access:   upcloud.IPAddressAccessPrivate,
				Family:   upcloud.IPAddressFamilyIPv6,
				MAC:      "AA-00-04-00-XX-YY",
				Zone:     "uk-lon1",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cachedIPs = nil
			mips := MockIpAddressService{}
			mips.On(methodName, &test.expected).Return(&ip, nil)

			c := commands.BuildCommand(AssignCommand(&mips), nil, config.New(viper.New()))
			c.SetFlags(test.flags)

			_, err := c.MakeExecuteCommand()([]string{})

			if err != nil {
				assert.Equal(t, test.error, err.Error())
			} else {
				mips.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
