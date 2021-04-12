package root

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/output"
	"runtime"
)

// VersionCommand reports the current version of upctl
type VersionCommand struct {
	*commands.BaseCommand
}

// Execute implements command.Command
func (s *VersionCommand) Execute(_ commands.Executor, _ string) (output.Output, error) {
	return output.Details{Sections: []output.DetailSection{
		{Rows: []output.DetailRow{
			{Title: "upctl version", Key: "version", Value: config.Version},
			{Title: "Build date:", Key: "build_date", Value: config.BuildDate},
			{Title: "Built with:", Key: "built_with", Value: runtime.Version()},
		},
		}},
	}, nil
}
