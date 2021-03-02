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

func TestRestoreBackupCommand(t *testing.T) {
	methodName := "RestoreBackup"

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
			mss := MockStorageService{}
			mss.On(methodName, mock.Anything).Return(nil, nil)

			tc := commands.BuildCommand(RestoreBackupCommand(&mss), nil, config.New(viper.New()))
			err := tc.SetFlags(test.args)
			assert.NoError(t, err)

			_, err = tc.MakeExecuteCommand()([]string{Storage2.UUID})
			assert.Nil(t, err)
			mss.AssertNumberOfCalls(t, methodName, test.methodCalls)
		})
	}
}
