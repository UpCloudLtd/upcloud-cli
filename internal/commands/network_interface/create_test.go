package network_interface

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

func TestCreateCommand(t *testing.T) {
	methodName := "CreateNetworkInterface"

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
			error: "ip-address is required",
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
			server.CachedServers = nil
			mns := MockNetworkService{}
			mns.On("GetNetworks").Return(&networks, nil)
			mns.On(methodName, &test.req).Return(&upcloud.Interface{}, nil)

			mss := MockServerService{}
			mss.On("GetServers").Return(&servers, nil)
			c := commands.BuildCommand(CreateCommand(&mss, &mns), nil, config.New(viper.New()))
			c.SetFlags(test.args)

			_, err := c.MakeExecuteCommand()([]string{s.UUID})

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mns.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}

}
