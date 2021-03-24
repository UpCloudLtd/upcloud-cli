package storage

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTemplatizeCommand(t *testing.T) {
	targetMethod := "TemplatizeStorage"
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
		expected request.TemplatizeStorageRequest
	}{
		{
			name:  "Backend called with no args",
			args:  []string{},
			error: "title is required",
		},
		{
			name: "Backend called with title",
			args: []string{"--title", "test-title"},
			expected: request.TemplatizeStorageRequest{
				UUID:  Storage2.UUID,
				Title: "test-title",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.MockService{}
			mService.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			mService.On(targetMethod, &test.expected).Return(&details, nil)

			tc := commands.BuildCommand(TemplatizeCommand(&mService), nil, config.New(viper.New()))
			err := tc.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = tc.MakeExecuteCommand()([]string{Storage2.UUID})
			if test.error != "" {
				assert.Errorf(t, err, "title is required")
			} else {
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
