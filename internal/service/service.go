package service

import (
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
)

// AllServices is a convenience interface for providing all services
type AllServices interface {
	service.Cloud
	service.Network
	service.NetworkPeering
	service.Storage
	service.Server
	service.Firewall
	service.IPAddress
	service.Account
	service.ManagedDatabaseServiceManager
	service.LoadBalancer
	service.Kubernetes
	service.ServerGroup
	service.ManagedObjectStorage
	service.Permission
	service.Gateway
	service.Host
	service.Partner
	service.Token
	service.Tag
	service.AuditLog
}
