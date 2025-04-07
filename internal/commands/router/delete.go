package router

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingRouter
	completion.Router
}

// DeleteCommand creates the "delete router" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a router",
			"upctl router delete 0497728e-76ef-41d0-997f-fa9449eb71bc",
			"upctl router delete my_router",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
}

// MaximumExecutions implements Command.MaximumExecutions
func (s *deleteCommand) MaximumExecutions() int {
	return maxRouterActions
}

func Delete(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Network()
	msg := fmt.Sprintf("Deleting router %s", uuid)
	exec.PushProgressStarted(msg)

	err := svc.DeleteRouter(exec.Context(), &request.DeleteRouterRequest{UUID: uuid})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	return Delete(exec, arg)
}
