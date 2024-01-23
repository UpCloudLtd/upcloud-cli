package account

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ListCommand creates the 'account list' command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List sub-accounts", "upctl account list"),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (l *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.Account()
	accounts, err := svc.GetAccountList(exec.Context())
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, a := range accounts {
		rows = append(rows, output.TableRow{
			a.Username,
			a.Type,
			a.Roles.Role,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: accounts,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "username", Header: "Username"},
				{Key: "type", Header: "Type"},
				{Key: "roles", Header: "Roles", Format: formatRoles},
			},
			Rows: rows,
		},
	}, nil
}

func formatRoles(val interface{}) (text.Colors, string, error) {
	roles, ok := val.([]string)
	if !ok {
		return nil, "", fmt.Errorf("cannot parse %T, expected []string", val)
	}

	return nil, strings.Join(roles, ", "), nil
}
