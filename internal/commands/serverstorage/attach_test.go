package serverstorage

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/mockexecute"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAttachStorageCommand(t *testing.T) {
	targetMethod := "AttachStorage"

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
			error: `required flag(s) "storage" not set`,
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
				"--boot-disk",
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
			conf := config.New()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}

			mService.On("GetServers", mock.Anything).Return(servers, nil)
			mService.On(targetMethod, &test.attacheReq).Return(&serverDetails, nil)
			mService.On("GetStorages", mock.Anything).Return(storages, nil)

			c := commands.BuildCommand(AttachCommand(), nil, conf)

			c.Cobra().SetArgs(append(test.args, Server1.UUID))
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" {
				assert.Error(t, err)
				assert.Equal(t, test.error, err.Error())
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
