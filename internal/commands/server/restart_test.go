package server

import (
	internal "github.com/UpCloudLtd/cli/internal/service"
	"testing"
	// "time"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRestartCommand(t *testing.T) {
	methodName := "RestartServer"

	var Server1 = upcloud.Server{
		State: "started",
		Title: "server-1-title",
		UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
	}

	var servers = &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
		},
	}

	details := upcloud.ServerDetails{
		Server: Server1,
	}

	details2 := upcloud.ServerDetails{
		Server: upcloud.Server{
			State: "started",
			Title: "server-1-title",
			UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		},
	}

	for _, test := range []struct {
		name       string
		args       []string
		restartReq request.RestartServerRequest
	}{
		{
			name: "use default values",
			args: []string{},
			restartReq: request.RestartServerRequest{
				UUID:          Server1.UUID,
				StopType:      defaultStopType,
				Timeout:       defaultRestartTimeout,
				TimeoutAction: defaultRestartTimeoutAction,
			},
		},
		{
			name: "flags mapped to the correct field",
			args: []string{
				"--stop-type", "hard",
			},
			restartReq: request.RestartServerRequest{
				UUID:          Server1.UUID,
				StopType:      "hard",
				Timeout:       defaultRestartTimeout,
				TimeoutAction: defaultRestartTimeoutAction,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := RestartCommand()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}
			mService.On("GetServers", mock.Anything).Return(servers, nil)
			mService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: Server1.UUID}).Return(&details2, nil)
			mService.On(methodName, &test.restartReq).Return(&details, nil)

			c := commands.BuildCommand(testCmd, nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.Command).Execute(commands.NewExecutor(conf, mService), Server1.UUID)
			assert.NoError(t, err)

			mService.AssertNumberOfCalls(t, methodName, 1)
		})
	}
}
