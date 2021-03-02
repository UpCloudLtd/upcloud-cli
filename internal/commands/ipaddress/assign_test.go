package ipaddress

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssignCommand(t *testing.T) {
	methodName := "AssignIPAddress"

	s1 := upcloud.Server{UUID: "f2e42635-b8b8-48ed-aa44-a494ef438f83", Title: "s1"}
	s2 := upcloud.Server{UUID: "7bc4c854-a87d-40c0-97b5-b0d17333248d", Title: "s2"}
	s3 := upcloud.Server{UUID: "1d2c4cfd-b835-4814-9d81-2904b74ad86d", Title: "s3"}

	servers := upcloud.Servers{Servers: []upcloud.Server{s1, s2, s3}}

	ip := upcloud.IPAddress{}

	for _, test := range []struct {
		name     string
		flags    []string
		error    string
		expected request.AssignIPAddressRequest
	}{
		{
			name:  "using default value",
			flags: []string{"--zone", "uk-lon1"},
			error: "server is required for non-floating IP",
		},
		{
			name:  "using default value with server",
			flags: []string{"--server", s2.UUID},
			expected: request.AssignIPAddressRequest{
				Access:     upcloud.IPAddressAccessPublic,
				Family:     upcloud.IPAddressFamilyIPv4,
				Floating:   upcloud.FromBool(false),
				ServerUUID: s2.UUID,
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
			mips := MockIPAddressService{}
			mips.On(methodName, &test.expected).Return(&ip, nil)
			mss := server.MockServerService{}
			mss.On("GetServers").Return(&servers, nil)

			c := commands.BuildCommand(AssignCommand(&mss, &mips), nil, config.New(viper.New()))
			err := c.SetFlags(test.flags)
			assert.NoError(t, err)

			_, err = c.MakeExecuteCommand()([]string{})

			if err != nil {
				assert.Equal(t, test.error, err.Error())
			} else {
				mips.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
