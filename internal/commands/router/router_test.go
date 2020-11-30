package router

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRouterService struct {
	mock.Mock
}

func (m *MockRouterService) GetRouters() (*upcloud.Routers, error) {
	args := m.Called()
	return args[0].(*upcloud.Routers), args.Error(1)
}
func (m *MockRouterService) GetRouterDetails(r *request.GetRouterDetailsRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Router), args.Error(1)
}
func (m *MockRouterService) CreateRouter(r *request.CreateRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Router), args.Error(1)
}
func (m *MockRouterService) ModifyRouter(r *request.ModifyRouterRequest) (*upcloud.Router, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Router), args.Error(1)
}
func (m *MockRouterService) DeleteRouter(r *request.DeleteRouterRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

type MockNetworkService struct {
	mock.Mock
}

func (m *MockNetworkService) GetNetworks() (*upcloud.Networks, error) {
	args := m.Called()
	return args[0].(*upcloud.Networks), args.Error(1)
}
func (m *MockNetworkService) GetNetworksInZone(r *request.GetNetworksInZoneRequest) (*upcloud.Networks, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Networks), args.Error(1)
}
func (m *MockNetworkService) CreateNetwork(r *request.CreateNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Network), args.Error(1)
}
func (m *MockNetworkService) GetNetworkDetails(r *request.GetNetworkDetailsRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Network), args.Error(1)
}
func (m *MockNetworkService) ModifyNetwork(r *request.ModifyNetworkRequest) (*upcloud.Network, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Network), args.Error(1)
}
func (m *MockNetworkService) GetServerNetworks(r *request.GetServerNetworksRequest) (*upcloud.Networking, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Networking), args.Error(1)
}
func (m *MockNetworkService) CreateNetworkInterface(r *request.CreateNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Interface), args.Error(1)
}
func (m *MockNetworkService) ModifyNetworkInterface(r *request.ModifyNetworkInterfaceRequest) (*upcloud.Interface, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Interface), args.Error(1)
}
func (m *MockNetworkService) DeleteNetwork(r *request.DeleteNetworkRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
func (m *MockNetworkService) DeleteNetworkInterface(r *request.DeleteNetworkInterfaceRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

const mockResponse = "mock-response"
const mockRequest = "mock-request"

type MockHandler struct{}

func (s MockHandler) Handle(requests []interface{}) (interface{}, error) {
	for _, r := range requests {
		if r != mockRequest {
			return nil, fmt.Errorf("upexpected request %q", r)
		}
	}
	return mockResponse, nil
}

func TestSearchRouter(t *testing.T) {

	Router1 := upcloud.Router{
		Name: "name-1",
		Type: "normal",
		UUID: "ffd3ab80-fe95-49c0-ab70-fbc987246c7a",
	}

	Router2 := upcloud.Router{
		Name: "name-2",
		Type: "normal",
		UUID: "f14dd3e7-3dbb-4e3c-92b9-d1cf5178a13e",
	}

	mss := MockRouterService{}

	getRouters := "GetRouters"
	mss.On(getRouters).Return(&upcloud.Routers{Routers: []upcloud.Router{Router1, Router2}}, nil)

	buildRequestFn := func(uuid string) interface{} {
		return mockRequest
	}

	for _, test := range []struct {
		name    string
		args    []string
		request Request
		calls   int
		error   string
	}{
		{
			name: "no router",
			args: []string{},
			request: Request{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mss,
				Handler:      MockHandler{},
			},
			calls: 0,
			error: "at least one router uuid is required",
		},
		{
			name: "single router with UUID",
			args: []string{Router2.UUID},
			request: Request{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mss,
				Handler:      MockHandler{},
			},
			calls: 0,
		},
		{
			name: "single router with Name",
			args: []string{Router2.Name},
			request: Request{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mss,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "single router, exactly once",
			args: []string{Router2.UUID},
			request: Request{
				ExactlyOne:   true,
				BuildRequest: buildRequestFn,
				Service:      &mss,
				Handler:      MockHandler{},
			},
			calls: 0,
		},
		{
			name: "multiple router",
			args: []string{Router1.Name, Router2.Name},
			request: Request{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mss,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "multiple router, exactly once",
			args: []string{Router1.UUID, Router2.UUID},
			request: Request{
				ExactlyOne:   true,
				BuildRequest: buildRequestFn,
				Service:      &mss,
				Handler:      MockHandler{},
			},
			error: "single router uuid is required",
			calls: 0,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cachedRouters = nil

			res, err := test.request.Send(test.args)

			if test.error != "" && err != nil {
				assert.Equal(t, test.error, err.Error())
				assert.Nil(t, res)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, mockResponse, res)
			}

			mss.AssertNumberOfCalls(t, getRouters, test.calls)
			mss.Calls = nil
		})
	}
}
