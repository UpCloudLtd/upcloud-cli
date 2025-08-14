package serverstorage

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/storage"
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

type attachCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
	params attachParams
}

type attachParams struct {
	request.AttachStorageRequest
	bootable config.OptionalBoolean
}

// AttachCommand creates the "server storage attach" command
func AttachCommand() commands.Command {
	return &attachCommand{
		BaseCommand: commands.New(
			"attach",
			"Attach a storage as a device to a server",
			"upctl server storage attach 00038afc-d526-4148-af0e-d2f1eeaded9b --storage 015899e0-0a68-4949-85bb-261a99de5fdd",
			"upctl server storage attach 00038afc-d526-4148-af0e-d2f1eeaded9b --storage 01a5568f-4766-4ce7-abf5-7d257903a735 --address virtio:2",
			`upctl server storage attach my_server1 --storage "My Storage"`,
		),
	}
}

var defaultAttachParams = &attachParams{
	AttachStorageRequest: request.AttachStorageRequest{
		Type:     "disk",
		BootDisk: 0,
		Address:  "virtio",
	},
}

// InitCommand implements Command.InitCommand
func (s *attachCommand) InitCommand() {
	types := []string{upcloud.StorageTypeDisk, upcloud.StorageTypeCDROM}
	s.params = attachParams{AttachStorageRequest: request.AttachStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Type, "type", defaultAttachParams.Type, "Type of the attached storage. Available: "+strings.Join(types, ", "))
	flagSet.StringVar(&s.params.Address, "address", defaultAttachParams.Address, "Address where the storage device is attached on the server. \nAddress is of the form busname:deviceindex where busname can be ide/scsi/virtio. (example: 'virtio:1')\nSpecify only the bus name to auto-select next available device index from that bus. (example: 'virtio')")
	flagSet.StringVar(&s.params.StorageUUID, "storage", defaultAttachParams.StorageUUID, "UUID of the storage to attach.")
	config.AddToggleFlag(flagSet, &s.params.bootable, "boot-disk", false, "Set attached device as the server's boot disk.")

	s.AddFlags(flagSet)
	commands.Must(s.Cobra().MarkFlagRequired("storage"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("type", cobra.FixedCompletions(types, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("address", cobra.NoFileCompletions))
}

func (s *attachCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("storage", namedargs.CompletionFunc(completion.StorageUUID{}, cfg)))
}

// MaximumExecutions implements command.Command
func (s *attachCommand) MaximumExecutions() int {
	return maxServerStorageActions
}

// ExecuteSingleArgument implements command.SingleArgumentCommand
func (s *attachCommand) ExecuteSingleArgument(exec commands.Executor, uuid string) (output.Output, error) {
	storageSvc := exec.Storage()

	strg, err := storage.SearchSingleStorage(s.params.StorageUUID, exec)
	if err != nil {
		return nil, err
	}

	s.params.StorageUUID = strg.UUID
	s.params.BootDisk = defaultAttachParams.BootDisk

	if s.params.bootable.Value() {
		s.params.BootDisk = 1
	}
	req := s.params.AttachStorageRequest
	req.ServerUUID = uuid

	msg := fmt.Sprintf("Attaching storage %q to server %q", req.StorageUUID, req.ServerUUID)
	exec.PushProgressStarted(msg)

	res, err := storageSvc.AttachStorage(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
