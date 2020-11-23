package ip_address

import (
  "github.com/UpCloudLtd/cli/internal/commands"
  "github.com/UpCloudLtd/cli/internal/ui"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
  "github.com/spf13/pflag"
)

type modifyCommand struct {
  *commands.BaseCommand
  service service.IpAddress
  req request.ModifyIPAddressRequest
}

func ModifyCommand(service service.IpAddress) commands.Command {
  return &modifyCommand{
    BaseCommand: commands.New("modify", "Modify an ip address"),
    service:     service,
  }
}

func (s *modifyCommand) InitCommand() {
 fs := &pflag.FlagSet{}
 fs.StringVar(&s.req.MAC, "mac", "", "MAC address of server interface to attach floating IP to.")
 fs.StringVar(&s.req.PTRRecord, "ptr-record", "", "A fully qualified domain name.")
 s.AddFlags(fs)
}

func (s *modifyCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
  return func(args []string) (interface{}, error) {
    return Request{
      BuildRequest: func(ip *upcloud.IPAddress) interface{} {
        return &request.ModifyIPAddressRequest{IPAddress: ip.Address}
      },
      Service: s.service,
      HandleContext: ui.HandleContext{
        RequestId:     func(in interface{}) string { return in.(*request.ModifyIPAddressRequest).IPAddress },
        MaxActions:    maxIpAddressActions,
        InteractiveUi: s.Config().InteractiveUI(),
        ActionMsg:     "Modifying IP Address",
        Action: func(req interface{}) (interface{}, error) {
          return s.service.ModifyIPAddress(req.(*request.ModifyIPAddressRequest))
        },
      },
    }.Send(args)
  }
}
