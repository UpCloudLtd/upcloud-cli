package storage

import (
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ListTestMock struct {
	mocks.MockStorageService
}

var storage1 = upcloud.Storage{
	UUID:   mocks.Uuid1,
	Title:  mocks.Title1,
	Access: "private",
}

var storage2 = upcloud.Storage{
	UUID:   mocks.Uuid2,
	Title:  mocks.Title2,
	Access: "private",
}

var storage3 = upcloud.Storage{
	UUID:   mocks.Uuid3,
	Title:  mocks.Title3,
	Access: "public",
}

func (m ListTestMock) GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	var storages []upcloud.Storage
	storages = append(storages, storage1, storage2, storage3)
	return &upcloud.Storages{Storages: storages}, nil
}

func TestListStorages(t *testing.T) {

	for _, testcase := range []struct {
		name    string
		private bool
		public  bool
		testFn  func(res upcloud.Storages, e error)
	}{
		{
			name:    "List storages",
			private: true,
			public:  true,
			testFn: func(res upcloud.Storages, e error) {
				assert.Equal(t, 3, len(res.Storages))
				assert.Nil(t, e)
			},
		},
		{
			name:    "List private storages",
			private: true,
			public:  false,
			testFn: func(res upcloud.Storages, e error) {
				assert.Equal(t, 2, len(res.Storages))
				assert.Nil(t, e)
			},
		},
		{
			name:    "List public storages",
			private: false,
			public:  true,
			testFn: func(res upcloud.Storages, e error) {
				assert.Equal(t, 1, len(res.Storages))
				assert.Nil(t, e)
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {

			lc := listCommand{
				BaseCommand: mocks.GetBaseCommand(),
				service:     ListTestMock{},
				private:     testcase.private,
				public:      testcase.public,
			}

			res, err := lc.MakeExecuteCommand()([]string{})
			result := res.(*upcloud.Storages)
			testcase.testFn(*result, err)
		})
	}
}
