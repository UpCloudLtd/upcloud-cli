package database

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	upcloudv9 "github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"
	"github.com/google/uuid"
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

func Delete(exec commands.Executor, uuidStr string, disableTerminationProtection, wait bool) (output.Output, error) {
	msg := fmt.Sprintf("Deleting database %s", uuidStr)
	exec.PushProgressStarted(msg)

	client, err := v9Client(exec)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}
	id, err := parseDatabaseUUID(uuidStr)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if disableTerminationProtection {
		b := false
		res, err := client.ModifyDatabaseWithResponse(exec.Context(), id, upcloudv9.ModifyDatabaseJSONRequestBody{
			TerminationProtection: &b,
		})
		if err != nil {
			return commands.HandleError(exec, msg, err)
		}
		if res.JSON200 == nil {
			return commands.HandleError(exec, msg, apierror.FromResponse(res.StatusCode(), res.Body))
		}
	}

	res, err := client.DeleteDatabaseWithResponse(exec.Context(), id)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}
	if status := res.StatusCode(); status < 200 || status > 299 {
		return commands.HandleError(exec, msg, apierror.FromResponse(status, res.Body))
	}

	if wait {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for database service %s to be deleted", uuidStr))
		err = waitUntilDatabaseDeleted(exec, client, id)
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

func waitUntilDatabaseDeleted(exec commands.Executor, client *upcloudv9.ClientWithResponses, id uuid.UUID) error {
	ticker := time.NewTicker(databasePollInterval)
	defer ticker.Stop()

	ctx := exec.Context()

	for {
		select {
		case <-ticker.C:
			res, err := client.GetDatabaseWithResponse(ctx, id)
			if err != nil {
				return err
			}
			if res.JSON200 == nil {
				err := apierror.FromResponse(res.StatusCode(), res.Body)
				if apierror.IsNotFound(err) {
					return nil
				}

				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
