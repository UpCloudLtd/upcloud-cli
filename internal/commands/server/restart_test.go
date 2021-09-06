package server

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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
			testSimpleServerCommand(t, RestartCommand(), servers, Server1, details, methodName, &test.restartReq, &details2, test.args)
		})
	}
}
