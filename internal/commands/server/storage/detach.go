package serverstorage

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type detachCommand struct {
	*commands.BaseCommand
	resolver.CachingServer
	completion.Server
	params detachParams
}

type detachParams struct {
	request.DetachStorageRequest
}

var defaultDetachParams = &detachParams{
	DetachStorageRequest: request.DetachStorageRequest{},
}

// DetachCommand creates the "server storage detach" command
func DetachCommand() commands.Command {
	return &detachCommand{
		BaseCommand: commands.New(
			"detach",
			"Detaches a storage resource from a server",
			"upctl server storage detach 00038afc-d526-4148-af0e-d2f1eeaded9b --address virtio:1",
			"upctl server storage detach my_server1 --address virtio:2",
		),
	}
}

// InitCommand implements Command.InitCommand
func (s *detachCommand) InitCommand() {
	s.params = detachParams{DetachStorageRequest: request.DetachStorageRequest{}}

	flagSet := &pflag.FlagSet{}
	flagSet.StringVar(&s.params.Address, "address", defaultDetachParams.Address, "Detach the storage attached to this address.")

	s.AddFlags(flagSet)
	commands.Must(s.Cobra().MarkFlagRequired("address"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("address", cobra.NoFileCompletions))
}

// MaximumExecutions implements command.Command
func (s *detachCommand) MaximumExecutions() int {
	return maxServerStorageActions
}

// ExecuteSingleArgument implements command.SingleArgumentCommand
func (s *detachCommand) ExecuteSingleArgument(exec commands.Executor, uuid string) (output.Output, error) {
	storageSvc := exec.Storage()

	req := s.params.DetachStorageRequest
	req.ServerUUID = uuid

	msg := fmt.Sprintf("Detaching storage address %q from server %q", req.Address, req.ServerUUID)
	exec.PushProgressStarted(msg)

	res, err := storageSvc.DetachStorage(exec.Context(), &req)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
