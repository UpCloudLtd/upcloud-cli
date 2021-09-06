package storage

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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

	for _, test := range []struct {
		name        string
		args        []string
		storage     upcloud.Storage
		details     upcloud.StorageDetails
		methodCalls int
		expected    request.ModifyStorageRequest
		error       string
	}{
		{
			name:        "without backup rule update of existing backup rule",
			args:        []string{"--size", "50"},
			storage:     Storage1,
			details:     StorageDetails1,
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
			details:     StorageDetails1,
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
		{
			name:        "modifying existing backup rule without time",
			args:        []string{"--size", "50", "--backup-interval", "mon"},
			storage:     Storage1,
			details:     StorageDetails1,
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
		{
			name:        "without backup rule update of non-existing backup rule",
			args:        []string{"--size", "50"},
			storage:     Storage2,
			details:     StorageDetails2,
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
			details:     StorageDetails2,
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
			details:     StorageDetails2,
			methodCalls: 1,
			error:       "backup-time is required",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil
			conf := config.New()
			testCmd := ModifyCommand()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}
			mService.On("GetStorages").Return(&upcloud.Storages{Storages: []upcloud.Storage{test.storage}}, nil)
			mService.On(targetMethod, &test.expected).Return(&test.details, nil)
			mService.On("GetStorageDetails", &request.GetStorageDetailsRequest{UUID: test.storage.UUID}).Return(&test.details, nil)

			c := commands.BuildCommand(testCmd, nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), test.storage.UUID)

			if test.error != "" {
				assert.Error(t, err)
				assert.Equal(t, test.error, err.Error())
			} else {
				assert.Nil(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, test.methodCalls)
			}
		})
	}
}
