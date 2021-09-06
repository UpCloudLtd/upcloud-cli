package root

import (
	"runtime"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
)

// VersionCommand reports the current version of upctl.
type VersionCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand.
func (s *VersionCommand) ExecuteWithoutArguments(_ commands.Executor) (output.Output, error) {
	return output.Details{Sections: []output.DetailSection{
		{Rows: []output.DetailRow{
			{Title: "Version", Key: "version", Value: config.Version},
			{Title: "Build date:", Key: "build_date", Value: config.BuildDate},
			{Title: "Built with:", Key: "built_with", Value: runtime.Version()},
		},
		}},
	}, nil
}
