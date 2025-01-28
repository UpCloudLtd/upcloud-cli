package tokens

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// ShowCommand creates the "token show" command
func ShowCommand() commands.Command {
	return &showCommand{
		BaseCommand: commands.New(
			"show",
			"Show API token details",
			"upctl token show 0cdabbf9-090b-4fc5-a6ae-3f76801ed171",
		),
	}
}

type showCommand struct {
	*commands.BaseCommand
}

// Execute implements commands.MultipleArgumentCommand
func (c *showCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Token()
	token, err := svc.GetTokenDetails(exec.Context(), &request.GetTokenDetailsRequest{ID: uuid})
	if err != nil {
		return nil, err
	}

	details := output.Details{
		Sections: []output.DetailSection{
			{
				Rows: []output.DetailRow{
					{Title: "Name:", Key: "name", Value: token.Name},
					{Title: "Token ID", Key: "id", Colour: ui.DefaultUUUIDColours, Value: token.ID},
					{Title: "Type", Key: "type", Value: token.Type},
					{Title: "Created At", Key: "created_at", Value: token.CreatedAt},
					{Title: "Last Used", Key: "last_used", Value: token.LastUsed},
					{Title: "Expires At", Key: "expires_at", Value: token.ExpiresAt},
					{Title: "Allowed IP Ranges", Key: "allowed_ip_ranges", Value: token.AllowedIPRanges},
					{Title: "Can Create Sub Tokens", Key: "can_create_sub_tokens", Value: token.CanCreateSubTokens},
				},
			},
		},
	}
	return output.MarshaledWithHumanOutput{
		Value:  token,
		Output: details,
	}, nil
}
