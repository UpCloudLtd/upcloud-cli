package server

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestRelocateCommand(t *testing.T) {
	targetMethod := "RelocateServer"

	server := upcloud.Server{
		State: upcloud.ServerStateStopped,
		Title: "server-1-title",
		UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		Zone:  "fi-priv-example",
	}
	serverDetails := upcloud.ServerDetails{
		Server: server,
	}

	for _, test := range []struct {
		name    string
		args    []string
		request request.RelocateServerRequest
		error   string
	}{
		{
			name:  "no args",
			args:  []string{},
			error: `required flag(s) "zone" not set`,
		},
		{
			name: "success",
			args: []string{"--zone", server.Zone, server.UUID},
			request: request.RelocateServerRequest{
				UUID: server.UUID,
				Zone: server.Zone,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			req := test.request
			mService.On(targetMethod, &req).Return(&serverDetails, nil)

			conf := config.New()
			command := commands.BuildCommand(RelocateCommand(), nil, conf)
			command.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(command, &mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
