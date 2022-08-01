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
	service.Zones
	service.ManagedDatabaseServiceManager
	service.LoadBalancer
}
