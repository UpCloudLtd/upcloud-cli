package dokku

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
)

func DestroyDokkuCommand() commands.Command {
	return &destroyDokkuCommand{
		BaseCommand: commands.New(
			"dokku",
			"Destroy a Dokku stack",
			"upctl stack destroy dokku --name <project-name> --zone <zone-name>",
			"upctl stack destroy dokku --name my-new-project --zone es-mad1",
		),
	}
}

type destroyDokkuCommand struct {
	*commands.BaseCommand
	zone string
	name string
}

func (s *destroyDokkuCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.zone, "zone", s.zone, "Zone for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Supabase stack name")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
// TODO: This is not deleting the LB that dokku creates. Need to find a way to identify it.
func (c *destroyDokkuCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	err := stack.DestroyStack(exec, c.name, c.zone, false, false, stack.StackTypeDokku)
	if err != nil {
		return nil, err
	}

	return output.None{}, nil
}
