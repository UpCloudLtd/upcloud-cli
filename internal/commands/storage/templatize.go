package storage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-cli/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/spf13/pflag"
)

type templatizeCommand struct {
	*commands.BaseCommand
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

	s.AddFlags(flagSet)
}

// Execute implements commands.MultipleArgumentCommand
func (s *templatizeCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	if s.params.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	svc := exec.Storage()
	req := s.params.TemplatizeStorageRequest
	req.UUID = uuid

	msg := fmt.Sprintf("Templatise storage %v", uuid)
	logline := exec.NewLogEntry(msg)

	logline.StartedNow()

	res, err := svc.TemplatizeStorage(&req)
	if err != nil {
		return commands.HandleError(logline, fmt.Sprintf("%s: failed", msg), err)
	}

	logline.SetMessage(fmt.Sprintf("%s: success", msg))
	logline.MarkDone()

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
