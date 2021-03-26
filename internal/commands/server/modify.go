package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

// ModifyCommand creates the "server modify" command
func ModifyCommand(service service.Server) commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New("modify", "Modifies the configuration of an existing server"),
		service:     service,
	}
}

type modifyCommand struct {
	*commands.BaseCommand
	service service.Server
	params  modifyParams
}

type modifyParams struct {
	request.ModifyServerRequest
	remoteAccessEnabled string
	metadata            string
}

var defaultModifyParams = modifyParams{
	ModifyServerRequest: request.ModifyServerRequest{},
}

// InitCommand implements Command.InitCommand
func (s *modifyCommand) InitCommand() {
	s.SetPositionalArgHelp(PositionalArgHelp)
	s.ArgCompletion(GetServerArgumentCompletionFunction(s.service))
	s.params = modifyParams{ModifyServerRequest: request.ModifyServerRequest{}}
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.params.BootOrder, "boot-order", defaultModifyParams.BootOrder, "The boot device order.")
	flags.IntVar(&s.params.CoreNumber, "cores", defaultModifyParams.CoreNumber, "Number of cores.")
	flags.StringVar(&s.params.Hostname, "hostname", defaultModifyParams.Hostname, "Hostname.")
	flags.StringVar(&s.params.Firewall, "firewall", defaultModifyParams.Firewall, "Enables or disables firewall on the server. You can manage firewall rules with the firewall command.\nAvailable: true, false")
	flags.IntVar(&s.params.MemoryAmount, "memory", defaultModifyParams.MemoryAmount, "Memory amount in MiB.")
	flags.StringVar(&s.params.metadata, "metadata", defaultModifyParams.metadata, "Enable metadata service.")
	flags.StringVar(&s.params.Plan, "plan", defaultModifyParams.Plan, "Server plan to use. Set this to custom to use custom core/memory amounts.")
	flags.StringVar(&s.params.SimpleBackup, "simple-backup", defaultModifyParams.SimpleBackup, "Simple backup rule. Format (HHMM,{dailies,weeklies,monthlies}).\nExample: 2300,dailies")
	flags.StringVar(&s.params.Title, "title", defaultModifyParams.Title, "Visible name.")
	flags.StringVar(&s.params.TimeZone, "time-zone", defaultModifyParams.TimeZone, "Time zone to set the RTC to.")
	flags.StringVar(&s.params.VideoModel, "video-model", defaultModifyParams.VideoModel, "Video interface model of the server.\nAvailable: vga,cirrus")
	flags.StringVar(&s.params.remoteAccessEnabled, "remote-access-enabled", defaultModifyParams.remoteAccessEnabled, "Enables or disables the remote access.\nAvailable: true, false")
	flags.StringVar(&s.params.RemoteAccessType, "remote-access-type", defaultModifyParams.RemoteAccessType, "The remote access type.")
	flags.StringVar(&s.params.RemoteAccessPassword, "remote-access-password", defaultModifyParams.RemoteAccessPassword, "The remote access password.")

	s.AddFlags(flags)
}

// MakeExecuteCommand implements Command.MakeExecuteCommand
func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
	return func(args []string) (interface{}, error) {

		remoteAccess := new(upcloud.Boolean)
		if err := remoteAccess.UnmarshalJSON([]byte(s.params.remoteAccessEnabled)); err != nil {
			return nil, err
		}
		s.params.RemoteAccessEnabled = *remoteAccess

		metadata := new(upcloud.Boolean)
		if err := metadata.UnmarshalJSON([]byte(s.params.metadata)); err != nil {
			return nil, err
		}
		s.params.Metadata = *metadata

		switch s.params.Firewall {
		case "true":
			s.params.Firewall = "on"
		case "false":
			s.params.Firewall = "off"
		}

		return Request{
			BuildRequest: func(uuid string) interface{} {
				req := s.params.ModifyServerRequest
				req.UUID = uuid
				return &req
			},
			Service: s.service,
			Handler: ui.HandleContext{
				RequestID:     func(in interface{}) string { return in.(*request.ModifyServerRequest).UUID },
				ResultUUID:    getServerDetailsUUID,
				InteractiveUI: s.Config().InteractiveUI(),
				MaxActions:    maxServerActions,
				ActionMsg:     "Modifying server",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyServer(req.(*request.ModifyServerRequest))
				},
			},
		}.Send(args)
	}
}
