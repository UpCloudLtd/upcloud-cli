package serverfirewall

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	for _, test := range []struct {
		name        string
		flags       []string
		arg         string
		expectedReq *request.CreateFirewallRuleRequest
		error       string
	}{
		{
			name:  "Empty arg",
			flags: []string{},
			error: "server is required",
		},
		{
			name:  "Empty info",
			flags: []string{},
			arg:   Server1.UUID,
			error: "direction is required",
		},
		{
			name: "Action is required",
			flags: []string{
				Server1.UUID,
				"--direction", "in",
			},
			arg:   Server1.UUID,
			error: "action is required",
		},
		{
			name: "Family is required",
			flags: []string{
				Server1.UUID,
				"--direction", "in",
				"--action", "accept",
			},
			arg:   Server1.UUID,
			error: "family (IPv4/IPv6) is required",
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
			err := cc.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = cc.(commands.Command).Execute(commands.NewExecutor(conf, &mService), test.arg)
			if test.error != "" {
				assert.Error(t, err)
				assert.Equal(t, test.error, err.Error())
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, "CreateFirewallRule", 1)
			}
		})
	}
}
