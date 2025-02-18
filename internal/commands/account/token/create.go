package token

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
)

// CreateCommand creates the "token create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create an API token",
			`upctl account token create --name test --expires-in 1h`,
			`upctl account token create --name test --expires-in 1h --allow-ip-range "0.0.0.0/0" --allow-ip-range "::/0"`,
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
	name            string
	expiresIn       time.Duration
	expiresAt       string
	canCreateTokens bool
	allowedIPRanges []string
}

func (s *createParams) processParams() error {
	if s.expiresIn == 0 && s.expiresAt == "" {
		return fmt.Errorf("either expires-in or expires-at must be set")
	}
	if s.expiresAt != "" {
		var err error
		s.ExpiresAt, err = time.Parse(time.RFC3339, s.expiresAt)
		if err != nil {
			return fmt.Errorf("invalid expires-at: %w", err)
		}
	} else {
		s.ExpiresAt = time.Now().Add(s.expiresIn)
	}
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
	fs.StringVar(&dst.expiresAt, "expires-at", def.expiresAt, "Exact time when the token expires in RFC3339 format. e.g. 2025-01-01T00:00:00Z")
	fs.DurationVar(&dst.expiresIn, "expires-in", def.expiresIn, "Duration until the token expires. e.g. 24h")
	fs.BoolVar(&dst.canCreateTokens, "can-create-tokens", def.canCreateTokens, "Allow token to be used to create further tokens.")
	fs.StringArrayVar(&dst.allowedIPRanges, "allow-ip-range", def.allowedIPRanges, "Allowed IP ranges for the token. If not defined, the token can not be used from any IP. To allow access from all IPs, use `0.0.0.0/0` as the value.")

	commands.Must(fs.SetAnnotation("name", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("expires-at", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("expires-in", commands.FlagAnnotationNoFileCompletions, nil))
	commands.Must(fs.SetAnnotation("allow-ip-range", commands.FlagAnnotationNoFileCompletions, nil))
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	s.flagSet = &pflag.FlagSet{}
	s.params = newCreateParams()
	applyCreateFlags(s.flagSet, &s.params, defaultCreateParams)

	s.AddFlags(s.flagSet)
	commands.Must(s.Cobra().MarkFlagRequired("name"))
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
