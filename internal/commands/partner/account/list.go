package partneraccount

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
)

func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New(
			"list",
			"List accounts associated with partner",
			"upctl partner account list",
		),
	}
}

type listCommand struct {
	*commands.BaseCommand
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (l *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	accounts, err := exec.All().GetPartnerAccounts(exec.Context())
	if err != nil {
		return nil, err
	}

	rows := []output.TableRow{}
	for _, a := range accounts {
		rows = append(rows, output.TableRow{
			a.Username,
			a.FirstName,
			a.LastName,
			a.Company,
		})
	}

	return output.MarshaledWithHumanOutput{
		Value: accounts,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "username", Header: "Username"},
				{Key: "first_name", Header: "First name"},
				{Key: "last_name", Header: "Last name"},
				{Key: "company", Header: "Company"},
			},
			Rows: rows,
		},
	}, nil
}
