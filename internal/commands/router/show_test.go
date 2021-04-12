package router

import (
	"bytes"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		AttachedNetworks: upcloud.RouterNetworkSlice{
			{NetworkUUID: networks[0].UUID},
		},
		Name: "test-router",
		Type: "normal",
		UUID: "37f5d657-195c-4b5e-ad61-112945ad184b",
	}

	expected := `  
  UUID: 37f5d657-195c-4b5e-ad61-112945ad184b 
  Name: test-router                          
  Type: normal                               

  Networks:

     UUID                                   Name           Router                                 Type      Zone    
    ────────────────────────────────────── ────────────── ────────────────────────────────────── ───────── ─────────
     ce6a9934-c0c6-4d84-9ad4-0611f5b95e79   test-network   79c0ad83-ac84-44f3-a2f8-06cbd524ee8c   utility   uk-lon1 
    
`
	mService := smock.Service{}
	mService.On("GetRouters", mock.Anything).Return(&upcloud.Routers{Routers: []upcloud.Router{router}}, nil)
	mService.On("GetNetworkDetails", mock.Anything).Return(networks[0], nil)

	conf := config.New()
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	c := commands.BuildCommand(ShowCommand(), nil, conf)

	// get resolver to trigger caching
	_, err := c.(resolver.ResolutionProvider).Get(&mService)
	assert.NoError(t, err)

	res, err := c.(commands.Command).Execute(commands.NewExecutor(conf, &mService), router.UUID)
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf, res)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())

}
