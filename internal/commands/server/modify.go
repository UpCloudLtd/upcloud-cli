package server

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/ui"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"github.com/spf13/pflag"
)

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

var DefaultModifyParams = modifyParams{
	ModifyServerRequest: request.ModifyServerRequest{},
}

func (s *modifyCommand) InitCommand() {
	s.params = modifyParams{ModifyServerRequest: request.ModifyServerRequest{}}
	flags := &pflag.FlagSet{}
	flags.StringVar(&s.params.BootOrder, "boot-order", DefaultModifyParams.BootOrder, "The boot device order.")
	flags.IntVar(&s.params.CoreNumber, "cores", DefaultModifyParams.CoreNumber, "Number of cores")
	flags.StringVar(&s.params.Hostname, "hostname", DefaultModifyParams.Hostname, "Hostname")
	flags.StringVar(&s.params.Firewall, "firewall", DefaultModifyParams.Firewall, "Sets the firewall on or off. You can manage firewall rules with the firewall command\nAvailable: on, off")
	flags.IntVar(&s.params.MemoryAmount, "memory", DefaultModifyParams.MemoryAmount, "Memory amount in MiB")
	flags.StringVar(&s.params.metadata, "metadata", DefaultModifyParams.metadata, "Enable metadata service")
	flags.StringVar(&s.params.Plan, "plan", DefaultModifyParams.Plan, "Server plan to use. Set this to custom to use custom core/memory amounts.")
	flags.StringVar(&s.params.SimpleBackup, "simple-backup", DefaultModifyParams.SimpleBackup, "Simple backup rule. Format (HHMM,{dailies,weeklies,monthlies}).\nExample: 2300,dailies")
	flags.StringVar(&s.params.Title, "title", DefaultModifyParams.Title, "Visible name")
	flags.StringVar(&s.params.TimeZone, "time-zone", DefaultModifyParams.TimeZone, "Time zone to set the RTC to")
	flags.StringVar(&s.params.VideoModel, "video-model", DefaultModifyParams.VideoModel, "Video interface model of the server.\nAvailable: vga,cirrus")
	flags.StringVar(&s.params.remoteAccessEnabled, "remote-access-enabled", DefaultModifyParams.remoteAccessEnabled, "Enables or disables the remote access\nAvailable: true, false")
	flags.StringVar(&s.params.RemoteAccessType, "remote-access-type", DefaultModifyParams.RemoteAccessType, "The remote access type")
	flags.StringVar(&s.params.RemoteAccessPassword, "remote-access-password", DefaultModifyParams.RemoteAccessPassword, "The remote access password")

	s.AddFlags(flags)
}

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

		return Request{
			BuildRequest: func(server *upcloud.Server) interface{} {
				req := s.params.ModifyServerRequest
				req.UUID = server.UUID
				return &req
			},
			Service:    s.service,
			HandleContext: ui.HandleContext{
				RequestId:     func(in interface{}) string { return in.(*request.ModifyServerRequest).UUID },
				ResultUuid:    getServerDetailsUuid,
				InteractiveUi: s.Config().InteractiveUI(),
				MaxActions:    maxServerActions,
				ActionMsg:     "Modifying server",
				Action: func(req interface{}) (interface{}, error) {
					return s.service.ModifyServer(req.(*request.ModifyServerRequest))
				},
			},
		}.Send(args)
	}
}
