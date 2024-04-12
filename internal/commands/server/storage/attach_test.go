package serverstorage

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAttachStorageCommand(t *testing.T) {
	targetMethod := "AttachStorage"

	Storage1 := upcloud.Storage{
		UUID:   UUID1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	storages := &upcloud.Storages{
		Storages: []upcloud.Storage{
			Storage1,
		},
	}

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

	serverDetails := upcloud.ServerDetails{
		Server: upcloud.Server{
			UUID:  UUID1,
			State: upcloud.ServerStateStarted,
		},
		VideoModel: "vga",
		Firewall:   "off",
	}

	for _, test := range []struct {
		name      string
		args      []string
		attachReq request.AttachStorageRequest
		error     string
	}{
		{
			name:  "storage is missing",
			args:  []string{},
			error: `required flag(s) "storage" not set`,
		},
		{
			name: "use default values",
			args: []string{"--storage", Storage1.Title},
			attachReq: request.AttachStorageRequest{
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
			attachReq: request.AttachStorageRequest{
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

			mService.On("GetServers", mock.Anything).Return(servers, nil)
			attachReq := test.attachReq
			mService.On(targetMethod, &attachReq).Return(&serverDetails, nil)
			mService.On("GetStorages", mock.Anything).Return(storages, nil)

			c := commands.BuildCommand(AttachCommand(), nil, conf)

			c.Cobra().SetArgs(append(test.args, Server1.UUID))
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
