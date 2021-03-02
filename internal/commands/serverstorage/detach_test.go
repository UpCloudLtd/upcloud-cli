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

func TestDetachCommand(t *testing.T) {
	methodName := "DetachStorage"

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

	for _, test := range []struct {
		name      string
		args      []string
		detachReq request.DetachStorageRequest
		error     string
	}{
		{
			name:  "Address missing",
			args:  []string{},
			error: "address is required",
		},
		{
			name: "Address provided",
			args: []string{"--address", "ide"},
			detachReq: request.DetachStorageRequest{
				Address:    "ide",
				ServerUUID: Server1.UUID,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			server.CachedServers = nil
			storage.CachedStorages = nil

			mServerService := server.MockServerService{}
			mServerService.On("GetServers", mock.Anything).Return(servers, nil)

			mStorageService := MockStorageService{}
			mStorageService.On(methodName, &test.detachReq).Return(&details, nil)

			tc := commands.BuildCommand(DetachCommand(&mServerService, &mStorageService), nil, config.New(viper.New()))
			err := tc.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = tc.MakeExecuteCommand()([]string{Server1.UUID})

			if test.error != "" {
				assert.Equal(t, test.error, err.Error())
			} else {
				mStorageService.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
