package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateBackupCommand(t *testing.T) {
	methodName := "CreateBackup"
	var Storage1 = upcloud.Storage{
		UUID:   Uuid1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	var Storage2 = upcloud.Storage{
		UUID:   Uuid2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	details := upcloud.StorageDetails{
		Storage: Storage1,
	}
	for _, test := range []struct {
		name     string
		args     []string
		expected request.CreateBackupRequest
		error    string
	}{
		{
			name:  "title is missing",
			args:  []string{},
			error: "title is required",
		},
		{
			name:     "title is provided",
			args:     []string{"--title", "test-title"},
			expected: request.CreateBackupRequest{Title: "test-title"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mss := mocks.MockStorageService{}
			mss.On("GetStorages", &request.GetStoragesRequest{}).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			mss.On(methodName, mock.Anything).Return(&details, nil)

			tc := commands.BuildCommand(CreateBackupCommand(&mss), nil, config.New(viper.New()))
			mocks.SetFlags(tc, test.args)

			_, err := tc.MakeExecuteCommand()([]string{Storage2.UUID})

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mss.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
