package databaseindex

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type deleteCommand struct {
	*commands.BaseCommand
	name string
	completion.Database
	resolver.CachingDatabase
}

// DeleteCommand creates the "database index delete" command
func DeleteCommand() commands.Command {
	return &deleteCommand{
		BaseCommand: commands.New(
			"delete",
			"Delete an index from the specified database.",
			"upctl database index delete 55199a44-4751-4e27-9394-7c7661910be3 --name .index-to-delete",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *deleteCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.name, "name", "", "Index name")
	s.AddFlags(flagSet)

	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("name", cobra.NoFileCompletions))
}

// Execute implements commands.MultipleArgumentCommand
func (s *deleteCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	msg := fmt.Sprintf("Deleting index %s from Managed Database %v", s.name, arg)
	exec.PushProgressStarted(msg)

	err := exec.All().DeleteManagedDatabaseIndex(exec.Context(), &request.DeleteManagedDatabaseIndexRequest{
		ServiceUUID: arg, IndexName: s.name,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
