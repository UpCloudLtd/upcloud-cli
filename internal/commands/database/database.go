package database

import (
	"context"
	"fmt"
	"time"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// BaseDatabaseCommand creates the base "database" command
func BaseDatabaseCommand() commands.Command {
	return &databaseCommand{
		commands.New("database", "Manage databases"),
	}
}

type databaseCommand struct {
	*commands.BaseCommand
}

// InitCommand implements Command.InitCommand
func (db *databaseCommand) InitCommand() {
	db.Cobra().Aliases = []string{"db"}
}

// waitForManagedDatabaseState waits for database to reach given state and updates progress message with key matching given msg. Finally, progress message is updated back to given msg and either done state or timeout warning.
func waitForManagedDatabaseState(uuid string, state upcloud.ManagedDatabaseState, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for database %s to be in %s state", uuid, state))

	ctx, cancel := context.WithTimeout(exec.Context(), 15*time.Minute)
	defer cancel()

	if _, err := exec.All().WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
		UUID:         uuid,
		DesiredState: state,
	}); err != nil {
		exec.PushProgressUpdate(messages.Update{
			Key:     msg,
			Message: msg,
			Status:  messages.MessageStatusWarning,
			Details: "Error: " + err.Error(),
		})
		return
	}

	exec.PushProgressUpdateMessage(msg, msg)
	exec.PushProgressSuccess(msg)
}
