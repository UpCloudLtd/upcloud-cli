package storage

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteStorageCommand(t *testing.T) {
	targetMethod := "DeleteStorage"
	Storage2 := upcloud.Storage{
		UUID:   UUID2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	for _, test := range []struct {
		name        string
		args        []string
		methodCalls int
		expected    request.DeleteStorageRequest
	}{
		{
			name:        "delete storage",
			args:        []string{},
			methodCalls: 1,
			expected:    request.DeleteStorageRequest{UUID: Storage2.UUID},
		},
		{
			name:        "delete storage and keep only latest backup",
			args:        []string{"--backups", "keep_latest"},
			methodCalls: 1,
			expected:    request.DeleteStorageRequest{UUID: Storage2.UUID, Backups: request.DeleteStorageBackupsModeKeepLatest},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil
			conf := config.New()
			testCmd := DeleteCommand()
			mService := new(smock.Service)

			expected := test.expected
			mService.On(targetMethod, &expected).Return(nil, nil)
			mService.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage2}}, nil)

			c := commands.BuildCommand(testCmd, nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), Storage2.UUID)
			assert.Nil(t, err)

			mService.AssertNumberOfCalls(t, targetMethod, test.methodCalls)
		})
	}
}
