package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeleteStorageCommand(t *testing.T) {
	methodName := "DeleteStorage"
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
	for _, test := range []struct {
		name        string
		args        []string
		methodCalls int
		expected    request.DeleteStorageRequest
	}{
		{
			name:        "Backend called",
			args:        []string{},
			methodCalls: 1,
			expected:    request.DeleteStorageRequest{UUID: Storage2.UUID},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil
			mss := mocks.MockStorageService{}
			mss.On(methodName, &test.expected).Return(nil, nil)
			mss.On("GetStorages", mock.Anything).Return(&upcloud.Storages{Storages: []upcloud.Storage{Storage2}}, nil)

			tc := commands.BuildCommand(DeleteCommand(&mss), nil, config.New(viper.New()))
			mocks.SetFlags(tc, test.args)

			_, err := tc.MakeExecuteCommand()([]string{Storage2.UUID})
			assert.Nil(t, err)

			mss.AssertNumberOfCalls(t, methodName, test.methodCalls)
		})
	}
}
