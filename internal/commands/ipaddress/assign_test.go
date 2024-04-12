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

func TestAssignCommand(t *testing.T) {
	targetMethod := "AssignIPAddress"

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
			mService := smock.Service{}
			expected := test.expected
			mService.On(targetMethod, &expected).Return(&ip, nil)
			mService.On("GetServers").Return(&servers, nil)
			for _, server := range servers.Servers {
				mService.On("GetServerDetails",
					&request.GetServerDetailsRequest{UUID: server.UUID},
				).Return(&upcloud.ServerDetails{Server: server}, nil)
			}
			conf := config.New()

			c := commands.BuildCommand(AssignCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = c.(commands.NoArgumentCommand).ExecuteWithoutArguments(commands.NewExecutor(conf, &mService, flume.New("test")))

			if err != nil {
				assert.Equal(t, test.error, err.Error())
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
