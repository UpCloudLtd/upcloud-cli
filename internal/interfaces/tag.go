package interfaces

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type Tag interface {
	GetTags() (*upcloud.Tags, error)
	CreateTag(r *request.CreateTagRequest) (*upcloud.Tag, error)
	ModifyTag(r *request.ModifyTagRequest) (*upcloud.Tag, error)
	DeleteTag(r *request.DeleteTagRequest) error
	TagServer(r *request.TagServerRequest) (*upcloud.ServerDetails, error)
	UntagServer(r *request.UntagServerRequest) (*upcloud.ServerDetails, error)
}
