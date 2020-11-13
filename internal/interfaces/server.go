package interfaces

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type Server interface {
	GetServerConfigurations() (*upcloud.ServerConfigurations, error)
	GetServers() (*upcloud.Servers, error)
	GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error)
	CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error)
	WaitForServerState(r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error)
	StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error)
	StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error)
	RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error)
	ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error)
	DeleteServer(r *request.DeleteServerRequest) error
	DeleteServerAndStorages(r *request.DeleteServerAndStoragesRequest) error
}
