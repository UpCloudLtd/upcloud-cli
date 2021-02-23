package networkinterface

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

func TestModifyCommand(t *testing.T) {
	methodName := "ModifyNetworkInterface"

	s := upcloud.Server{UUID: "c4cb35bc-3fb5-4cce-9951-79cab2225417"}
	servers := upcloud.Servers{Servers: []upcloud.Server{s}}
	network := upcloud.Network{UUID: "aa39e313-d908-418a-a959-459699bdc83a", Name: "test-network"}
	networks := upcloud.Networks{Networks: []upcloud.Network{network}}

	for _, test := range []struct {
		name  string
		args  []string
		error string
		req   request.ModifyNetworkInterfaceRequest
	}{
		{
			name:  "index is missing",
			args:  []string{},
			error: "index is required",
		},
		{
			name: "index is present, using default values",
			args: []string{
				"--index", "4",
			},
			req: request.ModifyNetworkInterfaceRequest{CurrentIndex: 4, ServerUUID: s.UUID},
		},
		{
			name: "index is present, all values modified",
			args: []string{
				"--index", "4",
				"--new-index", "5",
				"--bootable", "false",
				"--source-ip-filtering", "true",
				"--ip-addresses", "127.0.0.2,127.0.0.3,127.0.0.4",
			},
			req: request.ModifyNetworkInterfaceRequest{
				ServerUUID:        s.UUID,
				CurrentIndex:      4,
				NewIndex:          5,
				Bootable:          upcloud.FromBool(false),
				SourceIPFiltering: upcloud.FromBool(true),
				IPAddresses: request.CreateNetworkInterfaceIPAddressSlice{
					{Address: "127.0.0.2", Family: upcloud.IPAddressFamilyIPv4},
					{Address: "127.0.0.3", Family: upcloud.IPAddressFamilyIPv4},
					{Address: "127.0.0.4", Family: upcloud.IPAddressFamilyIPv4},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			server.CachedServers = nil

			mns := MockNetworkService{}
			mns.On("GetNetworks").Return(&networks, nil)
			mns.On(methodName, &test.req).Return(&upcloud.Interface{}, nil)

			mss := server.MockServerService{}
			mss.On("GetServers").Return(&servers, nil)
			c := commands.BuildCommand(ModifyCommand(&mns, &mss), nil, config.New(viper.New()))
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
