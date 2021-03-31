package service

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

// Wrapper is a temporary reimplementation fo upcloud-go-api which prevents
// casting service.Service struct for the different interfaces used by the commands
// XXX: this should be removed once the upcloud-go-api supports services in the correct way
type Wrapper struct {
	*service.Service
}

func (s Wrapper) Server() service.Server {
	return s.Service
}

func (s Wrapper) Storage() service.Storage {
	return s.Service
}

func (s Wrapper) Network() service.Network {
	return s.Service
}

func (s Wrapper) Firewall() service.Firewall {
	return s.Service
}

func (s Wrapper) IpAddress() service.IpAddress {
	return s.Service
}

func (s Wrapper) Account() service.Account {
	return s.Service
}

func (s Wrapper) Plan() service.Plans {
	return s.Service
}
