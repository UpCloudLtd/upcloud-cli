package serverfirewall

import (
	"fmt"
	"testing"

	"github.com/gemalto/flume"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteServerFirewallRuleCommand(t *testing.T) {
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
		name  string
		flags []string
		error string
	}{
		{
			name:  "no position",
			flags: []string{},
			error: "position is required",
		},
		{
			name:  "position 1",
			flags: []string{"--position", "1"},
		},
		{
			name:  "invalid position",
			flags: []string{"--position", "-1"},
			error: "invalid position (1-1000 allowed)",
		},
	} {
		deleteRuleMethod := "DeleteFirewallRule"
		t.Run(test.name, func(t *testing.T) {
			mService := new(smock.Service)
			mService.On(deleteRuleMethod, mock.Anything).Return(nil, nil)

			conf := config.New()
			cc := commands.BuildCommand(DeleteCommand(), nil, conf)
			err := cc.Cobra().Flags().Parse(test.flags)
			assert.NoError(t, err)

			_, err = cc.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), Server1.UUID)
			if test.error != "" {
				fmt.Println("ERROR", test.error, "==", err)
				assert.EqualError(t, err, test.error)
				mService.AssertNumberOfCalls(t, deleteRuleMethod, 0)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, deleteRuleMethod, 1)
			}
		})
	}
}
