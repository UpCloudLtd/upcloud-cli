package database

import (
	"context"
	"fmt"
	"time"

	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
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
func WaitForManagedDatabaseState(uuid string, state string, exec commands.Executor, msg string) {
	exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for database %s to be in %s state", uuid, state))

	ctx, cancel := context.WithTimeout(exec.Context(), 15*time.Minute)
	defer cancel()

	err := func() error {
		client, err := v9Client(exec)
		if err != nil {
			return err
		}
		id, err := parseDatabaseUUID(uuid)
		if err != nil {
			return err
		}
		return waitForDatabaseState(ctx, client, id, state)
	}()
	if err != nil {
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
