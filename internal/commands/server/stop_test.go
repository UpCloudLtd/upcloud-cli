package server

import (
	"testing"
	"time"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStopCommand(t *testing.T) {
	targetMethod := "StopServer"

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
		name    string
		args    []string
		stopReq request.StopServerRequest
	}{
		{
			name: "use default values",
			args: []string{},
			stopReq: request.StopServerRequest{
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
			stopReq: request.StopServerRequest{
				UUID:     Server1.UUID,
				Timeout:  dur10,
				StopType: upcloud.StopTypeHard,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedServers = nil
			conf := config.New()
			testCmd := StopCommand()
			mService := new(smock.MockService)

			conf.Service = mService
			mService.On("GetServers", mock.Anything).Return(servers, nil)
			mService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: Server1.UUID}).Return(&details2, nil)
			mService.On(targetMethod, &test.stopReq).Return(&details, nil)

			c := commands.BuildCommand(testCmd, nil, conf)
			err := c.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = c.MakeExecuteCommand()([]string{Server1.UUID})
			assert.NoError(t, err)

			mService.AssertNumberOfCalls(t, targetMethod, 1)
		})
	}
}
