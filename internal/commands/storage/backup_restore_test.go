package storage

import (
	"testing"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRestoreBackupCommand(t *testing.T) {
	targetMethod := "RestoreBackup"

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
			mService := smock.MockService{}
			mService.On(targetMethod, mock.Anything).Return(nil, nil)

			tc := commands.BuildCommand(RestoreBackupCommand(&mService), nil, config.New())
			err := tc.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = tc.MakeExecuteCommand()([]string{Storage2.UUID})
			assert.Nil(t, err)
			mService.AssertNumberOfCalls(t, targetMethod, test.methodCalls)
		})
	}
}
