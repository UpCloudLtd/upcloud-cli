package tokens

import (
	"time"

	"fmt"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// CreateCommand creates the "tokens create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create an API token",
			`upctl tokens create --name test --expires_in 1h`,
			`upctl tokens create --name test --expires_in 1h --allowed-ip-ranges "0.0.0.0/0" --allowed-ip-ranges "::/0"`,
		),
	}
}

var defaultCreateParams = &createParams{
	CreateTokenRequest: request.CreateTokenRequest{},
	name:               "",
	expiresIn:          0,
	allowedIPRanges:    []string{}, // TODO: should we default to empty or "0.0.0.0/0", "::/0"?
	canCreateTokens:    false,
}

func newCreateParams() createParams {
	return createParams{
		CreateTokenRequest: request.CreateTokenRequest{},
	}
}

type createParams struct {
	request.CreateTokenRequest
	name      string
	expiresIn time.Duration
	//	expiresAt       time.Time/string // TODO: is it necessary to be able to define exact time for expiry instead of duration?
	canCreateTokens bool
	allowedIPRanges []string
}

func (s *createParams) processParams() error {
	s.ExpiresAt = time.Now().Add(s.expiresIn)
	s.Name = s.name
	s.CanCreateSubTokens = s.canCreateTokens
	s.AllowedIPRanges = s.allowedIPRanges
	return nil
}

type createCommand struct {
	*commands.BaseCommand
	params  createParams
	flagSet *pflag.FlagSet
}

func applyCreateFlags(fs *pflag.FlagSet, dst, def *createParams) {
	fs.StringVar(&dst.name, "name", def.name, "Name for the token.")
	fs.DurationVar(&dst.expiresIn, "expires_in", def.expiresIn, "Duration until the token expires.")
	fs.BoolVar(&dst.canCreateTokens, "can-create-tokens", def.canCreateTokens, "Allow token to be used to create further tokens.")
	fs.StringArrayVar(&dst.allowedIPRanges, "allowed-ip-ranges", def.allowedIPRanges, "Allowed IP ranges for the token.")
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	s.flagSet = &pflag.FlagSet{}
	s.params = newCreateParams()
	applyCreateFlags(s.flagSet, &s.params, defaultCreateParams)

	s.AddFlags(s.flagSet)
	_ = s.Cobra().MarkFlagRequired("name")
	_ = s.Cobra().MarkFlagRequired("expires_in")
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.Token()

	if err := s.params.processParams(); err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Creating token %s", s.params.Name)
	exec.PushProgressStarted(msg)

	res, err := svc.CreateToken(exec.Context(), &s.params.CreateTokenRequest)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "API Token", Value: res.APIToken, Colour: ui.DefaultNoteColours},
		{Title: "Name", Value: res.Name},
		{Title: "ID", Value: res.ID, Colour: ui.DefaultUUUIDColours},
		{Title: "Type", Value: res.Type},
		{Title: "Created At", Value: res.CreatedAt.Format(time.RFC3339)},
		{Title: "Expires At", Value: res.ExpiresAt.Format(time.RFC3339)},
		{Title: "Allowed IP Ranges", Value: res.AllowedIPRanges},
		{Title: "Can Create Sub Tokens", Value: res.CanCreateSubTokens},
	}}, nil
}
