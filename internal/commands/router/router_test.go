package router

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
)

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

	mService := mock.Service{}

	getRouters := "GetRouters"
	mService.On(getRouters).Return(&upcloud.Routers{Routers: []upcloud.Router{Router1, Router2}}, nil)

	buildRequestFn := func(uuid string) interface{} {
		return mockRequest
	}

	for _, test := range []struct {
		name    string
		args    []string
		request routerRequest
		calls   int
		error   string
	}{
		{
			name: "no router",
			args: []string{},
			request: routerRequest{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 0,
			error: "at least one router uuid is required",
		},
		{
			name: "single router with UUID",
			args: []string{Router2.UUID},
			request: routerRequest{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 0,
		},
		{
			name: "single router with Name",
			args: []string{Router2.Name},
			request: routerRequest{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "single router, exactly once",
			args: []string{Router2.UUID},
			request: routerRequest{
				ExactlyOne:   true,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 0,
		},
		{
			name: "multiple router",
			args: []string{Router1.Name, Router2.Name},
			request: routerRequest{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "multiple router, exactly once",
			args: []string{Router1.UUID, Router2.UUID},
			request: routerRequest{
				ExactlyOne:   true,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			error: "single router uuid is required",
			calls: 0,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			cachedRouters = nil

			res, err := test.request.send(test.args)

			if test.error != "" && err != nil {
				assert.Equal(t, test.error, err.Error())
				assert.Nil(t, res)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, mockResponse, res)
			}

			mService.AssertNumberOfCalls(t, getRouters, test.calls)
			mService.Calls = nil
		})
	}
}
