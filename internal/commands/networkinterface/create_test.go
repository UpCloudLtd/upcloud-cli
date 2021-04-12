package networkinterface

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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
				ServerUUID:        server.UUID,
				Bootable:          upcloud.FromBool(false),
				SourceIPFiltering: upcloud.FromBool(false),
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
				ServerUUID:        server.UUID,
				Bootable:          upcloud.FromBool(false),
				SourceIPFiltering: upcloud.FromBool(false),
				NetworkUUID:       network.UUID,
				IPAddresses: request.CreateNetworkInterfaceIPAddressSlice{
					{Family: upcloud.IPAddressFamilyIPv4},
				},
				Type: upcloud.NetworkTypePrivate,
			},
		},
		{
			name: "ip-family unsupported for private network",
			flags: []string{
				"--network", network.UUID,
				"--family", "IPv6",
			},
			error: "Currently only IPv4 is supported in private networks",
		},
		{
			name: "set ip-family for public network",
			flags: []string{
				"--family", "IPv6",
				"--type", "public",
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:        server.UUID,
				Bootable:          upcloud.FromBool(false),
				SourceIPFiltering: upcloud.FromBool(false),
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
				ServerUUID:        server.UUID,
				Bootable:          upcloud.FromBool(false),
				SourceIPFiltering: upcloud.FromBool(false),
				NetworkUUID:       network.UUID,
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
				"--bootable",
				"--source-ip-filtering",
				"--index", "4",
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:        server.UUID,
				Bootable:          upcloud.FromBool(true),
				SourceIPFiltering: upcloud.FromBool(true),
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
			mService.On(targetMethod, &test.req).Return(&upcloud.Interface{}, nil)

			mService.On("GetNetworkDetails", &request.GetNetworkDetailsRequest{UUID: network.UUID}).Return(&network, nil)
			conf := config.New()

			c := commands.BuildCommand(CreateCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = c.(commands.Command).Execute(commands.NewExecutor(conf, &mService), server.UUID)

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}
