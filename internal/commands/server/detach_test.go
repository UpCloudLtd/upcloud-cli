package server

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

func TestDetachCommand(t *testing.T) {
	methodName := "DetachStorage"
	details := upcloud.ServerDetails{
		Server: Server1,
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
			serverSvc := MockServerService()
			storageSvc := MockStorageService()
			serverSvc.On("GetServers", mock.Anything).Return(servers, nil)
			storageSvc.On(methodName, mock.Anything).Return(&details, nil)

			tc := commands.BuildCommand(DetachCommand(serverSvc, storageSvc), nil, config.New(viper.New()))
			mocks.SetFlags(tc, test.args)

			results, err := tc.MakeExecuteCommand()([]string{Server1.UUID})
			for _, result := range results.([]interface{}) {
				assert.Equal(t, &details, result.(*upcloud.ServerDetails))
			}

			assert.Nil(t, err)

			storageSvc.AssertNumberOfCalls(t, methodName, test.methodCalls)
		})
	}
}
