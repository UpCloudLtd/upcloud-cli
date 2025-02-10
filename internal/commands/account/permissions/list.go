package permissions

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/pflag"
)

// ListCommand creates the 'permissions list' command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List permissions", "upctl account permissions list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
	username string
}

// InitCommand implements Command.InitCommand
func (l *listCommand) InitCommand() {
	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&l.username, "username", "", "Filter permissions by username.")

	l.AddFlags(flagSet)
}

func (l *listCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(l.Cobra().RegisterFlagCompletionFunc("username", namedargs.CompletionFunc(completion.Username{}, cfg)))
}

func (l *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	permissions, err := svc.GetPermissions(exec.Context(), &request.GetPermissionsRequest{})
	if err != nil {
		return nil, err
	}

	filtered := make([]upcloud.Permission, 0)
	rows := []output.TableRow{}
	for _, permission := range permissions {
		if permission.User == l.username || l.username == "" {
			filtered = append(filtered, permission)
			rows = append(rows, output.TableRow{
				permission.User,
				permission.TargetType,
				permission.TargetIdentifier,
				permission.Options,
			})
		}
	}
	return output.MarshaledWithHumanOutput{
		Value: filtered,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "username", Header: "Username"},
				{Key: "target_type", Header: "Target type"},
				{Key: "target_identifier", Header: "Target identifier", Colour: ui.DefaultUUUIDColours},
				{Key: "options", Header: "Options", Format: formatOptions},
			},
			Rows: rows,
		},
	}, nil
}

func formatOptions(val interface{}) (text.Colors, string, error) {
	options, ok := val.(*upcloud.PermissionOptions)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse %T, expected *upcloud.PermissionOptions", val)
	}

	if options == nil {
		return nil, "", nil
	}

	colors, value, _ := format.Boolean(options.Storage)
	return nil, fmt.Sprintf("Storage: %s", colors.Sprint(value)), nil
}
