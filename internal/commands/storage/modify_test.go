package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModifyCommandExistingBackupRule(t *testing.T) {
	methodName := "ModifyStorage"
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
	var StorageDetails1 = upcloud.StorageDetails{
		Storage: Storage1,
		BackupRule: &upcloud.BackupRule{
			Interval:  "sun",
			Time:      "0800",
			Retention: 5,
		},
	}
	var Storage2 = upcloud.Storage{
		UUID:   UUID2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}
	var StorageDetails2 = upcloud.StorageDetails{
		Storage:    Storage2,
		BackupRule: &upcloud.BackupRule{Time: "", Interval: "", Retention: 0},
	}

	for _, test1 := range []struct {
		name        string
		args        []string
		storage     upcloud.Storage
		methodCalls int
		expected    request.ModifyStorageRequest
		error       string
	}{
		{
			name:        "without backup rule update of existing backup rule",
			args:        []string{"--size", "50"},
			storage:     Storage1,
			methodCalls: 1,
			expected: request.ModifyStorageRequest{
				UUID:       Storage1.UUID,
				Size:       50,
				BackupRule: StorageDetails1.BackupRule,
			},
		},
		{
			name:        "modifying existing backup rule without time",
			args:        []string{"--size", "50", "--backup-interval", "mon"},
			storage:     Storage1,
			methodCalls: 1,
			expected: request.ModifyStorageRequest{
				UUID: Storage1.UUID,
				Size: 50,
				BackupRule: &upcloud.BackupRule{
					Interval:  "mon",
					Time:      StorageDetails1.BackupRule.Time,
					Retention: StorageDetails1.BackupRule.Retention,
				},
			},
		},
	} {
		t.Run(test1.name, func(t *testing.T) {
			CachedStorages = nil
			mss := MockStorageService{}
			mss.On("GetStorages").Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1}}, nil)
			mss.On(methodName, &test1.expected).Return(&StorageDetails1, nil)
			mss.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: Storage1.UUID}).Return(&StorageDetails1, nil)

			tc := commands.BuildCommand(ModifyCommand(&mss), nil, config.New(viper.New()))
			err := tc.SetFlags(test1.args)
			assert.NoError(t, err)

			_, err = tc.MakeExecuteCommand()([]string{test1.storage.UUID})

			if test1.error != "" {
				assert.Equal(t, test1.error, err.Error())
			} else {
				assert.Nil(t, err)
				mss.AssertNumberOfCalls(t, methodName, test1.methodCalls)
			}
		})
	}

	for _, test2 := range []struct {
		name        string
		args        []string
		storage     upcloud.Storage
		methodCalls int
		expected    request.ModifyStorageRequest
		error       string
	}{
		{
			name:        "modifying existing backup rule without time",
			args:        []string{"--size", "50", "--backup-interval", "mon"},
			storage:     Storage1,
			methodCalls: 1,
			expected: request.ModifyStorageRequest{
				UUID: Storage1.UUID,
				Size: 50,
				BackupRule: &upcloud.BackupRule{
					Interval:  "mon",
					Time:      StorageDetails1.BackupRule.Time,
					Retention: StorageDetails1.BackupRule.Retention,
				},
			},
		},
	} {
		t.Run(test2.name, func(t *testing.T) {
			CachedStorages = nil
			mss := MockStorageService{}
			mss.On("GetStorages").Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1}}, nil)
			mss.On(methodName, &test2.expected).Return(&StorageDetails1, nil)
			mss.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: Storage1.UUID}).Return(&StorageDetails1, nil)

			tc := commands.BuildCommand(ModifyCommand(&mss), nil, config.New(viper.New()))
			err := tc.SetFlags(test2.args)
			assert.NoError(t, err)

			_, err = tc.MakeExecuteCommand()([]string{test2.storage.UUID})

			if test2.error != "" {
				assert.Equal(t, test2.error, err.Error())
			} else {
				assert.Nil(t, err)
				mss.AssertNumberOfCalls(t, methodName, test2.methodCalls)
			}
		})
	}

	for _, test3 := range []struct {
		name        string
		args        []string
		storage     upcloud.Storage
		methodCalls int
		expected    request.ModifyStorageRequest
		error       string
	}{
		{
			name:        "without backup rule update of non-existing backup rule",
			args:        []string{"--size", "50"},
			storage:     Storage2,
			methodCalls: 1,
			expected: request.ModifyStorageRequest{
				UUID: Storage2.UUID,
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
					Time:      "1200",
					Retention: defaultBackupRuleParams.Retention,
					Interval:  defaultBackupRuleParams.Interval,
				},
			},
		},
		{
			name:        "adding backup rule without backup time",
			args:        []string{"--size", "50", "--backup-retention", "10"},
			storage:     Storage2,
			methodCalls: 1,
			error:       "backup-time is required",
		},
	} {
		t.Run(test3.name, func(t *testing.T) {
			CachedStorages = nil
			mss := MockStorageService{}
			mss.On("GetStorages").Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage2}}, nil)
			mss.On(methodName, &test3.expected).Return(&StorageDetails2, nil)
			mss.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: Storage2.UUID}).Return(&StorageDetails2, nil)

			tc := commands.BuildCommand(ModifyCommand(&mss), nil, config.New(viper.New()))
			err := tc.SetFlags(test3.args)
			assert.NoError(t, err)

			_, err = tc.MakeExecuteCommand()([]string{test3.storage.UUID})

			if test3.error != "" {
				assert.Equal(t, test3.error, err.Error())
			} else {
				assert.Nil(t, err)
				mss.AssertNumberOfCalls(t, methodName, test3.methodCalls)
			}
		})
	}
}
