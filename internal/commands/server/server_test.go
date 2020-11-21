package server

import (
  "github.com/UpCloudLtd/cli/internal/mocks"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "testing"
)

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

func MockServerService() *mocks.MockServerService {
  mss := new(mocks.MockServerService)
  mss.On("GetServers", mock.Anything).Return(servers, nil)
  return mss
}

func TestSearchServer(t *testing.T) {
  for _, testcase := range []struct {
    name string
    args []string
    expected []*upcloud.Server
    unique bool
    additional []upcloud.Server
    backendCalls int
    errMsg string
  } {
    {
     name: "SingleUuid",
     args: []string{Server2.UUID},
     expected: []*upcloud.Server{&Server2},
     backendCalls: 1,
    },
    {
     name: "MultipleUuidSearched",
     args: []string{Server2.UUID, Server3.UUID},
     expected: []*upcloud.Server{&Server2, &Server3},
     backendCalls: 1,
    },
    {
     name: "SingleTitle",
     args: []string{Server2.Title},
     expected: []*upcloud.Server{&Server2},
     backendCalls: 1,
    },
    {
     name: "MultipleTitlesSearched",
     args: []string{Server2.Title, Server3.Title},
     expected: []*upcloud.Server{&Server2, &Server3},
     backendCalls: 1,
    },
    {
     name: "MultipleTitlesFound",
     args: []string{Server1.Title},
     expected: []*upcloud.Server{&Server1, &Server4},
     backendCalls: 1,
    },
    {
     name: "MultipleTitlesNotAllowed",
     args: []string{Server1.Title},
     expected: []*upcloud.Server{&Server1, &Server4},
     backendCalls: 1,
     unique: true,
     errMsg: "multiple servers matched to query \"" + Server1.Title + "\", use UUID to specify",
    },
    {
     name: "SingleHostname",
     args: []string{Server2.Hostname},
     expected: []*upcloud.Server{&Server2},
     backendCalls: 1,
    },
    {
     name: "MultipleHostnamesSearched",
     args: []string{Server2.Hostname, Server3.Hostname},
     expected: []*upcloud.Server{&Server2, &Server3},
     backendCalls: 1,
    },
    {
      name: "MultipleHostnamesFound",
      args: []string{Server4.Hostname},
      expected: []*upcloud.Server{&Server4, &Server5},
      backendCalls: 1,
    },
    {
      name: "MultipleHostnamesNotAllowed",
      args: []string{Server4.Hostname},
      expected: []*upcloud.Server{&Server4, &Server5},
      backendCalls: 1,
      unique: true,
      errMsg: "multiple servers matched to query \"" + Server4.Hostname + "\", use UUID to specify",
    },
  } {
    t.Run(testcase.name, func(t *testing.T) {
      cachedServers = []upcloud.Server{}
      mss := new(mocks.MockServerService)
      mss.On("GetServers", mock.Anything).Return(servers, nil)

      result, err := searchAllArgs(testcase.args, mss, testcase.unique)

      if testcase.errMsg == "" {
        assert.Nil(t, err)
        assert.ElementsMatch(t, testcase.expected, result)
      } else {
        assert.Nil(t, result)
        assert.EqualError(t, err, testcase.errMsg)
      }
      mss.AssertNumberOfCalls(t, "GetServers", testcase.backendCalls)
    })
  }
}
