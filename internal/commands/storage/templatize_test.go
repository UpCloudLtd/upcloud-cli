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

func TestTemplatizeCommand(t *testing.T) {
	methodName := "TemplatizeStorage"
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

			tc := commands.BuildCommand(TemplatizeCommand(mss), nil, config.New(viper.New()))
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
