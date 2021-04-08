package root

import (
	"bytes"
	"fmt"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/output"
)

// CompletionCommand creates shell completion scripts
type CompletionCommand struct {
	*commands.BaseCommand
}

// Execute implements command.NewCommand
func (s *CompletionCommand) Execute(exec commands.Executor, arg string) (output.Output, error) {
	if arg == "" {
		return nil, fmt.Errorf("shell name is requred")
	}
	completion := new(bytes.Buffer)
	if arg == "bash" {
		err := s.Cobra().Root().GenBashCompletion(completion)
		return output.Raw(completion.Bytes()), err
	}

	return nil, fmt.Errorf("completion for %s is not supported", arg)
}
