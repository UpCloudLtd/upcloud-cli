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
			error: `required flag(s) "title" not set`,
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
			CachedStorages = nil
			conf := config.New()
			mService := new(smock.Service)

			conf.Service = internal.Wrapper{Service: mService}

			mService.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			mService.On(targetMethod, &test.expected).Return(&details, nil)

			c := commands.BuildCommand(TemplatizeCommand(), nil, conf)

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
