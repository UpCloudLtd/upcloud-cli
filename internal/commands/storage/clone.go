package storage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type cloneCommand struct {
	*commands.BaseCommand
	resolver.CachingStorage
	completion.Storage
	params cloneParams
}

type cloneParams struct {
	request.CloneStorageRequest
	encrypted config.OptionalBoolean
}

// CloneCommand creates the "storage clone" command
func CloneCommand() commands.Command {
	return &cloneCommand{
		BaseCommand: commands.New(
			"clone",
			"Clone a storage",
			"upctl storage clone 015899e0-0a68-4949-85bb-261a99de5fdd --title my_storage_clone --zone fi-hel1",
			"upctl storage clone 015899e0-0a68-4949-85bb-261a99de5fdd --title my_storage_clone2 --zone pl-waw1  --tier maxiops",
			`upctl storage clone "My Storage" --title my_storage_clone3 --zone pl-waw1  --tier maxiops`,
		),
	}
}

var defaultCloneParams = &cloneParams{
	CloneStorageRequest: request.CloneStorageRequest{
		Tier: upcloud.StorageTierHDD,
	},
}

// InitCommand implements Command.InitCommand
func (s *cloneCommand) InitCommand() {
	s.params = cloneParams{CloneStorageRequest: request.CloneStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Tier, "tier", defaultCloneParams.Tier, "The storage tier to use.")
	flagSet.StringVar(&s.params.Title, "title", defaultCloneParams.Title, "A short, informational description.")
	flagSet.StringVar(&s.params.Zone, "zone", defaultCloneParams.Zone, namedargs.ZoneDescription("storage"))
	config.AddToggleFlag(flagSet, &s.params.encrypted, "encrypt", false, "Encrypt the new storage.")

	s.AddFlags(flagSet)
	commands.Must(s.Cobra().MarkFlagRequired("title"))
	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("tier", cobra.FixedCompletions(tiers, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("title", cobra.NoFileCompletions))
}

func (s *cloneCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("zone", namedargs.CompletionFunc(completion.Zone{}, cfg)))
}

// Execute implements commands.MultipleArgumentCommand
func (s *cloneCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Storage()
	req := s.params.CloneStorageRequest
	req.UUID = uuid
	req.Encrypted = s.params.encrypted.AsUpcloudBoolean()

	msg := fmt.Sprintf("Cloning storage %v", uuid)
	exec.PushProgressStarted(msg)

	res, err := svc.CloneStorage(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
