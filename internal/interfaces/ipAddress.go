package interfaces

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type IpAddress interface {
	GetIPAddresses() (*upcloud.IPAddresses, error)
	GetIPAddressDetails(r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error)
	AssignIPAddress(r *request.AssignIPAddressRequest) (*upcloud.IPAddress, error)
	ModifyIPAddress(r *request.ModifyIPAddressRequest) (*upcloud.IPAddress, error)
	ReleaseIPAddress(r *request.ReleaseIPAddressRequest) error
}
