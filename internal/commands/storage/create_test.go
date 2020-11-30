package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCommand(t *testing.T) {
	methodName := "CreateStorage"
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
		error    string
		expected request.CreateStorageRequest
	}{
		{
			name: "create with default values, no backup rule",
			args: []string{
				"--title", "create-storage-test", "" +
					"--zone", "abc"},
			expected: request.CreateStorageRequest{
				Size:       DefaultCreateParams.Size,
				Tier:       DefaultCreateParams.Tier,
				Title:      "create-storage-test",
				Zone:       "abc",
				BackupRule: nil,
			},
		},
		{
			name: "create with default values, with backup rule",
			args: []string{"--title", "create-storage-test", "--zone", "abc", "--backup-time", "09:00"},
			expected: request.CreateStorageRequest{
				Size:  DefaultCreateParams.Size,
				Tier:  DefaultCreateParams.Tier,
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
			mss := mocks.MockStorageService{}
			mss.On(methodName, &test.expected).Return(&details, nil)

			tc := commands.BuildCommand(CreateCommand(&mss), nil, config.New(viper.New()))
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
