package storage

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"
	internal "github.com/UpCloudLtd/cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
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
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}
			mService.On(targetMethod, &test.expected).Return(&details, nil)

			c := commands.BuildCommand(testCmd, nil, config.New())
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.Command).Execute(commands.NewExecutor(conf, mService), "")
			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
