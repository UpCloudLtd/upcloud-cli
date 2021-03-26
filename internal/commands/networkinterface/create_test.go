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

	s := upcloud.Server{UUID: "c4cb35bc-3fb5-4cce-9951-79cab2225417"}
	servers := upcloud.Servers{Servers: []upcloud.Server{s}}
	network := upcloud.Network{UUID: "aa39e313-d908-418a-a959-459699bdc83a", Name: "test-network"}
	networks := upcloud.Networks{Networks: []upcloud.Network{network}}

	for _, test := range []struct {
		name  string
		args  []string
		error string
		req   request.CreateNetworkInterfaceRequest
	}{
		{
			name: "network is missing",
			args: []string{"--type", "public"},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:        s.UUID,
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
			args: []string{
				"--network", network.Name,
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:        s.UUID,
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
			args: []string{
				"--network", network.Name,
				"--family", "IPv6",
			},
			error: "Currently only IPv4 is supported in private networks",
		},
		{
			name: "set ip-family for public network",
			args: []string{
				"--family", "IPv6",
				"--type", "public",
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:        s.UUID,
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
			args: []string{
				"--network", network.Name,
				"--ip-addresses", "1000.40.210.253",
			},
			error: "1000.40.210.253 is an invalid ip address",
		},
		{
			name: "using default values",
			args: []string{
				"--network", network.Name,
				"--ip-addresses", "127.0.0.1",
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:        s.UUID,
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
			args: []string{
				"--network", network.Name,
				"--ip-addresses", "127.0.0.1,127.0.0.2,127.0.0.3/22",
				"--bootable",
				"--source-ip-filtering",
				"--index", "4",
			},
			req: request.CreateNetworkInterfaceRequest{
				ServerUUID:        s.UUID,
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
			mService := smock.MockService{}
			mService.On("GetNetworks").Return(&networks, nil)
			mService.On(targetMethod, &test.req).Return(&upcloud.Interface{}, nil)

			mService.On("GetServers").Return(&servers, nil)
			c := commands.BuildCommand(CreateCommand(&mService, &mService), nil, config.New())
			err := c.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = c.MakeExecuteCommand()([]string{s.UUID})

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}

}
