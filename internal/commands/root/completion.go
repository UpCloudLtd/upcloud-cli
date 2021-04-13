package root

import (
	"bytes"
	"fmt"
	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
)

// CompletionCommand creates shell completion scripts
type CompletionCommand struct {
	*commands.BaseCommand
	resolver.CompletionResolver
}

// ExecuteSingleArgument implements commands.SingleArgumentCommand
func (s *CompletionCommand) ExecuteSingleArgument(_ commands.Executor, arg string) (output.Output, error) {
	if arg == "bash" {
		completion := new(bytes.Buffer)
		err := s.Cobra().Root().GenBashCompletion(completion)

		return output.Raw(completion.Bytes()), err
	}

	return nil, fmt.Errorf("completion for %s is not supported", arg)
}
