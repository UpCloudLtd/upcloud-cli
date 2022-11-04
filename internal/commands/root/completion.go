package root

import (
	"bytes"
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"
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

// DoesNotUseServices implements commands.OfflineCommand as this command does not use services
func (s *CompletionCommand) DoesNotUseServices() {}
