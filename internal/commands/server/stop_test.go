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

func TestStopCommand(t *testing.T) {
	methodName := "StopServer"

	var Server1 = upcloud.Server{
		State: upcloud.ServerStateMaintenance,
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
			State: upcloud.ServerStateStopped,
			Title: "server-1-title",
			UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		},
	}

	dur120, _ := time.ParseDuration("120s")
	dur10, _ := time.ParseDuration("10s")

	for _, test := range []struct {
		name     string
		args     []string
		startReq request.StopServerRequest
	}{
		{
			name: "use default values",
			args: []string{},
			startReq: request.StopServerRequest{
				UUID:     Server1.UUID,
				Timeout:  dur120,
				StopType: upcloud.StopTypeSoft,
			},
		},
		{
			name: "flags mapped to the correct field",
			args: []string{
				"--timeout", "10",
				"--type", "hard",
			},
			startReq: request.StopServerRequest{
				UUID:     Server1.UUID,
				Timeout:  dur10,
				StopType: upcloud.StopTypeHard,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedServers = nil

			mServerService := MockServerService{}
			mServerService.On("GetServers", mock.Anything).Return(servers, nil)
			mServerService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: Server1.UUID}).Return(&details2, nil)
			mServerService.On(methodName, &test.startReq).Return(&details, nil)

			c := commands.BuildCommand(StopCommand(&mServerService), nil, config.New(viper.New()))
			c.SetFlags(test.args)

			c.MakeExecuteCommand()([]string{Server1.UUID})

			mServerService.AssertNumberOfCalls(t, methodName, 1)
		})
	}
}
