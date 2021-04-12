package network

import (
	"bytes"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestShowCommand(t *testing.T) {

	server1 := upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server1.example.com",
		License:      0,
		MemoryAmount: 2048,
		State:        "started",
		Plan:         "1xCPU-2GB",
		Title:        "server1",
		UUID:         "0077fa3d-32db-4b09-9f5f-30d9e9afb568",
		Zone:         "fi-hel1",
		Tags: []string{
			"DEV",
			"Ubuntu",
		},
	}

	server2 := upcloud.Server{
		CoreNumber:   2,
		Hostname:     "server2.example.com",
		License:      0,
		MemoryAmount: 2048,
		State:        "stopped",
		Plan:         "1xCPU-2GB",
		Title:        "server2",
		UUID:         "0077fa3d-32db-4b09-9f5f-30d9e9afb569",
		Zone:         "fi-hel1",
		Tags: []string{
			"DEV",
			"Ubuntu",
		},
	}

	servers := []upcloud.Server{server1, server2}

	network := upcloud.Network{
		IPNetworks: upcloud.IPNetworkSlice{
			{
				Address:          "196.12.0.1",
				DHCP:             upcloud.FromBool(true),
				DHCPDefaultRoute: upcloud.FromBool(true),
				DHCPDns:          []string{"196.12.0.3", "196.12.0.4"},
				Family:           upcloud.IPAddressFamilyIPv4,
				Gateway:          "196.12.0.5",
			},
			{
				Address:          "196.15.0.1",
				DHCP:             upcloud.FromBool(true),
				DHCPDefaultRoute: upcloud.FromBool(false),
				DHCPDns:          []string{"196.15.0.3", "196.15.0.4"},
				Family:           upcloud.IPAddressFamilyIPv4,
				Gateway:          "196.15.0.5",
			},
		},
		Name:   "test-network",
		Type:   "utility",
		UUID:   "ce6a9934-c0c6-4d84-9ad4-0611f5b95e79",
		Zone:   "uk-lon1",
		Router: "79c0ad83-ac84-44f3-a2f8-06cbd524ee8c",
		Servers: []upcloud.NetworkServer{
			{ServerUUID: server1.UUID, ServerTitle: server1.Title},
			{ServerUUID: server2.UUID, ServerTitle: server2.Title},
		},
	}

	expected := `  
  Common
    UUID:   ce6a9934-c0c6-4d84-9ad4-0611f5b95e79 
    Name:   test-network                         
    Router: 79c0ad83-ac84-44f3-a2f8-06cbd524ee8c 
    Type:   utility                              
    Zone:   uk-lon1                              

  IP Networks:

     Address      Family   DHCP   DHCP Def Router   DHCP DNS              
    ──────────── ──────── ────── ───────────────── ───────────────────────
     196.12.0.1   IPv4     yes    yes               196.12.0.3 196.12.0.4 
     196.15.0.1   IPv4     yes    no                196.15.0.3 196.15.0.4 
    
  Servers:

     UUID                                   Title     Hostname              State   
    ────────────────────────────────────── ───────── ───────────────────── ─────────
     0077fa3d-32db-4b09-9f5f-30d9e9afb568   server1   server1.example.com   started 
     0077fa3d-32db-4b09-9f5f-30d9e9afb569   server2   server2.example.com   stopped 
    
`

	cachedNetworks = nil
	mService := smock.Service{}
	mService.On("GetNetworks").Return(&upcloud.Networks{Networks: []upcloud.Network{network}}, nil)
	for _, server := range servers {
		mService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: server.UUID}).Return(&upcloud.ServerDetails{Server: server}, nil)
	}

	conf := config.New()
	// force human output
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(ShowCommand(), nil, conf)

	// get resolver to initialize command cache
	_, err := command.(*showCommand).Get(&mService)
	if err != nil {
		t.Fatal(err)
	}
	res, err := command.(commands.Command).Execute(commands.NewExecutor(conf, &mService), network.UUID)

	assert.Nil(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf, res)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
}
