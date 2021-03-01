package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestStartCommand(t *testing.T) {
	methodName := "StartServer"

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
			State: upcloud.ServerStateStarted,
			Title: "server-1-title",
			UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		},
	}

	dur120, _ := time.ParseDuration("120s")
	dur10, _ := time.ParseDuration("10s")

	for _, test := range []struct {
		name     string
		args     []string
		startReq request.StartServerRequest
	}{
		{
			name: "use default values",
			args: []string{},
			startReq: request.StartServerRequest{
				UUID:      Server1.UUID,
				Timeout:   dur120,
				AvoidHost: 0,
				Host:      0,
			},
		},
		{
			name: "flags mapped to the correct field",
			args: []string{
				"--avoid-host", "5678",
				"--timeout", "10",
				"--host", "1234",
			},
			startReq: request.StartServerRequest{
				UUID:      Server1.UUID,
				Timeout:   dur10,
				AvoidHost: 5678,
				Host:      1234,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedServers = nil

			mServerService := MockServerService{}
			mServerService.On("GetServers", mock.Anything).Return(servers, nil)
			mServerService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: Server1.UUID}).Return(&details2, nil)
			mServerService.On(methodName, &test.startReq).Return(&details, nil)

			c := commands.BuildCommand(StartCommand(&mServerService), nil, config.New(viper.New()))
			err := c.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = c.MakeExecuteCommand()([]string{Server1.UUID})
			assert.NoError(t, err)

			mServerService.AssertNumberOfCalls(t, methodName, 1)
		})
	}
}
