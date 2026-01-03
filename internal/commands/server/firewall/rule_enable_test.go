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
				Position:                1,
				Direction:               "in",
				Action:                  "accept",
				Protocol:                "tcp",
				DestinationAddressStart: "0.0.0.0",
				DestinationAddressEnd:   "255.255.255.255",
				SourceAddressStart:      "0.0.0.0",
				SourceAddressEnd:        "255.255.255.255",
			},
			{
				Position:                2,
				Direction:               "in",
				Action:                  "drop",
				Protocol:                "tcp",
				DestinationAddressStart: "0.0.0.0",
				DestinationAddressEnd:   "255.255.255.255",
				SourceAddressStart:      "0.0.0.0",
				SourceAddressEnd:        "255.255.255.255",
				Comment:                 "Catch-all drop rule",
			},
			{
				Position:                3,
				Direction:               "in",
				Action:                  "accept",
				Protocol:                "tcp",
				DestinationAddressStart: "0.0.0.0",
				DestinationAddressEnd:   "255.255.255.255",
				SourceAddressStart:      "0.0.0.0",
				SourceAddressEnd:        "255.255.255.255",
				Comment:                 "Test rule",
			},
			{
				Position:                4,
				Direction:               "out",
				Action:                  "accept",
				Protocol:                "udp",
				DestinationAddressStart: "0.0.0.0",
				DestinationAddressEnd:   "255.255.255.255",
				SourceAddressStart:      "0.0.0.0",
				SourceAddressEnd:        "255.255.255.255",
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
			expectedErr: "would enable 2 firewall rules (exceeds skip-confirmation=1)",
		},
		{
			name:        "Invalid position - too low",
			flags:       []string{"--position", "0"},
			arg:         serverUUID,
			expectedErr: "would enable 2 firewall rules (exceeds skip-confirmation=1)",
		},
		{
			name:        "Invalid position - too high",
			flags:       []string{"--position", "1001"},
			arg:         serverUUID,
			expectedErr: "no disabled firewall rules matched the specified filters",
		},
		{
			name:  "Successfully enable rule at position 3",
			flags: []string{"--position", "3"},
			arg:   serverUUID,
			checkMocks: func(t *testing.T, mService *smock.Service) {
				mService.AssertCalled(t, "GetFirewallRules", &request.GetFirewallRulesRequest{
					ServerUUID: serverUUID,
				})

				// Should delete the old rule at position 3 (after catch-all)
				mService.AssertCalled(t, "DeleteFirewallRule", &request.DeleteFirewallRuleRequest{
					ServerUUID: serverUUID,
					Position:   3,
				})

				// Should create the rule before catch-all (position 2)
				mService.AssertCalled(t, "CreateFirewallRule", mock.MatchedBy(func(req interface{}) bool {
					r, ok := req.(*request.CreateFirewallRuleRequest)
					if !ok {
						return false
					}
					return r.ServerUUID == serverUUID && r.FirewallRule.Position == 2
				}))
			},
		},
		{
			name:        "Rule position not found",
			flags:       []string{"--position", "99"},
			arg:         serverUUID,
			expectedErr: "no disabled firewall rules matched the specified filters",
		},
		{
			name:        "Skip confirmation set to 0 requires confirmation for single rule",
			flags:       []string{"--comment", "Test", "--skip-confirmation", "0"},
			arg:         serverUUID,
			expectedErr: "would enable 1 firewall rules (exceeds skip-confirmation=0)",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}

			mService.On("GetFirewallRules", mock.MatchedBy(func(req interface{}) bool {
				r, ok := req.(*request.GetFirewallRulesRequest)
				return ok && r.ServerUUID == serverUUID
			})).Return(currentRules, nil)

			mService.On("DeleteFirewallRule", mock.Anything).Return(nil).Maybe()
			mService.On("CreateFirewallRule", mock.Anything).Return(&upcloud.FirewallRule{}, nil).Maybe()

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
