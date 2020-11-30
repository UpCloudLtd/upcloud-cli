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

func TestModifyCommand(t *testing.T) {
	methodName := "ModifyStorage"
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
	var StorageDetails1 = upcloud.StorageDetails{
		Storage: Storage1,
		BackupRule: &upcloud.BackupRule{
			Interval:  "sun",
			Time:      "0800",
			Retention: 5,
		},
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
	var StorageDetails2 = upcloud.StorageDetails{
		Storage: Storage2,
	}

	for _, test := range []struct {
		name        string
		args        []string
		storage     upcloud.Storage
		methodCalls int
		expected    request.ModifyStorageRequest
	}{
		{
			name:        "without backup rule",
			args:        []string{"--size", "50"},
			storage:     Storage1,
			methodCalls: 1,
			expected: request.ModifyStorageRequest{
				UUID: Storage1.UUID,
				Size: 50,
			},
		},
		{
			name:        "adding backup rule",
			args:        []string{"--size", "50", "--backup-time", "12:00"},
			storage:     Storage2,
			methodCalls: 1,
			expected: request.ModifyStorageRequest{
				UUID: Storage2.UUID,
				Size: 50,
				BackupRule: &upcloud.BackupRule{
					Time: "1200",
				},
			},
		},
		{
			name:        "modifying existing backup rule",
			args:        []string{"--size", "50", "--backup-time", "12:00", "--backup-interval", "mon"},
			storage:     Storage1,
			methodCalls: 1,
			expected: request.ModifyStorageRequest{
				UUID: Storage1.UUID,
				Size: 50,
				BackupRule: &upcloud.BackupRule{
					Interval: "mon",
					Time:     "1200",
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil
			mss := mocks.MockStorageService{}
			mss.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			mss.On(methodName, mock.Anything).Return(&StorageDetails1, nil)
			mss.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: Storage1.UUID}).Return(&StorageDetails1, nil)
			mss.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: Storage2.UUID}).Return(&StorageDetails2, nil)

			tc := commands.BuildCommand(ModifyCommand(&mss), nil, config.New(viper.New()))
			mocks.SetFlags(tc, test.args)

			_, err := tc.MakeExecuteCommand()([]string{test.storage.UUID})
			assert.Nil(t, err)

			mss.AssertNumberOfCalls(t, methodName, test.methodCalls)
		})
	}
}
