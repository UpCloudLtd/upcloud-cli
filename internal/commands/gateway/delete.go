package gateway

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/utils"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "gateway delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a gateway",
			"upctl gateway delete 8abc8009-4325-4b23-4321-b1232cd81231",
			"upctl gateway delete my-gateway",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingGateway
	completion.Gateway

	wait config.OptionalBoolean
}

// InitCommand implements Command.InitCommand
func (c *deleteCommand) InitCommand() {
	flags := &pflag.FlagSet{}
	config.AddToggleFlag(flags, &c.wait, "wait", false, "Wait until the gateway has been deleted before returning.")
	c.AddFlags(flags)
}

func Delete(exec commands.Executor, uuid string, wait bool) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting gateway %v", uuid)
	exec.PushProgressStarted(msg)

	err := svc.DeleteGateway(exec.Context(), &request.DeleteGatewayRequest{
		UUID: uuid,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if wait {
		exec.PushProgressUpdateMessage(msg, fmt.Sprintf("Waiting for gateway %s to be deleted", uuid))
		err = waitUntilGatewayDeleted(exec, uuid)
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
	return Delete(exec, arg, c.wait.Value())
}

func waitUntilGatewayDeleted(exec commands.Executor, uuid string) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx := exec.Context()
	svc := exec.All()

	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			_, err := svc.GetGateway(exec.Context(), &request.GetGatewayRequest{
				UUID: uuid,
			})
			if err != nil {
				if utils.IsNotFoundError(err) {
					return nil
				}

				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
