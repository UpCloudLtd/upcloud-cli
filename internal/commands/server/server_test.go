package server

import (
	"fmt"
	"testing"

	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	Title1 = "mock-storage-title1"
	Title2 = "mock-storage-title2"
	UUID1  = "0127dfd6-3884-4079-a948-3a8881df1a7a"
	UUID2  = "012bde1d-f0e7-4bb2-9f4a-74e1f2b49c07"
	UUID3  = "012c61a6-b8f0-48c2-a63a-b4bf7d26a655"
)

func TestSearchServer(t *testing.T) {

	var Server1 = upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server-1-hostname",
		License:      0,
		MemoryAmount: 1024,
		Plan:         "server-1-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-1-title",
		UUID:         "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		Zone:         "fi-hel1",
	}

	var Server2 = upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server-2-hostname",
		License:      0,
		MemoryAmount: 1024,
		Plan:         "server-2-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-2-title",
		UUID:         "f77a5b25-84af-4f52-bc40-581930091fad",
		Zone:         "fi-hel1",
	}

	var Server3 = upcloud.Server{
		CoreNumber:   2,
		Hostname:     "server-3-hostname",
		License:      0,
		MemoryAmount: 4096,
		Plan:         "server-3-plan",
		Progress:     0,
		State:        "stopped",
		Tags:         nil,
		Title:        "server-3-title",
		UUID:         "f0131b8f-ffe0-4271-83a8-c75b99e168c3",
		Zone:         "hu-bud1",
	}

	var Server4 = upcloud.Server{
		CoreNumber:   4,
		Hostname:     "server-4-hostname",
		License:      0,
		MemoryAmount: 5120,
		Plan:         "server-4-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        Server1.Title,
		UUID:         "e5b3a855-cd8a-45b6-8cef-c7c860a02217",
		Zone:         "uk-lon1",
	}

	var Server5 = upcloud.Server{
		CoreNumber:   4,
		Hostname:     Server4.Hostname,
		License:      0,
		MemoryAmount: 5120,
		Plan:         "server-5-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-5-title",
		UUID:         "39bc2725-213d-46c8-8b25-49990c6966a7",
		Zone:         "uk-lon1",
	}

	var servers = &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
			Server2,
			Server3,
			Server4,
			Server5,
		},
	}

	for _, testcase := range []struct {
		name         string
		args         []string
		expected     []string
		unique       bool
		additional   []upcloud.Server
		backendCalls int
		errMsg       string
	}{
		{
			name:         "SingleUuid",
			args:         []string{Server2.UUID},
			expected:     []string{Server2.UUID},
			backendCalls: 0,
		},
		{
			name:         "MultipleUuidSearched",
			args:         []string{Server2.UUID, Server3.UUID},
			expected:     []string{Server2.UUID, Server3.UUID},
			backendCalls: 0,
		},
		{
			name:         "SingleTitle",
			args:         []string{Server2.Title},
			expected:     []string{Server2.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleTitlesSearched",
			args:         []string{Server2.Title, Server3.Title},
			expected:     []string{Server2.UUID, Server3.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleTitlesFound",
			args:         []string{Server1.Title},
			expected:     []string{Server1.UUID, Server4.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleTitlesNotAllowed",
			args:         []string{Server1.Title},
			expected:     []string{Server1.UUID, Server4.UUID},
			backendCalls: 1,
			unique:       true,
			errMsg:       "multiple servers matched to query \"" + Server1.Title + "\", use UUID to specify",
		},
		{
			name:         "SingleHostname",
			args:         []string{Server2.Hostname},
			expected:     []string{Server2.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleHostnamesSearched",
			args:         []string{Server2.Hostname, Server3.Hostname},
			expected:     []string{Server2.UUID, Server3.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleHostnamesFound",
			args:         []string{Server4.Hostname},
			expected:     []string{Server4.UUID, Server5.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleHostnamesNotAllowed",
			args:         []string{Server4.Hostname},
			expected:     []string{Server4.UUID, Server5.UUID},
			backendCalls: 1,
			unique:       true,
			errMsg:       "multiple servers matched to query \"" + Server4.Hostname + "\", use UUID to specify",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			CachedServers = []upcloud.Server{}
			mService := smock.Service{}
			mService.On("GetServers", mock.Anything).Return(servers, nil)

			result, err := SearchAllServers(testcase.args, &mService, testcase.unique)

			if testcase.errMsg == "" {
				assert.Nil(t, err)
				assert.ElementsMatch(t, testcase.expected, result)
			} else {
				assert.Nil(t, result)
				assert.EqualError(t, err, testcase.errMsg)
			}
			mService.AssertNumberOfCalls(t, "GetServers", testcase.backendCalls)
		})
	}
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

func TestSendServerRequest(t *testing.T) {
	var Server1 = upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server-1-hostname",
		License:      0,
		MemoryAmount: 1024,
		Plan:         "server-1-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-1-title",
		UUID:         "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		Zone:         "fi-hel1",
	}

	var Server2 = upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server-2-hostname",
		License:      0,
		MemoryAmount: 1024,
		Plan:         "server-2-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-2-title",
		UUID:         "f77a5b25-84af-4f52-bc40-581930091fad",
		Zone:         "fi-hel1",
	}

	var Server3 = upcloud.Server{
		CoreNumber:   2,
		Hostname:     "server-3-hostname",
		License:      0,
		MemoryAmount: 4096,
		Plan:         "server-3-plan",
		Progress:     0,
		State:        "stopped",
		Tags:         nil,
		Title:        "server-3-title",
		UUID:         "f0131b8f-ffe0-4271-83a8-c75b99e168c3",
		Zone:         "hu-bud1",
	}

	var Server4 = upcloud.Server{
		CoreNumber:   4,
		Hostname:     "server-4-hostname",
		License:      0,
		MemoryAmount: 5120,
		Plan:         "server-4-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        Server1.Title,
		UUID:         "e5b3a855-cd8a-45b6-8cef-c7c860a02217",
		Zone:         "uk-lon1",
	}

	var Server5 = upcloud.Server{
		CoreNumber:   4,
		Hostname:     Server4.Hostname,
		License:      0,
		MemoryAmount: 5120,
		Plan:         "server-5-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-5-title",
		UUID:         "39bc2725-213d-46c8-8b25-49990c6966a7",
		Zone:         "uk-lon1",
	}

	var servers = &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
			Server2,
			Server3,
			Server4,
			Server5,
		},
	}

	mService := smock.Service{}

	getServers := "GetServers"
	mService.On(getServers, mock.Anything).Return(servers, nil)

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
			name: "no server",
			args: []string{},
			request: Request{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 0,
			error: "at least one server uuid is required",
		},
		{
			name: "single server",
			args: []string{Server2.Title},
			request: Request{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "single server, exactly once",
			args: []string{Server2.Hostname},
			request: Request{
				ExactlyOne:   true,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "multiple servers",
			args: []string{Server1.Hostname, Server2.Title},
			request: Request{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "multiple servers, exactly once",
			args: []string{Server1.UUID, Server2.UUID},
			request: Request{
				ExactlyOne:   true,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			error: "single server uuid is required",
			calls: 0,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedServers = nil

			res, err := test.request.Send(test.args)

			if test.error != "" && err != nil {
				assert.Equal(t, test.error, err.Error())
				assert.Nil(t, res)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, mockResponse, res)
			}

			mService.AssertNumberOfCalls(t, getServers, test.calls)
			mService.Calls = nil
		})
	}
}
