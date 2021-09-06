package server

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	internal "github.com/UpCloudLtd/upcloud-cli/internal/service"
)

func testSimpleServerCommand(t *testing.T, testCmd commands.Command, servers *upcloud.Servers, server upcloud.Server, details upcloud.ServerDetails, methodName string, madeRequest interface{}, requestResponse interface{}, args []string) {
	t.Helper()
	conf := config.New()
	mService := new(smock.Service)

	conf.Service = internal.Wrapper{Service: mService}
	mService.On("GetServers", mock.Anything).Return(servers, nil)
	mService.On("GetServerDetails", &request.GetServerDetailsRequest{UUID: server.UUID}).Return(&details, nil)
	mService.On(methodName, madeRequest).Return(requestResponse, nil)

	c := commands.BuildCommand(testCmd, nil, conf)
	err := c.Cobra().Flags().Parse(args)
	assert.NoError(t, err)

	_, err = c.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, mService, flume.New("test")), server.UUID)
	assert.NoError(t, err)

	mService.AssertNumberOfCalls(t, methodName, 1)
}
