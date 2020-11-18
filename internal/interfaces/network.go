package interfaces

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type Network interface {
	GetNetworks() (*upcloud.Networks, error)
	GetNetworksInZone(r *request.GetNetworksInZoneRequest) (*upcloud.Networks, error)
	CreateNetwork(r *request.CreateNetworkRequest) (*upcloud.Network, error)
	GetNetworkDetails(r *request.GetNetworkDetailsRequest) (*upcloud.Network, error)
	ModifyNetwork(r *request.ModifyNetworkRequest) (*upcloud.Network, error)
	DeleteNetwork(r *request.DeleteNetworkRequest) error
	GetServerNetworks(r *request.GetServerNetworksRequest) (*upcloud.Networking, error)
	CreateNetworkInterface(r *request.CreateNetworkInterfaceRequest) (*upcloud.Interface, error)
	ModifyNetworkInterface(r *request.ModifyNetworkInterfaceRequest) (*upcloud.Interface, error)
	DeleteNetworkInterface(r *request.DeleteNetworkInterfaceRequest) error
}
