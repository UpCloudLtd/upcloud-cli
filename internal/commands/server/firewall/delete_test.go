package serverfirewall

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteServerFirewallRuleCommand(t *testing.T) {
	Server1 := upcloud.Server{
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

	// Mock firewall rules
	mockRules := &upcloud.FirewallRules{
		FirewallRules: []upcloud.FirewallRule{
			{
				Position:                1,
				Direction:               "in",
				Action:                  "accept",
				Protocol:                "tcp",
				DestinationPortStart:    "22",
				DestinationPortEnd:      "22",
				DestinationAddressStart: "0.0.0.0",
				DestinationAddressEnd:   "255.255.255.255",
				SourceAddressStart:      "0.0.0.0",
				SourceAddressEnd:        "255.255.255.255",
				Comment:                 "SSH",
			},
			{
				Position:                2,
				Direction:               "in",
				Action:                  "accept",
				Protocol:                "tcp",
				DestinationPortStart:    "80",
				DestinationPortEnd:      "80",
				DestinationAddressStart: "0.0.0.0",
				DestinationAddressEnd:   "255.255.255.255",
				SourceAddressStart:      "0.0.0.0",
				SourceAddressEnd:        "255.255.255.255",
				Comment:                 "HTTP",
			},
		},
	}

	for _, test := range []struct {
		name  string
		flags []string
		error string
	}{
		{
			name:  "no filters",
			flags: []string{},
			error: `would delete 2 firewall rules (exceeds skip-confirmation=1)`,
		},
		{
			name:  "position 1",
			flags: []string{"--position", "1"},
		},
		{
			name:  "invalid position",
			flags: []string{"--position", "-1"},
			error: "no firewall rules matched the specified filters",
		},
		{
			name:  "position too high",
			flags: []string{"--position", "1001"},
			error: "no firewall rules matched the specified filters",
		},
	} {
		deleteRuleMethod := "DeleteFirewallRule"
		getFirewallRulesMethod := "GetFirewallRules"
		t.Run(test.name, func(t *testing.T) {
			mService := new(smock.Service)

			// Mock GetFirewallRules call
			mService.On(getFirewallRulesMethod, mock.Anything).Return(mockRules, nil).Maybe()

			// Mock DeleteFirewallRule call
			mService.On(deleteRuleMethod, mock.Anything).Return(nil, nil).Maybe()

			conf := config.New()
			cc := commands.BuildCommand(DeleteCommand(), nil, conf)

			cc.Cobra().SetArgs(append(test.flags, Server1.UUID))
			_, err := mockexecute.MockExecute(cc, mService, conf)

			if test.error != "" {
				fmt.Println("ERROR", test.error, "==", err)
				assert.ErrorContains(t, err, test.error)
				mService.AssertNumberOfCalls(t, deleteRuleMethod, 0)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, getFirewallRulesMethod, 1)
				mService.AssertNumberOfCalls(t, deleteRuleMethod, 1)
			}
		})
	}
}
