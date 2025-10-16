package supabase

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
)

func DestroySupabaseCommand() commands.Command {
	return &destroySupabaseCommand{
		BaseCommand: commands.New(
			"supabase",
			"Destroy a Supabase stack",
			"upctl stack destroy supabase --name <project-name> --zone <zone-name>",
			"upctl stack destroy supabase --name my-new-project --zone es-mad1",
		),
	}
}

type destroySupabaseCommand struct {
	*commands.BaseCommand
	zone                string
	name                string
	deleteStorage       bool
	deleteObjectStorage bool
}

func (s *destroySupabaseCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.zone, "zone", s.zone, "Zone for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Supabase stack name")
	fs.BoolVar(&s.deleteStorage, "delete-storage", false, "Delete associated UpCloud storage resources")
	fs.BoolVar(&s.deleteObjectStorage, "delete-object-storage", false, "Delete associated UpCloud object storage resources")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *destroySupabaseCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	err := stack.DestroyStack(exec, c.name, c.zone, c.deleteStorage, c.deleteObjectStorage, stack.StackTypeSupabase)
	if err != nil {
		return nil, err
	}

	return output.None{}, nil
}
