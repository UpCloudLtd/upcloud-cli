package storage

import (
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeleteStorage(t *testing.T) {

	for _, testcase := range []struct {
		name   string
		args   []string
		testFn func(e error)
	}{
		{
			name:   "Storage with given title found and deleted successfully",
			args:   []string{Title1},
			testFn: func(e error) { assert.Nil(t, e) },
		},
		{
			name:   "Storage with given uuid found and deleted successfully",
			args:   []string{Uuid1},
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
