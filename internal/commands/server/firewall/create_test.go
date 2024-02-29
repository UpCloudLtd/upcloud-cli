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

func TestCreateFirewallRuleCommand(t *testing.T) {
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

	for _, test := range []struct {
		name        string
		flags       []string
		arg         string
		expectedReq *request.CreateFirewallRuleRequest
		error       string
	}{
		{
			name:  "Empty info",
			flags: []string{},
			arg:   Server1.UUID,
			error: `required flag(s) "direction", "action" not set`,
		},
		{
			name: "Action is required",
			flags: []string{
				Server1.UUID,
				"--direction", "in",
			},
			arg:   Server1.UUID,
			error: `required flag(s) "action" not set`,
		},
		{
			name: "FirewallRule, drop incoming by default",
			flags: []string{
				Server1.UUID,
				"--direction", "in",
				"--action", "drop",
			},
			arg: Server1.UUID,
			expectedReq: &request.CreateFirewallRuleRequest{
				FirewallRule: upcloud.FirewallRule{
					Direction: "in",
					Action:    "drop",
				},
				ServerUUID: Server1.UUID,
			},
		},
		{
			name: "FirewallRule, accept incoming IPv6",
			flags: []string{
				Server1.UUID,
				"--direction", "in",
				"--action", "accept",
				"--family", "IPv6",
			},
			arg: Server1.UUID,
			expectedReq: &request.CreateFirewallRuleRequest{
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
			mService := smock.Service{}
			if test.expectedReq != nil {
				mService.On("CreateFirewallRule", test.expectedReq).Return(&upcloud.FirewallRule{}, nil)
			} else {
				mService.On("CreateFirewallRule", mock.Anything).Return(&upcloud.FirewallRule{}, nil)
			}

			conf := config.New()
			cc := commands.BuildCommand(CreateCommand(), nil, conf)

			cc.Cobra().SetArgs(append(test.flags, test.arg))
			_, err := mockexecute.MockExecute(cc, &mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, "CreateFirewallRule", 1)
			}
		})
	}
}
