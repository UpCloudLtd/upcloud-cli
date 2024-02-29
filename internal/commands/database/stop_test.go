package database

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestStopCommand(t *testing.T) {
	targetMethod := "ShutdownManagedDatabase"

	db := upcloud.ManagedDatabase{
		State: upcloud.ManagedDatabaseStateRunning,
		Title: "database-title",
		UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
	}

	req := request.ShutdownManagedDatabaseRequest{
		UUID: db.UUID,
	}

	conf := config.New()
	testCmd := StopCommand()
	mService := new(smock.Service)

	mService.On(targetMethod, &req).Return(&db, nil)

	command := commands.BuildCommand(testCmd, nil, conf)

	command.Cobra().SetArgs([]string{db.UUID})
	_, err := mockexecute.MockExecute(command, mService, conf)

	assert.NoError(t, err)
	mService.AssertNumberOfCalls(t, targetMethod, 1)
}
