package mockexecute

import (
	"bytes"

	"github.com/gemalto/flume"
	"github.com/spf13/cobra"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/service"
)

func MockExecute(command commands.Command, service service.AllServices, config *config.Config) (string, error) {
	buf := bytes.NewBuffer(nil)
	command.Cobra().SetErr(buf)
	command.Cobra().SetOut(buf)

	command.Cobra().RunE = func(_ *cobra.Command, args []string) error {
		return mockRunE(command, service, config, args)
	}
	err := command.Cobra().Execute()

	return buf.String(), err
}

func mockRunE(command commands.Command, service service.AllServices, config *config.Config, args []string) error {
	executor := commands.NewExecutor(config, service, flume.New("test"))

	var err error
	var out output.Output
	switch typedCommand := command.(type) {
	case commands.NoArgumentCommand:
		out, err = typedCommand.ExecuteWithoutArguments(executor)
	case commands.SingleArgumentCommand:
		out, err = typedCommand.ExecuteSingleArgument(executor, args[0])
	case commands.MultipleArgumentCommand:
		out, err = typedCommand.Execute(executor, args[0])
	}

	_ = output.Render(command.Cobra().OutOrStdout(), config, out)

	return err
}
