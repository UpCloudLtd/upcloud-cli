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

func MockExecute(command commands.Command, service service.AllServices, conf *config.Config) (string, error) {
	buf := bytes.NewBuffer(nil)
	command.Cobra().SetErr(buf)
	command.Cobra().SetOut(buf)

	// Use human output if nothing else is defined
	if !conf.IsSet(config.KeyOutput) {
		conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)
	}

	command.Cobra().RunE = func(_ *cobra.Command, args []string) error {
		return mockRunE(command, service, conf, args)
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
		// Panics if called with nil args
		out, err = typedCommand.ExecuteSingleArgument(executor, args[0])
	case commands.MultipleArgumentCommand:
		// Panics if called with nil args
		out, err = typedCommand.Execute(executor, args[0])
	}
	if err != nil {
		return err
	}
	return output.Render(command.Cobra().OutOrStdout(), config, out)
}
