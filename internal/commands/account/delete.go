package account

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// DeleteCommand creates the "account delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete a sub-account",
			"upctl account delete my-sub-account",
		),
	}
}

type deleteCommand struct {
	*commands.BaseCommand
	resolver.CachingAccount
	completion.Account
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	svc := exec.All()
	msg := fmt.Sprintf("Deleting sub-account %v", arg)
	exec.PushProgressStarted(msg)

	err := svc.DeleteSubaccount(exec.Context(), &request.DeleteSubaccountRequest{
		Username: arg,
	})
	if err != nil {
		return commands.HandleError(exec, fmt.Sprintf("%s: failed", msg), err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
