package storage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type templatizeCommand struct {
	*commands.BaseCommand
	wait config.OptionalBoolean
	resolver.CachingStorage
	completion.Storage
	params templatizeParams
}

type templatizeParams struct {
	request.TemplatizeStorageRequest
}

// TemplatizeCommand creates the "storage templatise" command
// TODO: figure out consistent naming, one way or the other.
func TemplatizeCommand() commands.Command {
	return &templatizeCommand{
		BaseCommand: commands.New(
			"templatise",
			"Templatise a storage",
			`upctl storage templatise 01271548-2e92-44bb-9774-d282508cc762 --title "My Template"`,
			`upctl storage templatise 01271548-2e92-44bb-9774-d282508cc762 --title "My Template" --wait`,
			`upctl storage templatise "My Storage" --title super_template`,
		),
	}
}

var defaultTemplatizeParams = &templatizeParams{
	TemplatizeStorageRequest: request.TemplatizeStorageRequest{},
}

// InitCommand implements Command.InitCommand
func (s *templatizeCommand) InitCommand() {
	s.params = templatizeParams{TemplatizeStorageRequest: request.TemplatizeStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Title, "title", defaultTemplatizeParams.Title, "A short, informational description.")
	config.AddToggleFlag(flagSet, &s.wait, "wait", false, "Wait for storage to be in online state before returning.")

	s.AddFlags(flagSet)
	commands.Must(s.Cobra().MarkFlagRequired("title"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("title", cobra.NoFileCompletions))
}

// Execute implements commands.MultipleArgumentCommand
func (s *templatizeCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()
	req := s.params.TemplatizeStorageRequest
	req.UUID = uuid

	msg := fmt.Sprintf("Templatise storage %v", uuid)
	exec.PushProgressStarted(msg)

	res, err := svc.TemplatizeStorage(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	if s.wait.Value() {
		waitForStorageState(res.UUID, upcloud.StorageStateOnline, exec, msg)
	} else {
		exec.PushProgressSuccess(msg)
	}

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
