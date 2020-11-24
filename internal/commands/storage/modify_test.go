package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestModifyCommand(t *testing.T) {
	methodName := "ModifyStorage"
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
		name        string
		args        []string
		methodCalls int
	}{
		{
			name:        "Backend called, details returned",
			args:        []string{},
			methodCalls: 1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mss := MockStorageService()
			mss.On(methodName, mock.Anything).Return(&details, nil)
			mss.On("GetStorageDetails", mock.Anything).Return(&details, nil)

			tc := commands.BuildCommand(ModifyCommand(mss), nil, config.New(viper.New()))
			mocks.SetFlags(tc, test.args)

			results, err := tc.MakeExecuteCommand()([]string{Storage2.UUID})
			for _, result := range results.([]interface{}) {
				assert.Equal(t, &details, result.(*upcloud.StorageDetails))
			}

			assert.Nil(t, err)

			mss.AssertNumberOfCalls(t, methodName, test.methodCalls)
		})
	}
}
