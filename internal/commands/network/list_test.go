package network

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	Network1 := upcloud.Network{
		Name: "network-1",
		UUID: "28e15cf5-8817-42ab-b017-970666be96ec",
		Type: upcloud.NetworkTypeUtility,
		Zone: "fi-hel1",
	}

	Network2 := upcloud.Network{
		Name: "network-2",
		UUID: "f9f5ad16-a63a-4670-8449-c01d1e97281e",
		Type: upcloud.NetworkTypePrivate,
		Zone: "fi-hel1",
	}

	Network3 := upcloud.Network{
		Name: "network-3",
		UUID: "e157ce0a-eeb0-49fc-9f2c-a05c3ac57066",
		Type: upcloud.NetworkTypeUtility,
		Zone: "uk-lon1",
	}

	Network4 := upcloud.Network{
		Name: Network1.Name,
		UUID: "b3e49768-f13a-42c3-bea7-4e2471657f2f",
		Type: upcloud.NetworkTypePublic,
		Zone: "uk-lon1",
	}

	networks := &upcloud.Networks{Networks: []upcloud.Network{Network1, Network2, Network3, Network4}}

	for _, test := range []struct {
		name     string
		flags    []string
		expected []upcloud.Network
	}{
		{
			name:     "get all",
			flags:    []string{"--all"},
			expected: []upcloud.Network{Network1, Network2, Network3, Network4},
		},
		{
			name:     "filter where type is utility",
			flags:    []string{"--utility"},
			expected: []upcloud.Network{Network1, Network3},
		},
		{
			name:     "filter where zone is uk-lon1",
			flags:    []string{"--zone", "uk-lon1", "--all"},
			expected: []upcloud.Network{Network3, Network4},
		},
		{
			name:     "filter where zone is uk-lon1 and type is utility",
			flags:    []string{"--zone", "uk-lon1", "--utility"},
			expected: []upcloud.Network{Network3},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On("GetNetworks").Return(networks, nil)
			mService.On("GetNetworksInZone", &request.GetNetworksInZoneRequest{Zone: "fi-hel1"}).Return(&upcloud.Networks{Networks: []upcloud.Network{Network1, Network2}}, nil)
			mService.On("GetNetworksInZone", &request.GetNetworksInZoneRequest{Zone: "uk-lon1"}).Return(&upcloud.Networks{Networks: []upcloud.Network{Network3, Network4}}, nil)

			cfg := config.New()
			c := commands.BuildCommand(ListCommand(), nil, cfg)
			err := c.Cobra().Flags().Parse(test.flags)

			assert.NoError(t, err)

			res, err := c.(commands.NoArgumentCommand).ExecuteWithoutArguments(commands.NewExecutor(cfg, &mService, flume.New("test")))

			assert.Nil(t, err)
			assert.Equal(t, createOutput(test.expected), res)
		})
	}
}

func createOutput(networks []upcloud.Network) output.MarshaledWithHumanOutput {
	rows := []output.TableRow{}
	for _, network := range networks {
		rows = append(rows,
			output.TableRow{network.UUID, network.Name, network.Router, network.Type, network.Zone},
		)
	}

	return output.MarshaledWithHumanOutput{
		Value: upcloud.Networks{Networks: networks},
		Output: output.Table{
			HideHeader: false,
			Columns: []output.TableColumn{
				{Header: "UUID", Key: "uuid", Hidden: false, Colour: ui.DefaultUUUIDColours},
				{Header: "Name", Key: "name", Hidden: false},
				{Header: "Router", Key: "router", Hidden: false, Colour: ui.DefaultUUUIDColours},
				{Header: "Type", Key: "type", Hidden: false},
				{Header: "Zone", Key: "zone", Hidden: false},
			},
			Rows: rows,
		},
	}
}
