package database

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
)

func TestStartCommand(t *testing.T) {
	targetMethod := "StartManagedDatabase"

	var db = upcloud.ManagedDatabase{
		State: upcloud.ManagedDatabaseStateRunning,
		Title: "database-title",
		UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
	}

	req := request.StartManagedDatabaseRequest{
		UUID: db.UUID,
	}

	conf := config.New()
	testCmd := StartCommand()
	mService := new(smock.Service)

	conf.Service = internal.Wrapper{Service: mService}
	mService.On(targetMethod, &req).Return(&db, nil)

	command := commands.BuildCommand(testCmd, nil, conf)

	command.Cobra().SetArgs([]string{db.UUID})
	_, err := mockexecute.MockExecute(command, mService, conf)

	assert.NoError(t, err)
	mService.AssertNumberOfCalls(t, targetMethod, 1)
}
