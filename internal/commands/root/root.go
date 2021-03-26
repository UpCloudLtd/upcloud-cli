package root

import (
	"fmt"
	"os"
	"runtime"

	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
)

type CompletionCommand struct {
	*commands.BaseCommand
}

func (s *CompletionCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("shell name is requred")
		}
		shellName := args[0]

		if shellName == "bash" {
			err := s.Cobra().Root().GenBashCompletion(os.Stdout)
			return nil, err
		}

		return nil, fmt.Errorf("completion for %s is not supported", shellName)
	}
}

type VersionCommand struct {
	*commands.BaseCommand
}

func (s *VersionCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {
		return fmt.Printf(
			"Upctl %v\n\tBuild date: %v\n\tBuilt with: %v",
			config.Version, config.BuildDate, runtime.Version(),
		)
	}
}
