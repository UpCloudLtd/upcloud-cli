package storage

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/mockexecute"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCloneCommand(t *testing.T) {
	targetMethod := "CloneStorage"

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
		expected request.CloneStorageRequest
	}{
		{
			name: "using default tier",
			args: []string{"--zone", "test-zone", "--title", "test-title"},
			expected: request.CloneStorageRequest{
				UUID:  Storage2.UUID,
				Zone:  "test-zone",
				Tier:  "hdd",
				Title: "test-title",
			},
		},
		{
			name: "tier from args",
			args: []string{"--zone", "test-zone", "--title", "test-title", "--tier", "abc"},
			expected: request.CloneStorageRequest{
				UUID:  Storage2.UUID,
				Zone:  "test-zone",
				Tier:  "abc",
				Title: "test-title",
			},
		},
		{
			name: "title is missing",
			args: []string{
				"--zone", "zone",
			},
			error: `required flag(s) "title" not set`,
		},
		{
			name: "zone is missing",
			args: []string{
				"--title", "title",
			},
			error: `required flag(s) "zone" not set`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil
			conf := config.New()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}

			mService.On(targetMethod, &test.expected).Return(&details, nil)
			mService.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)

			c := commands.BuildCommand(CloneCommand(), nil, conf)

			c.Cobra().SetArgs(append(test.args, Storage2.UUID))
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" {
				assert.Error(t, err)
				assert.Equal(t, test.error, err.Error())
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
