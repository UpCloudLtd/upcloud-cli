package ip_address

import (
  "github.com/UpCloudLtd/cli/internal/commands"
  "github.com/UpCloudLtd/cli/internal/ui"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

type releaseCommand struct {
 *commands.BaseCommand
 service service.IpAddress
}

func ReleaseCommand(service service.IpAddress) commands.Command {
 return &releaseCommand{
   BaseCommand: commands.New("release", "Release an ip address"),
   service:     service,
 }
}

func (s *releaseCommand) MakeExecuteCommand() func(args []string) (interface{}, error) {
 return func(args []string) (interface{}, error) {
   return Request{
     BuildRequest: func(ip *upcloud.IPAddress) interface{} {
       return &request.ReleaseIPAddressRequest{IPAddress: ip.Address}
     },
     Service: s.service,
     HandleContext: ui.HandleContext{
       RequestId:     func(in interface{}) string { return in.(*request.ReleaseIPAddressRequest).IPAddress },
       MaxActions:    maxIpAddressActions,
       InteractiveUi: s.Config().InteractiveUI(),
       ActionMsg:     "Releasing IP Address",
       Action: func(req interface{}) (interface{}, error) {
         return nil, s.service.ReleaseIPAddress(req.(*request.ReleaseIPAddressRequest))
       },
     },
   }.Send(args)
 }
}
