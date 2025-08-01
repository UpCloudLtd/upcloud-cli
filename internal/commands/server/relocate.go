package server

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/spf13/pflag"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// RelocateCommand creates the "server relocate" command
func RelocateCommand() commands.Command {
	return &relocateCommand{
		BaseCommand: commands.New(
			"relocate",
			"Relocate a server",
			"upctl server relocate 00038afc-d526-4148-af0e-d2f1eeaded9b --zone fi-priv-example",
			"upctl server relocate 00038afc-d526-4148-af0e-d2f1eeaded9b 0053a6f5-e6d1-4b0b-b9dc-b90d0894e8d0 --zone fi-priv-example",
			"upctl server relocate my_server1 --zone fi-priv-example",
		),
	}
}

type relocateCommand struct {
	*commands.BaseCommand
	completion.StoppedServer
	resolver.CachingServer
	zone string
}

// InitCommand implements Command.InitCommand
func (s *relocateCommand) InitCommand() {
	s.Cobra().Long = commands.WrapLongDescription(`Relocate a server

	Relocates server with its storages and their backups to another zone. This feature can be used to move server from public zone to private cloud zone or vice versa. It can also be used to move server from one private cloud zone to another.

	For the relocation to succeed, both source and destination zones need to reside in the same physical location (ie. datacenter). The server cannot be attached to a SDN private network and it cannot have IP addresses from dedicated, customer-owned IP networks. Server needs to be in stopped state while its storages and their backups need to be in online state.`)

	fs := &pflag.FlagSet{}
	fs.StringVar(&s.zone, "zone", "", "The zone where the server should be relocated to.")
	s.AddFlags(fs)
	commands.Must(s.Cobra().MarkFlagRequired("zone"))
}

func (s *relocateCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("zone", namedargs.CompletionFunc(completion.Zone{}, cfg)))
}

// Execute implements commands.MultipleArgumentCommand
func (s *relocateCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.Server()
	msg := fmt.Sprintf("Relocating server %v to zone %v", uuid, s.zone)
	exec.PushProgressStarted(msg)

	res, err := svc.RelocateServer(exec.Context(), &request.RelocateServerRequest{
		UUID: uuid,
		Zone: s.zone,
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}
