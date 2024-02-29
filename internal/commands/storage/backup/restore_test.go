package storagebackup

import (
	"testing"

	"github.com/gemalto/flume"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRestoreBackupCommand(t *testing.T) {
	targetMethod := "RestoreBackup"

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
	for _, test := range []struct {
		name        string
		args        []string
		methodCalls int
		expected    request.RestoreBackupRequest
	}{
		{
			name:        "Backend called",
			args:        []string{},
			methodCalls: 1,
			expected:    request.RestoreBackupRequest{UUID: Storage2.UUID},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			mService := new(smock.Service)

			mService.On(targetMethod, mock.Anything).Return(nil, nil)

			c := commands.BuildCommand(RestoreBackupCommand(), nil, conf)
			err := c.Cobra().Flags().Parse(test.args)
			assert.NoError(t, err)

			_, err = c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), Storage2.UUID)
			assert.Nil(t, err)
			mService.AssertNumberOfCalls(t, targetMethod, test.methodCalls)
		})
	}
}
