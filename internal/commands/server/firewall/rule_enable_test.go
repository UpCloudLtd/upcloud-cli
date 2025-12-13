package serverfirewall

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRuleEnableCommand(t *testing.T) {
	serverUUID := "1fdfda29-ead1-4855-b71f-1e33eb2ca9de"

	currentRules := &upcloud.FirewallRules{
		FirewallRules: []upcloud.FirewallRule{
			{
				Position:  1,
				Direction: "in",
				Action:    "drop",
				Protocol:  "tcp",
			},
			{
				Position:  2,
				Direction: "in",
				Action:    "accept",
				Protocol:  "tcp",
			},
			{
				Position:  5,
				Direction: "out",
				Action:    "drop",
				Protocol:  "udp",
			},
		},
	}

	for _, test := range []struct {
		name        string
		flags       []string
		arg         string
		expectedErr string
		checkMocks  func(*testing.T, *smock.Service)
	}{
		{
			name:        "Missing position flag",
			flags:       []string{},
			arg:         serverUUID,
			expectedErr: `required flag(s) "position" not set`,
		},
		{
			name:        "Invalid position - too low",
			flags:       []string{"--position", "0"},
			arg:         serverUUID,
			expectedErr: "invalid position (1-1000 allowed)",
		},
		{
			name:        "Invalid position - too high",
			flags:       []string{"--position", "1001"},
			arg:         serverUUID,
			expectedErr: "invalid position (1-1000 allowed)",
		},
		{
			name:  "Successfully enable rule at position 5",
			flags: []string{"--position", "5"},
			arg:   serverUUID,
			checkMocks: func(t *testing.T, mService *smock.Service) {
				mService.AssertCalled(t, "GetFirewallRules", &request.GetFirewallRulesRequest{
					ServerUUID: serverUUID,
				})

				mService.AssertCalled(t, "CreateFirewallRules", mock.MatchedBy(func(req interface{}) bool {
					r, ok := req.(*request.CreateFirewallRulesRequest)
					if !ok {
						return false
					}
					if r.ServerUUID != serverUUID {
						return false
					}
					for _, rule := range r.FirewallRules {
						if rule.Position == 5 {
							return rule.Action == upcloud.FirewallRuleActionAccept
						}
					}
					return false
				}))
			},
		},
		{
			name:        "Rule position not found",
			flags:       []string{"--position", "99"},
			arg:         serverUUID,
			expectedErr: "firewall rule at position 99 not found on server",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}

			mService.On("GetFirewallRules", mock.MatchedBy(func(req interface{}) bool {
				r, ok := req.(*request.GetFirewallRulesRequest)
				return ok && r.ServerUUID == serverUUID
			})).Return(currentRules, nil)

			mService.On("CreateFirewallRules", mock.Anything).Return(nil)

			conf := config.New()
			cc := commands.BuildCommand(RuleEnableCommand(), nil, conf)

			cc.Cobra().SetArgs(append(test.flags, test.arg))
			_, err := mockexecute.MockExecute(cc, &mService, conf)

			if test.expectedErr != "" {
				assert.ErrorContains(t, err, test.expectedErr)
			} else {
				assert.NoError(t, err)
				if test.checkMocks != nil {
					test.checkMocks(t, &mService)
				}
			}
		})
	}
}
