package server

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStopCommand(t *testing.T) {
	targetMethod := "StopServer"

	Server1 := upcloud.Server{
		State: upcloud.ServerStateMaintenance,
		Title: "server-1-title",
		UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
	}

	servers := &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
		},
	}

	details := upcloud.ServerDetails{
		Server: Server1,
	}

	details2 := upcloud.ServerDetails{
		Server: upcloud.Server{
			State: upcloud.ServerStateStopped,
			Title: "server-1-title",
			UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		},
	}

	for _, test := range []struct {
		name    string
		args    []string
		stopReq request.StopServerRequest
	}{
		{
			name: "use default values",
			args: []string{},
			stopReq: request.StopServerRequest{
				UUID:     Server1.UUID,
				StopType: defaultStopType,
			},
		},
		{
			name: "flags mapped to the correct field",
			args: []string{
				"--type", "hard",
			},
			stopReq: request.StopServerRequest{
				UUID:     Server1.UUID,
				StopType: upcloud.StopTypeHard,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := StopCommand()
			mService := new(smock.Service)

			mService.On("GetServers", mock.Anything).Return(servers, nil)
			mService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: Server1.UUID}).Return(&details2, nil)
			stopReq := test.stopReq
			mService.On(targetMethod, &stopReq).Return(&details, nil)

			c := commands.BuildCommand(testCmd, nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(
				commands.NewExecutor(conf, mService, flume.New("test")),
				Server1.UUID,
			)
			assert.NoError(t, err)

			mService.AssertNumberOfCalls(t, targetMethod, 1)
		})
	}
}
