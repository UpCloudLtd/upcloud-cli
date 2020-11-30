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

func TestEjectCDROMCommand(t *testing.T) {
	methodName := "EjectCDROM"

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
		t.Run(test.name, func(t *testing.T) {
			CachedServers = nil
			storage.CachedStorages = nil

			mServerService := MockServerService{}
			mServerService.On("GetServers", mock.Anything).Return(servers, nil)

			mStorageService := MockStorageService{}
			mStorageService.On(methodName, &test.ejectReq).Return(&details, nil)

			tc := commands.BuildCommand(EjectCommand(&mServerService, &mStorageService), nil, config.New(viper.New()))
			tc.SetFlags(test.args)

			_, err := tc.MakeExecuteCommand()([]string{Server1.UUID})

			assert.Nil(t, err)
			mStorageService.AssertNumberOfCalls(t, methodName, 1)
		})
	}
}
