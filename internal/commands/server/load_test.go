package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/storage"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestLoadCDROMCommand(t *testing.T) {
	methodName := "LoadCDROM"

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
			CachedServers = nil
			storage.CachedStorages = nil

			mServerService := MockServerService{}
			mServerService.On("GetServers", mock.Anything).Return(servers, nil)

			mStorageService := MockStorageService{}
			mStorageService.On(methodName, &test.loadReq).Return(&details, nil)
			mStorageService.On("GetStorages", mock.Anything).Return(storages, nil)

			cc := commands.BuildCommand(LoadCommand(&mServerService, &mStorageService), nil, config.New(viper.New()))
			cc.SetFlags(test.args)

			_, err := cc.MakeExecuteCommand()([]string{Server1.UUID})

			if test.error != "" {
				assert.Equal(t, test.error, err.Error())
			} else {
				mStorageService.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
