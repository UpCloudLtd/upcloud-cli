package storage

import (
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var storages = &upcloud.Storages{Storages: []upcloud.Storage{
	{
		UUID:   mocks.Uuid1,
		Title:  mocks.Title1,
		Access: "private",
	},
	{
		UUID:   mocks.Uuid2,
		Title:  mocks.Title2,
		Access: "private",
	},
	{
		UUID:   mocks.Uuid3,
		Title:  mocks.Title3,
		Access: "public",
	},
},
}

func TestDeleteStorage(t *testing.T) {

	for _, testcase := range []struct {
		name   string
		args   []string
		testFn func(e error)
	}{
		{
			name:   "Storage with given title found and deleted successfully",
			args:   []string{mocks.Title1},
			testFn: func(e error) { assert.Nil(t, e) },
		},
		{
			name:   "Storage with given uuid found and deleted successfully",
			args:   []string{mocks.Uuid1},
			testFn: func(e error) { assert.Nil(t, e) },
		},
		{
			name: "Storage with given title does not exist",
			args: []string{"asdf"},
			testFn: func(e error) {
				assert.Equal(t, "no storage with uuid, name or title \"asdf\" was found", e.Error())
			},
		},
		{
			name: "No title or uuid given",
			args: []string{},
			testFn: func(e error) {
				assert.Equal(t, "server hostname, title or uuid is required", e.Error())
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			mss := new(mocks.MockStorageService)
			mss.On("GetStorages", mock.Anything).Return(storages, nil)
			mss.On("DeleteStorage", mock.Anything).Return(nil)
			dc := DeleteCommand(mss)

			_, err := dc.MakeExecuteCommand()(testcase.args)

			testcase.testFn(err)
		})
	}
}
