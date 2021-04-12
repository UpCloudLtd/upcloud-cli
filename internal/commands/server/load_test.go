package server

import (
	internal "github.com/UpCloudLtd/cli/internal/service"
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoadCDROMCommand(t *testing.T) {
	targetMethod := "LoadCDROM"

	var Server1 = upcloud.Server{
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
	var servers = &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
		},
	}

	details := upcloud.ServerDetails{
		Server: Server1,
	}

	var Storage1 = upcloud.Storage{
		UUID:   UUID1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	var storages = &upcloud.Storages{
		Storages: []upcloud.Storage{
			Storage1,
		},
	}

	for _, test := range []struct {
		name    string
		args    []string
		loadReq request.LoadCDROMRequest
		error   string
	}{
		{
			name:  "storage is missing",
			args:  []string{},
			error: "storage is required",
		},
		{
			name: "storage is provided",
			args: []string{"--storage", Storage1.UUID},
			loadReq: request.LoadCDROMRequest{
				ServerUUID:  Server1.UUID,
				StorageUUID: Storage1.UUID,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := LoadCommand()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}

			mService.On("GetServers", mock.Anything).Return(servers, nil)
			mService.On("GetStorages", mock.Anything).Return(storages, nil)
			mService.On(targetMethod, &test.loadReq).Return(&details, nil)

			c := commands.BuildCommand(testCmd, nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.Command).Execute(
				commands.NewExecutor(conf, mService),
				Server1.UUID,
			)

			if test.error != "" {
				if err == nil {
					t.Errorf("expected error '%v', got nil", test.error)
				} else {
					assert.Equal(t, test.error, err.Error())
				}
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
