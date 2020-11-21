package storage

import (
  "github.com/UpCloudLtd/cli/internal/mocks"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "testing"
)

var (
  Title1 = "mock-storage-title1"
  Title2 = "mock-storage-title2"
  Title3 = "mock-storage-title3"
  Uuid1  = "0127dfd6-3884-4079-a948-3a8881df1a7a"
  Uuid2  = "012bde1d-f0e7-4bb2-9f4a-74e1f2b49c07"
  Uuid3  = "012c61a6-b8f0-48c2-a63a-b4bf7d26a655"
  Uuid4  = "012c61a6-er4g-mf2t-b63a-b4be4326a655"
)

var Storage1 = upcloud.Storage{
  UUID:   Uuid1,
  Title:  Title1,
  Access: "private",
  State:  "maintenance",
  Type:   "backup",
  Zone:   "fi-hel1",
  Size:   40,
  Tier:   "maxiops",
}

var Storage2 = upcloud.Storage{
  UUID:   Uuid2,
  Title:  Title2,
  Access: "private",
  State:  "online",
  Type:   "normal",
  Zone:   "fi-hel1",
  Size:   40,
  Tier:   "maxiops",
}

var Storage3 = upcloud.Storage{
  UUID:   Uuid3,
  Title:  Title3,
  Access: "public",
  State:  "online",
  Type:   "normal",
  Zone:   "fi-hel1",
  Size:   10,
  Tier:   "maxiops",
}

var Storage4 = upcloud.Storage{
  UUID:   Uuid4,
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
  },
}

func TestSearchStorage(t *testing.T) {
  var storages = &upcloud.Storages{
    Storages: []upcloud.Storage{
      Storage1,
      Storage2,
      Storage3,
      Storage4,
    },
  }
  for _, testcase := range []struct {
    name string
    args []string
    expected []*upcloud.Storage
    unique bool
    additional []upcloud.Storage
    backendCalls int
    errMsg string
  } {
    {
      name: "SingleUuid",
      args: []string{Storage2.UUID},
      expected: []*upcloud.Storage{&Storage2},
      backendCalls: 1,
    },
    {
      name: "MultipleUuidSearched",
      args: []string{Storage2.UUID, Storage3.UUID},
      expected: []*upcloud.Storage{&Storage2, &Storage3},
      backendCalls: 1,
    },
    {
      name: "SingleTitle",
      args: []string{Storage2.Title},
      expected: []*upcloud.Storage{&Storage2},
      backendCalls: 1,
    },
    {
      name: "MultipleTitlesSearched",
      args: []string{Storage2.Title, Storage3.Title},
      expected: []*upcloud.Storage{&Storage2, &Storage3},
      backendCalls: 1,
    },
    {
      name: "MultipleTitlesFound",
      args: []string{Title1},
      expected: []*upcloud.Storage{&Storage1, &Storage4},
      backendCalls: 1,
    },
    {
      name: "MultipleTitlesFound",
      args: []string{Title1},
      expected: []*upcloud.Storage{&Storage1, &Storage4},
      backendCalls: 1,
      unique: true,
      errMsg: "multiple storages matched to query \"" + Title1 + "\", use UUID to specify",
    },
  } {
    t.Run(testcase.name, func(t *testing.T) {
      cachedStorages = []upcloud.Storage{}
      mss := new(mocks.MockStorageService)
      mss.On("GetStorages", mock.Anything).Return(storages, nil)

      result, err := searchAllArgs(testcase.args, mss, testcase.unique)

      if testcase.errMsg == "" {
        assert.Nil(t, err)
        assert.ElementsMatch(t, testcase.expected, result)
      } else {
        assert.Nil(t, result)
        assert.EqualError(t, err, testcase.errMsg)
      }
      mss.AssertNumberOfCalls(t, "GetStorages", testcase.backendCalls)
    })
  }
}


