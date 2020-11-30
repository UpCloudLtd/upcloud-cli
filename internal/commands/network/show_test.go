package network

import (
	"bytes"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowCommand(t *testing.T) {

	servers := []*upcloud.Server{
		{
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
		},
		{
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
		},
	}

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
		Name:    "test-network",
		Type:    "utility",
		UUID:    "ce6a9934-c0c6-4d84-9ad4-0611f5b95e79",
		Zone:    "uk-lon1",
		Router:  "79c0ad83-ac84-44f3-a2f8-06cbd524ee8c",
		Servers: nil,
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

	buf := new(bytes.Buffer)
	command := ShowCommand(&MockNetworkService{}, &MockServerService{})
	err := command.HandleOutput(buf, networkWithServers{
		network: &network,
		servers: servers,
	})

	assert.Nil(t, err)
	assert.Equal(t, expected, buf.String())
}
