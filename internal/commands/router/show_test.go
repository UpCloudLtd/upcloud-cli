package router

import (
	"bytes"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowCommand(t *testing.T) {

	networks := []*upcloud.Network{
		{
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
		},
	}

	router := upcloud.Router{
		AttachedNetworks: nil,
		Name:             "test-router",
		Type:             "normal",
		UUID:             "37f5d657-195c-4b5e-ad61-112945ad184b",
	}

	expected := `  
  Common
    UUID: 37f5d657-195c-4b5e-ad61-112945ad184b 
    Name: test-router                          
    Type: normal                               
  
  Networks:
     UUID                                   Name           Router                                 Type      Zone    
    ────────────────────────────────────── ────────────── ────────────────────────────────────── ───────── ─────────
     ce6a9934-c0c6-4d84-9ad4-0611f5b95e79   test-network   79c0ad83-ac84-44f3-a2f8-06cbd524ee8c   utility   uk-lon1 
`

	buf := new(bytes.Buffer)
	command := ShowCommand(&MockRouterService{}, &MockNetworkService{})
	err := command.HandleOutput(buf, &routerWithNetworks{
		router:   &router,
		networks: networks,
	})

	assert.Nil(t, err)
	assert.Equal(t, expected, buf.String())
}
