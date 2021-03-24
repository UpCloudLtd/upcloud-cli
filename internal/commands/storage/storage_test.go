package storage

import (
	"fmt"
	"testing"

	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	targetMethod = "GetStorages"
	Title1       = "mock-storage-title1"
	Title2       = "mock-storage-title2"
	Title3       = "mock-storage-title3"
	UUID1        = "0127dfd6-3884-4079-a948-3a8881df1a7a"
	UUID2        = "012bde1d-f0e7-4bb2-9f4a-74e1f2b49c07"
	UUID3        = "012c61a6-b8f0-48c2-a63a-b4bf7d26a655"
	UUID4        = "012c61a6-er4g-mf2t-b63a-b4be4326a655"
)

func TestSearchStorage(t *testing.T) {
	var Storage1 = upcloud.Storage{
		UUID:   UUID1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	var Storage2 = upcloud.Storage{
		UUID:   UUID2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	var Storage3 = upcloud.Storage{
		UUID:   UUID3,
		Title:  Title3,
		Access: "public",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   10,
		Tier:   "maxiops",
	}

	var Storage4 = upcloud.Storage{
		UUID:   UUID4,
		Title:  Title1,
		Access: "public",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   20,
		Tier:   "maxiops",
	}

	var storages = &upcloud.Storages{
		Storages: []upcloud.Storage{
			Storage1,
			Storage2,
			Storage3,
			Storage4,
		},
	}
	for _, testcase := range []struct {
		name         string
		args         []string
		expected     []string
		unique       bool
		additional   []upcloud.Storage
		backendCalls int
		errMsg       string
	}{
		{
			name:         "SingleUuid",
			args:         []string{Storage2.UUID},
			expected:     []string{Storage2.UUID},
			backendCalls: 0,
		},
		{
			name:         "MultipleUuidSearched",
			args:         []string{Storage2.UUID, Storage3.UUID},
			expected:     []string{Storage2.UUID, Storage3.UUID},
			backendCalls: 0,
		},
		{
			name:         "SingleTitle",
			args:         []string{Storage2.Title},
			expected:     []string{Storage2.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleTitlesSearched",
			args:         []string{Storage2.Title, Storage3.Title},
			expected:     []string{Storage2.UUID, Storage3.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleTitlesFound",
			args:         []string{Title1},
			expected:     []string{Storage1.UUID, Storage4.UUID},
			backendCalls: 1,
		},
		{
			name:         "MultipleTitlesFound_UniqueWanted",
			args:         []string{Title1},
			expected:     []string{Storage1.UUID, Storage4.UUID},
			backendCalls: 1,
			unique:       true,
			errMsg:       "multiple storages matched to query \"" + Title1 + "\", use UUID to specify",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			CachedStorages = nil

			mService := smock.MockService{}
			mService.On(targetMethod, mock.Anything).Return(storages, nil)

			result, err := searchAllStorages(testcase.args, &mService, testcase.unique)

			if testcase.errMsg == "" {
				assert.Nil(t, err)
				assert.ElementsMatch(t, testcase.expected, result)
			} else {
				assert.Nil(t, result)
				assert.EqualError(t, err, testcase.errMsg)
			}
			mService.AssertNumberOfCalls(t, targetMethod, testcase.backendCalls)
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

func TestSendStorageRequest(t *testing.T) {
	mService := smock.MockService{}

	var Storage1 = upcloud.Storage{
		UUID:   UUID1,
		Title:  Title1,
		Access: "private",
		State:  "maintenance",
		Type:   "backup",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	var Storage2 = upcloud.Storage{
		UUID:   UUID2,
		Title:  Title2,
		Access: "private",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   40,
		Tier:   "maxiops",
	}

	var Storage3 = upcloud.Storage{
		UUID:   UUID3,
		Title:  Title3,
		Access: "public",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   10,
		Tier:   "maxiops",
	}

	var Storage4 = upcloud.Storage{
		UUID:   UUID4,
		Title:  Title1,
		Access: "public",
		State:  "online",
		Type:   "normal",
		Zone:   "fi-hel1",
		Size:   20,
		Tier:   "maxiops",
	}

	var storages = &upcloud.Storages{
		Storages: []upcloud.Storage{
			Storage1,
			Storage2,
			Storage3,
			Storage4,
		},
	}

	mService.On(targetMethod, mock.Anything).Return(storages, nil)

	buildRequestFn := func(storage string) (interface{}, error) {
		return mockRequest, nil
	}

	for _, test := range []struct {
		name    string
		args    []string
		request storageRequest
		calls   int
		error   string
	}{
		{
			name: "no storage",
			args: []string{},
			request: storageRequest{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 0,
			error: "at least one storage uuid is required",
		},
		{
			name: "single storage",
			args: []string{Storage2.Title},
			request: storageRequest{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "single storage, exactly once",
			args: []string{Storage2.Title},
			request: storageRequest{
				ExactlyOne:   true,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "multiple storages",
			args: []string{Storage3.Title, Storage2.UUID},
			request: storageRequest{
				ExactlyOne:   false,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			calls: 1,
		},
		{
			name: "multiple storages, exactly once",
			args: []string{Storage1.UUID, Storage2.UUID},
			request: storageRequest{
				ExactlyOne:   true,
				BuildRequest: buildRequestFn,
				Service:      &mService,
				Handler:      MockHandler{},
			},
			error: "single storage uuid is required",
			calls: 0,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			CachedStorages = nil

			res, err := test.request.send(test.args)

			if test.error != "" && err != nil {
				assert.Equal(t, test.error, err.Error())
				assert.Nil(t, res)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, mockResponse, res)
			}

			mService.AssertNumberOfCalls(t, targetMethod, test.calls)
			mService.Calls = nil
		})
	}
}
