package token

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/paging"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"

	"github.com/spf13/pflag"
)

// ListCommand creates the 'token list' command
func ListCommand() commands.Command {
	return &listCommand{
		BaseCommand: commands.New("list", "List API tokens", "upctl account token list"),
	}
}

func (l *listCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	l.ConfigureFlags(fs)
	l.AddFlags(fs)
}

type listCommand struct {
	*commands.BaseCommand
	paging.PageParameters
}

func (l *listCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()
	tokens, err := svc.GetTokens(exec.Context(), &request.GetTokensRequest{
		Page: l.Page(),
	})
	if err != nil {
		return nil, err
	}

	var rows []output.TableRow
	for _, token := range *tokens {
		rows = append(rows, output.TableRow{
			token.ID,
			token.Name,
			token.Type,
			token.LastUsed,
			token.ExpiresAt,
		})
	}
	return output.MarshaledWithHumanOutput{
		Value: tokens,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "id", Header: "UUID", Colour: ui.DefaultUUUIDColours},
				{Key: "name", Header: "Name"},
				{Key: "type", Header: "Type"},
				{Key: "last_used", Header: "Last Used"},
				{Key: "expires_at", Header: "Expires At"},
			},
			Rows: rows,
		},
	}, nil
}
