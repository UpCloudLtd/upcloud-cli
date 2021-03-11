package serverfirewall_test

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/commands/serverfirewall"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateFirewallRuleCommand(t *testing.T) {

	var Server1 = upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server-1-hostname",
		License:      0,
		MemoryAmount: 1024,
		Plan:         "server-1-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-1-title",
		UUID:         "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		Zone:         "fi-hel1",
	}

	var servers = &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
		},
	}

	for _, test := range []struct {
		name        string
		args        []string
		expectedReq request.CreateFirewallRuleRequest
		error       string
	}{
		{
			name:  "Empty info",
			args:  []string{},
			error: "Direction is required.",
		},
		{
			name: "FirewallRule, accept incoming IPv6",
			args: []string{
				Server1.UUID,
				"--direction", "in",
				"--action", "accept",
				"--family", "IPv6",
			},
			expectedReq: request.CreateFirewallRuleRequest{
				FirewallRule: upcloud.FirewallRule{
					Direction: "in",
					Action:    "accept",
					Family:    "IPv6",
				},
				ServerUUID: Server1.UUID,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			server.CachedServers = nil

			mServerService := server.MockServerService{}
			mServerService.On("GetServers", mock.Anything).Return(servers, nil)

			mFirewallRuleService := MockFirewallRuleService{}

			cc := commands.BuildCommand(serverfirewall.CreateCommand(&mServerService, &mFirewallRuleService), nil, config.New(viper.New()))
			cc.SetFlags(test.args)

			_, err := cc.MakeExecuteCommand()([]string{Server1.UUID})

			if test.error != "" {
				assert.Equal(t, test.error, err.Error())
			}
		})
	}
}
