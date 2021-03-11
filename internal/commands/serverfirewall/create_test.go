package serverfirewall

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

func TestAttachStorageCommand(t *testing.T) {
	methodName := "CreateFirewallRule"

	var Rule1 = upcloud.FirewallRule{
		Action:               upcloud.FirewallRule.FirewallRuleActionAccept,
		Comment:              "Allow HTTP from anywhere",
		DestinationPortStart: "80",
		DestinationPortEnd:   "80",
		Direction:            upcloud.FirewallRule.FirewallRuleDirectionIn,
		Family:               upcloud.FirewallRule.IPAddressFamilyIPv4,
		Position:             1,
	}

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

	var serverDetails = upcloud.ServerDetails{
		Server: upcloud.Server{
			UUID:  UUID1,
			State: upcloud.ServerStateStarted,
		},
		VideoModel: "vga",
		Firewall:   "off",
	}

	for _, test := range []struct {
		name       string
		args       []string
		createruleReq request.CreateFirewallRuleRequest
		error      string
	}{
		{
			name:  "Empty info",
			args:  []string{},
			error: "Info is required",
		},
		{
			name: "FirewallRule, accept incoming IPv6",
			args: []string{
				Server1.UUID,
				"--direction", "in",
				"--action", "accept",
				"--family", "IPv6",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			server.CachedServers = nil

			mServerService := server.MockServerService{}
			mServerService.On("GetServers", mock.Anything).Return(servers, nil)

			mFirewallRuleService := MockFirewallRuleService{}

			cc := commands.BuildCommand(CreateCommand(&mServerService, &mFirewallRuleService), nil, config.New(viper.New()))
			cc.SetFlags(test.args)

			_, err := cc.MakeExecuteCommand()([]string{Server1.UUID})

			if test.error != "" {
				assert.Equal(t, test.error, err.Error())
			} else {
				mFirewallRuleService.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
