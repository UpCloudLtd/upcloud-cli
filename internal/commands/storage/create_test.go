package storage

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	targetMethod := "CreateStorage"
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
	details := upcloud.StorageDetails{
		Storage: Storage1,
	}
	for _, test := range []struct {
		name     string
		args     []string
		error    string
		expected request.CreateStorageRequest
	}{
		{
			name: "create with default values, no backup rule",
			args: []string{
				"--title", "create-storage-test", "" +
					"--zone", "abc"},
			expected: request.CreateStorageRequest{
				Size:       defaultCreateParams.Size,
				Tier:       defaultCreateParams.Tier,
				Title:      "create-storage-test",
				Zone:       "abc",
				BackupRule: nil,
			},
		},
		{
			name: "create with default values, with backup rule",
			args: []string{"--title", "create-storage-test", "--zone", "abc", "--backup-time", "09:00"},
			expected: request.CreateStorageRequest{
				Size:  defaultCreateParams.Size,
				Tier:  defaultCreateParams.Tier,
				Title: "create-storage-test",
				Zone:  "abc",
				BackupRule: &upcloud.BackupRule{
					Interval:  "daily",
					Time:      "0900",
					Retention: 7,
				},
			},
		},
		{
			name: "create with non default values",
			args: []string{
				"--title", "create-storage-test",
				"--zone", "abc",
				"--size", "30",
				"--tier", "xyz",
				"--backup-time", "09:00",
				"--backup-retention", "10",
				"--backup-interval", "mon",
			},
			expected: request.CreateStorageRequest{
				Size:  30,
				Tier:  "xyz",
				Title: "create-storage-test",
				Zone:  "abc",
				BackupRule: &upcloud.BackupRule{
					Interval:  "mon",
					Time:      "0900",
					Retention: 10,
				},
			},
		},
		{
			name: "title is missing",
			args: []string{
				"--size", "10",
				"--zone", "zone",
			},
			error: "size, title and zone are required",
		},
		{
			name: "zone is missing",
			args: []string{
				"--title", "title",
				"--size", "10",
			},
			error: "size, title and zone are required",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := mock.MockService{}
			mService.On(targetMethod, &test.expected).Return(&details, nil)

			tc := commands.BuildCommand(CreateCommand(&mService), nil, config.New(viper.New()))
			err := tc.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = tc.MakeExecuteCommand()([]string{Storage2.UUID})
			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
