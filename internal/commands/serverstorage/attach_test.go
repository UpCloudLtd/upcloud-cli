package serverstorage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/commands/server"
	"github.com/UpCloudLtd/cli/internal/commands/storage"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAttachStorageCommand(t *testing.T) {
	methodName := "AttachStorage"

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

	var serverDetails = upcloud.ServerDetails{
		Server: upcloud.Server{
			UUID:  UUID1,
			State: upcloud.ServerStateStarted,
		},
		VideoModel: "vga",
		Firewall:   "off",
	}

	for _, test := range []struct {
		name       string
		args       []string
		attacheReq request.AttachStorageRequest
		error      string
	}{
		{
			name:  "storage is missing",
			args:  []string{},
			error: "storage is required",
		},
		{
			name: "use default values",
			args: []string{"--storage", Storage1.Title},
			attacheReq: request.AttachStorageRequest{
				ServerUUID:  Server1.UUID,
				Type:        upcloud.StorageTypeDisk,
				Address:     "virtio",
				StorageUUID: Storage1.UUID,
				BootDisk:    0,
			},
		},
		{
			name: "flags mapped to the correct field",
			args: []string{
				"--storage", Storage1.Title,
				"--type", "cdrom",
				"--address", "ide",
				"--boot-disk", "1",
			},
			attacheReq: request.AttachStorageRequest{
				ServerUUID:  Server1.UUID,
				Type:        upcloud.StorageTypeCDROM,
				Address:     "ide",
				StorageUUID: Storage1.UUID,
				BootDisk:    1,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			server.CachedServers = nil
			storage.CachedStorages = nil

			mServerService := server.MockServerService{}
			mServerService.On("GetServers", mock.Anything).Return(servers, nil)

			mStorageService := MockStorageService{}
			mStorageService.On(methodName, &test.attacheReq).Return(&serverDetails, nil)
			mStorageService.On("GetStorages", mock.Anything).Return(storages, nil)

			cc := commands.BuildCommand(AttachCommand(&mServerService, &mStorageService), nil, config.New(viper.New()))
			err := cc.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = cc.MakeExecuteCommand()([]string{Server1.UUID})

			if test.error != "" {
				assert.Equal(t, test.error, err.Error())
			} else {
				mStorageService.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
