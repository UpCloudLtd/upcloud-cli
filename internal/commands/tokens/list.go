package tokens

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
		BaseCommand: commands.New("list", "List API tokens", "upctl token list"),
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
			token.Name,
			token.ID,
			token.Type,
			token.CreatedAt,
			token.LastUsed,
			token.ExpiresAt,
			token.AllowedIPRanges,
			token.CanCreateSubTokens,
		})
	}
	return output.MarshaledWithHumanOutput{
		Value: tokens,
		Output: output.Table{
			Columns: []output.TableColumn{
				{Key: "name", Header: "Name"},
				{Key: "id", Header: "Token ID", Colour: ui.DefaultUUUIDColours},
				{Key: "type", Header: "Type"},
				{Key: "created_at", Header: "Created At"},
				{Key: "last_used", Header: "Last Used"},
				{Key: "expires_at", Header: "Expires At"},
				{Key: "allowed_ip_ranges", Header: "Allowed IP Ranges"},
				{Key: "can_create_sub_tokens", Header: "Can Create Sub Tokens"},
			},
			Rows: rows,
		},
	}, nil
}
