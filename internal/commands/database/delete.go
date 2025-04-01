package database

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingDatabase
	completion.Database

	disableTerminationProtection config.OptionalBoolean
	wait                         config.OptionalBoolean
}

// DeleteCommand creates the "delete database" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a database",
			"upctl database delete 0497728e-76ef-41d0-997f-fa9449eb71bc",
			"upctl database delete my_database",
		),
	}
}

// InitCommand implements Command.InitCommand
func (c *deleteCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &c.disableTerminationProtection, "disable-termination-protection", false, "Disable termination-protection before deleting the database instance.")
	config.AddToggleFlag(flags, &c.wait, "wait", false, "Wait until the database instance has been deleted before returning.")
	c.AddFlags(flags)
}

func Delete(exec commands.Executor, uuid string, disableTerminationProtection, wait bool) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting database %s", uuid)
	exec.PushProgressStarted(msg)

	if disableTerminationProtection {
		b := false
		_, err := svc.ModifyManagedDatabase(exec.Context(), &request.ModifyManagedDatabaseRequest{
			UUID:                  uuid,
			TerminationProtection: &b,
		})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
	}

	err := svc.DeleteManagedDatabase(exec.Context(), &request.DeleteManagedDatabaseRequest{UUID: uuid})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if wait {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for database service %s to be deleted", uuid))
		err = waitUntilDatabaseDeleted(exec, uuid)
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
		exec.PushProgressUpdateMessage(msg, msg)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (c *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	return Delete(exec, arg, c.disableTerminationProtection.Value(), c.wait.Value())
}

func waitUntilDatabaseDeleted(exec commands.Executor, uuid string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx := exec.Context()
	svc := exec.All()

	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			_, err := svc.GetManagedDatabase(exec.Context(), &request.GetManagedDatabaseRequest{
				UUID: uuid,
			})
			if err != nil {
				var ucErr *upcloud.Problem
				if errors.As(err, &ucErr) && ucErr.Status == http.StatusNotFound {
					return nil
				}

				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
