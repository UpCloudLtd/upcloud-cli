package server_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands/server"
)

func TestStartCommand(t *testing.T) {
	t.Parallel()
	targetMethod := "StartServer"

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
			State: upcloud.ServerStateStarted,
			Title: "server-1-title",
			UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		},
	}

	for _, test := range []struct {
		name     string
		args     []string
		startReq request.StartServerRequest
	}{
		{
			name: "use default values",
			args: []string{},
			startReq: request.StartServerRequest{
				UUID: Server1.UUID,
			},
		},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server.SimpleServerCommandTests(t, server.StartCommand(), servers, Server1, details2, targetMethod, &test.startReq, &details, test.args)
		})
	}
}
