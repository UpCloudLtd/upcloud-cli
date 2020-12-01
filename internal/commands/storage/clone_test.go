package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCloneCommand(t *testing.T) {
	methodName := "CloneStorage"

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
			error: "title and zone are required",
		},
		{
			name: "zone is missing",
			args: []string{
				"--title", "title",
			},
			error: "title and zone are required",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil
			mss := MockStorageService{}
			mss.On(methodName, &test.expected).Return(&details, nil)
			mss.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)

			tc := commands.BuildCommand(CloneCommand(&mss), nil, config.New(viper.New()))
			tc.SetFlags(test.args)

			_, err := tc.MakeExecuteCommand()([]string{Storage2.UUID})

			if test.error != "" {
				assert.Errorf(t, err, test.error)
			} else {
				mss.AssertNumberOfCalls(t, methodName, 1)
			}
		})
	}
}
