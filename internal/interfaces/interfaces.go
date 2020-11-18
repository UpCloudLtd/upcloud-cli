package interfaces

import "github.com/UpCloudLtd/upcloud-go-api/upcloud/service"

type ServerAndStorage interface {
	service.Server
	service.Storage
}
