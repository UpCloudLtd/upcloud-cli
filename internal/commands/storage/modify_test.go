package storage

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
)

func TestModifyCommandExistingBackupRule(t *testing.T) {
	targetMethod := "ModifyStorage"
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
	StorageDetails1 := upcloud.StorageDetails{
		Storage: Storage1,
		BackupRule: &upcloud.BackupRule{
			Interval:  "sun",
			Time:      "0800",
			Retention: 5,
		},
	}
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
	StorageDetails2 := upcloud.StorageDetails{
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
			conf := config.New()
			testCmd := ModifyCommand()
			mService := new(smock.Service)

			mService.On("GetStorages").Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1}}, nil)
			mService.On(targetMethod, &test1.expected).Return(&StorageDetails1, nil)
			mService.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: Storage1.UUID}).Return(&StorageDetails1, nil)

			c := commands.BuildCommand(testCmd, nil, conf)
			err := c.Cobra().Flags().Parse(test1.args)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), test1.storage.UUID)

			if test1.error != "" {
				assert.EqualError(t, err, test1.error)
			} else {
				assert.Nil(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, test1.methodCalls)
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
			conf := config.New()
			testCmd := ModifyCommand()
			mService := new(smock.Service)

			mService.On("GetStorages").Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1}}, nil)
			mService.On(targetMethod, &test2.expected).Return(&StorageDetails1, nil)
			mService.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: Storage1.UUID}).Return(&StorageDetails1, nil)

			c := commands.BuildCommand(testCmd, nil, config.New())
			err := c.Cobra().Flags().Parse(test2.args)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), test2.storage.UUID)

			if test2.error != "" {
				assert.EqualError(t, err, test2.error)
			} else {
				assert.Nil(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, test2.methodCalls)
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
			conf := config.New()
			testCmd := ModifyCommand()
			mService := new(smock.Service)

			mService.On("GetStorages").Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage2}}, nil)
			mService.On(targetMethod, &test3.expected).Return(&StorageDetails2, nil)
			mService.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: Storage2.UUID}).Return(&StorageDetails2, nil)

			c := commands.BuildCommand(testCmd, nil, config.New())
			err := c.Cobra().Flags().Parse(test3.args)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), test3.storage.UUID)

			if test3.error != "" {
				assert.EqualError(t, err, test3.error)
			} else {
				assert.Nil(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, test3.methodCalls)
			}
		})
	}
}

func TestModifyCommandAutoresize(t *testing.T) {
	t.Run("modifying storage size with filesystem autoresize enabled", func(t *testing.T) {
		conf := config.New()
		testCmd := ModifyCommand()
		mService := new(smock.Service)
		UUID := "some_storage_id"

		mGetDetailsResponse := upcloud.StorageDetails{
			Storage:    upcloud.Storage{Size: 45},
			BackupRule: &upcloud.BackupRule{},
		}

		mModifyResponse := upcloud.StorageDetails{
			Storage: upcloud.Storage{
				Size: 50,
			},
		}

		mResizeResponse := upcloud.ResizeStorageFilesystemBackup{
			UUID: "resize_backup",
		}

		mService.On("ModifyStorage", &request.ModifyStorageRequest{UUID: UUID, Size: 50}).Return(&mModifyResponse, nil)
		mService.On("ResizeStorageFilesystem", &request.ResizeStorageFilesystemRequest{UUID: UUID}).Return(&mResizeResponse, nil)
		mService.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: UUID}).Return(&mGetDetailsResponse, nil)

		c := commands.BuildCommand(testCmd, nil, conf)
		err := c.Cobra().Flags().Parse([]string{"--size", "50", "--enable-filesystem-autoresize"})
		assert.NoError(t, err)

		output, err := c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), UUID)
		assert.NoError(t, err)
		mService.AssertNumberOfCalls(t, "ModifyStorage", 1)
		mService.AssertNumberOfCalls(t, "ResizeStorageFilesystem", 1)

		json, err := output.MarshalJSON()
		assert.NoError(t, err)
		assert.Contains(t, string(json), "latest_resize_backup")
		assert.Contains(t, string(json), "resize_backup")
	})
}
