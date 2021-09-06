package server_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands/server"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEjectCDROMCommand(t *testing.T) {
	t.Parallel()
	targetMethod := "EjectCDROM"

	Server1 := upcloud.Server{
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
	servers := &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
		},
	}

	details := upcloud.ServerDetails{
		Server: Server1,
	}

	for _, test := range []struct {
		name     string
		args     []string
		ejectReq request.EjectCDROMRequest
	}{
		{
			name: "Backend called, details returned",
			args: []string{},
			ejectReq: request.EjectCDROMRequest{
				ServerUUID: Server1.UUID,
			},
		},
	} {
		// grab a local reference for parallel tests
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			conf := config.New()
			testCmd := server.EjectCommand()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}
			mService.On("GetServers", mock.Anything).Return(servers, nil)
			mService.On(targetMethod, &test.ejectReq).Return(&details, nil)

			c := commands.BuildCommand(testCmd, nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(
				commands.NewExecutor(conf, mService, flume.New("test")),
				Server1.UUID,
			)

			assert.Nil(t, err)
			mService.AssertNumberOfCalls(t, targetMethod, 1)
		})
	}
}
