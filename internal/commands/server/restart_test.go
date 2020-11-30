package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
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

	dur120, _ := time.ParseDuration("120s")
	dur10, _ := time.ParseDuration("10s")

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
				StopType:      "soft",
				Timeout:       dur120,
				TimeoutAction: "ignore",
				Host:          0,
			},
		},
		{
			name: "flags mapped to the correct field",
			args: []string{
				"--stop-type", "hard",
				"--timeout-action", "destroy",
				"--timeout", "10",
				"--host", "1234",
			},
			restartReq: request.RestartServerRequest{
				UUID:          Server1.UUID,
				StopType:      "hard",
				Timeout:       dur10,
				TimeoutAction: "destroy",
				Host:          1234,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedServers = nil

			mServerService := MockServerService{}
			mServerService.On("GetServers", mock.Anything).Return(servers, nil)
			mServerService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: Server1.UUID}).Return(&details2, nil)
			mServerService.On(methodName, &test.restartReq).Return(&details, nil)

			c := commands.BuildCommand(RestartCommand(&mServerService), nil, config.New(viper.New()))
			c.SetFlags(test.args)

			c.MakeExecuteCommand()([]string{Server1.UUID})

			mServerService.AssertNumberOfCalls(t, methodName, 1)
		})
	}
}
