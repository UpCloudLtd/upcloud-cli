package root

import (
	"runtime"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

// VersionCommand reports the current version of upctl
type VersionCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *VersionCommand) ExecuteWithoutArguments(_ commands.Executor) (output.Output, error) {
	return output.Details{
		Sections: []output.DetailSection{
			{
				Rows: []output.DetailRow{
					{Title: "Version:", Key: "version", Value: config.GetVersion()},
					{Title: "Build date:", Key: "build_date", Value: config.BuildDate},
					{Title: "Built with:", Key: "built_with", Value: runtime.Version()},
					{Title: "System:", Key: "operating_system", Value: runtime.GOOS},
					{Title: "Architecture:", Key: "architecture", Value: runtime.GOARCH},
				},
			},
		},
	}, nil
}

// DoesNotUseServices implements commands.OfflineCommand as this command does not use services
func (s *VersionCommand) DoesNotUseServices() {}
