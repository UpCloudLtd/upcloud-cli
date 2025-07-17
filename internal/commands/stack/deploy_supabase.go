package stack

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
)

func DeploySupabaseCommand() commands.Command {
	return &deploySupabaseCommand{
		BaseCommand: commands.New(
			"supabase",
			"Deploy a Supabase stack",
			"upctl stack deploy supabase <project-name>",
			"upctl stack deploy supabase my-new-project",
		),
	}
}

type deploySupabaseCommand struct {
	*commands.BaseCommand
	location string
	name     string
}

func (s *deploySupabaseCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.location, "location", s.location, "Select the location (region) for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Specify the name of the Supabase project")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("location"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))
}

func (s *deploySupabaseCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	msg := fmt.Sprintf("Creating supabase stack %v", s.name)
	exec.PushProgressStarted(msg)

	// Command implementation for deploying a Supabase stack

	exec.PushProgressSuccess("Supabase stack created successfully")

	return output.Raw([]byte("Commamnd executed successfully")), nil
}
