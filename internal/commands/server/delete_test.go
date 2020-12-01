package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeleteServerCommand(t *testing.T) {
	deleteServer := "DeleteServer"
	deleteServerAndStorages := "DeleteServerAndStorages"

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

	for _, test := range []struct {
		name                   string
		args                   []string
		deleteServCalls        int
		deleteServStorageCalls int
	}{
		{
			name:                   "Delete-storages true",
			args:                   []string{"--delete-storages"},
			deleteServCalls:        0,
			deleteServStorageCalls: 1,
		},
		{
			name:                   "Delete-storages false",
			args:                   []string{},
			deleteServCalls:        1,
			deleteServStorageCalls: 0,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mss := MockServerService{}
			mss.On(deleteServer, mock.Anything).Return(nil, nil)
			mss.On(deleteServerAndStorages, mock.Anything).Return(nil, nil)
			mss.On("GetServers", mock.Anything).Return(servers, nil)

			tc := commands.BuildCommand(DeleteCommand(&mss), nil, config.New(viper.New()))
			tc.SetFlags(test.args)

			results, err := tc.MakeExecuteCommand()([]string{Server1.UUID})
			for _, result := range results.([]interface{}) {
				assert.Nil(t, result)
			}

			assert.Nil(t, err)

			mss.AssertNumberOfCalls(t, deleteServer, test.deleteServCalls)
			mss.AssertNumberOfCalls(t, deleteServerAndStorages, test.deleteServStorageCalls)
		})
	}
}
