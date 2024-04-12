package networkinterface

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

func TestCreateCommand(t *testing.T) {
	targetMethod := "CreateNetworkInterface"

	server := upcloud.Server{UUID: "c4cb35bc-3fb5-4cce-9951-79cab2225417"}
	network := upcloud.Network{UUID: "aa39e313-d908-418a-a959-459699bdc83a", Name: "test-network"}
	networks := upcloud.Networks{Networks: []upcloud.Network{network}}

	for _, test := range []struct {
		name  string
		flags []string
		error string
		req   request.CreateNetworkInterfaceRequest
	}{
		{
			name:  "network is missing",
			flags: []string{"--type", "public"},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID: server.UUID,
				IPAddresses: request.CreateNetworkInterfaceIPAddressSlice{
					{Family: upcloud.IPAddressFamilyIPv4},
				},
				Type: upcloud.NetworkTypePublic,
			},
		},
		{
			name: "ip-address is missing",
			flags: []string{
				"--network", network.UUID,
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:  server.UUID,
				NetworkUUID: network.UUID,
				IPAddresses: request.CreateNetworkInterfaceIPAddressSlice{
					{Family: upcloud.IPAddressFamilyIPv4},
				},
				Type: upcloud.NetworkTypePrivate,
			},
		},
		{
			name: "set ip-family for public network",
			flags: []string{
				"--family", "IPv6",
				"--type", "public",
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID: server.UUID,
				IPAddresses: request.CreateNetworkInterfaceIPAddressSlice{
					{Family: upcloud.IPAddressFamilyIPv6},
				},
				Type: upcloud.NetworkTypePublic,
			},
		},
		{
			name: "invalid ip-address",
			flags: []string{
				"--network", network.UUID,
				"--ip-addresses", "1000.40.210.253",
			},
			error: "1000.40.210.253 is an invalid ip address",
		},
		{
			name: "using default values",
			flags: []string{
				"--network", network.UUID,
				"--ip-addresses", "127.0.0.1",
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:  server.UUID,
				NetworkUUID: network.UUID,
				IPAddresses: request.CreateNetworkInterfaceIPAddressSlice{
					{Address: "127.0.0.1", Family: upcloud.IPAddressFamilyIPv4},
				},
				Type: upcloud.NetworkTypePrivate,
			},
		},
		{
			name: "using network name",
			flags: []string{
				"--network", network.Name,
				"--ip-addresses", "127.0.0.1",
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:  server.UUID,
				NetworkUUID: network.UUID,
				IPAddresses: request.CreateNetworkInterfaceIPAddressSlice{
					{Address: "127.0.0.1", Family: upcloud.IPAddressFamilyIPv4},
				},
				Type: upcloud.NetworkTypePrivate,
			},
		},
		{
			name: "set optional fields",
			flags: []string{
				"--network", network.UUID,
				"--ip-addresses", "127.0.0.1,127.0.0.2,127.0.0.3/22",
				"--enable-bootable",
				"--disable-source-ip-filtering",
				"--index", "4",
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:        server.UUID,
				Bootable:          upcloud.True,
				SourceIPFiltering: upcloud.False,
				NetworkUUID:       network.UUID,
				IPAddresses: request.CreateNetworkInterfaceIPAddressSlice{
					{Address: "127.0.0.1", Family: upcloud.IPAddressFamilyIPv4},
					{Address: "127.0.0.2", Family: upcloud.IPAddressFamilyIPv4},
					{Address: "127.0.0.3/22", Family: upcloud.IPAddressFamilyIPv4},
				},
				Type:  upcloud.NetworkTypePrivate,
				Index: 4,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On("GetNetworks").Return(&networks, nil)
			req := test.req
			mService.On(targetMethod, &req).Return(&upcloud.Interface{}, nil)

			mService.On("GetNetworkDetails", &request.GetNetworkDetailsRequest{UUID: network.UUID}).Return(&network, nil)
			conf := config.New()

			c := commands.BuildCommand(CreateCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = c.(commands.SingleArgumentCommand).ExecuteSingleArgument(commands.NewExecutor(conf, &mService, flume.New("test")), server.UUID)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
