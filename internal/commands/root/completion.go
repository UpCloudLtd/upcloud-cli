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

// Execute implements command.Command
func (s *CompletionCommand) Execute(_ commands.Executor, arg string) (output.Output, error) {
	if arg == "" {
		return nil, fmt.Errorf("shell name is requred")
	}
	if arg == "bash" {
		completion := new(bytes.Buffer)
		err := s.Cobra().Root().GenBashCompletion(completion)
		return output.Raw(completion.Bytes()), err
	}

	return nil, fmt.Errorf("completion for %s is not supported", arg)
}
