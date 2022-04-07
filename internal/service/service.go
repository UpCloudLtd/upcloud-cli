package service

import (
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/service"
)

// AllServices is a convenience interface for providing all services
type AllServices interface {
	service.Network
	service.Storage
	service.Server
	service.Firewall
	service.IpAddress
	service.Plans
	service.Account
}

// Wrapper is a temporary reimplementation fo upcloud-go-api which prevents
// casting service.Service struct for the different interfaces used by the commands
// XXX: this should be removed once the upcloud-go-api supports services in the correct way
type Wrapper struct {
	Service AllServices
}

// Server returns sub service Server
func (s Wrapper) Server() service.Server {
	return s.Service
}

// Storage returns sub service Storage
func (s Wrapper) Storage() service.Storage {
	return s.Service
}

// Network returns sub service Network
func (s Wrapper) Network() service.Network {
	return s.Service
}

// Firewall returns sub service Firewall
func (s Wrapper) Firewall() service.Firewall {
	return s.Service
}

// IPAddress returns sub service IpAddress
func (s Wrapper) IPAddress() service.IpAddress {
	return s.Service
}

// Account returns sub service Account
func (s Wrapper) Account() service.Account {
	return s.Service
}

// Plan returns sub service Plans
func (s Wrapper) Plan() service.Plans {
	return s.Service
}
