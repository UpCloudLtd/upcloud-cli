package storagebackup

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBackupCommand(t *testing.T) {
	targetMethod := "CreateBackup"
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
	details := upcloud.StorageDetails{
		Storage: Storage1,
	}
	for _, test := range []struct {
		name     string
		args     []string
		expected request.CreateBackupRequest
		error    string
	}{
		{
			name:  "title is missing",
			args:  []string{},
			error: `required flag(s) "title" not set`,
		},
		{
			name:     "title is provided",
			args:     []string{"--title", "test-title"},
			expected: request.CreateBackupRequest{Title: "test-title"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			mService := new(smock.Service)

			mService.On("GetStorages", &request.GetStoragesRequest{}).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage1, Storage2}}, nil)
			mService.On(targetMethod, mock.Anything).Return(&details, nil)

			c := commands.BuildCommand(CreateBackupCommand(), nil, conf)

			c.Cobra().SetArgs(append(test.args, Storage2.UUID))
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
