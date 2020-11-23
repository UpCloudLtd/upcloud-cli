package mocks

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/mock"
)

type MockTagService struct {
	mock.Mock
}

func (m *MockTagService) GetTags() (*upcloud.Tags, error) {
	args := m.Called()
	return args[0].(*upcloud.Tags), args.Error(1)
}
func (m *MockTagService) CreateTag(r *request.CreateTagRequest) (*upcloud.Tag, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Tag), args.Error(1)
}
func (m *MockTagService) ModifyTag(r *request.ModifyTagRequest) (*upcloud.Tag, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Tag), args.Error(1)
}
func (m *MockTagService) DeleteTag(r *request.DeleteTagRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
func (m *MockTagService) TagServer(r *request.TagServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockTagService) UntagServer(r *request.UntagServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
