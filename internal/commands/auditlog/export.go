package auditlog

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

var formats = []string{
	request.ExportAuditLogFormatJSON,
	request.ExportAuditLogFormatCSV,
}

// ExportCommand creates the "audit-log export" command
func ExportCommand() commands.Command {
	return &exportCommand{
		BaseCommand: commands.New("export", "Export audit logs", "upctl audit-log export", "upctl audit-log export --output csv >audit-log.csv"),
	}
}

type exportCommand struct {
	*commands.BaseCommand
	params request.ExportAuditLogRequest
}

func (c *exportCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&c.params.Format, "output", formats[0], fmt.Sprintf("Export format (typically %s). Note that this overrides the global --output flag.", strings.Join(formats, "|")))
	c.AddFlags(fs)
}

// ExecuteWithoutArguments implements [commands.NoArgumentCommand].
func (c *exportCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	r, err := exec.All().ExportAuditLog(exec.Context(), &c.params)
	if err != nil {
		return nil, err
	}

	return output.Raw{Source: r}, nil
}
