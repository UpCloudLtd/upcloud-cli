package storage

import (
	"bytes"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type ListTestMock struct {
	mocks.MockStorageService
}

func (m ListTestMock) GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	var storages []upcloud.Storage
	storages = append(storages, mocks.Storage1, mocks.Storage2, mocks.Storage3)
	return &upcloud.Storages{Storages: storages}, nil
}

func TestListStorages(t *testing.T) {

	for _, testcase := range []struct {
		name    string
		private bool
		public  bool
		args    []string
		testFn  func(res upcloud.Storages, e error)
	}{
		{
			name: "List storages",
			args: []string{"--private", "--public"},
			testFn: func(res upcloud.Storages, e error) {
				assert.Equal(t, 3, len(res.Storages))
				assert.Nil(t, e)
			},
		},
		{
			name: "List private storages",
			args: []string{"--private"},
			testFn: func(res upcloud.Storages, e error) {
				assert.Equal(t, 2, len(res.Storages))
				assert.Nil(t, e)
			},
		},
		{
			name: "List public storages",
			args: []string{"--public"},
			testFn: func(res upcloud.Storages, e error) {
				assert.Equal(t, 3, len(res.Storages))
				assert.Nil(t, e)
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			stgs := upcloud.Storages{Storages: []upcloud.Storage{mocks.Storage1, mocks.Storage2, mocks.Storage3}}
			mss := new(mocks.MockStorageService)
			mss.On("GetStorages", mock.Anything).Return(&stgs, nil)

			lc := commands.BuildCommand(ListCommand(mss), nil, config.New(viper.New()))
			_ = mocks.SetFlags(lc, testcase.args...)
			res, err := lc.MakeExecuteCommand()([]string{})
			result := res.(*upcloud.Storages)
			testcase.testFn(*result, err)
		})
	}
}

func TestListStoragesOutput(t *testing.T) {
	storages := &upcloud.Storages{
		Storages: []upcloud.Storage{
			mocks.Storage1,
			mocks.Storage2,
			mocks.Storage3,
		},
	}

	lc := commands.BuildCommand(ListCommand(new(mocks.MockStorageService)), nil, config.New(viper.New()))

	expected := `
  UUID            Title                   Zone        State           Type       Size     Tier        Created  
─────────────── ─────────────────────── ─────────── ─────────────── ────────── ──────── ─────────── ───────────
  mock-uuid-1     mock-storage-title1     fi-hel1     maintenance     backup       40     maxiops              
  mock-uuid-2     mock-storage-title2     fi-hel1     online          normal       40     maxiops              
  mock-uuid-3     mock-storage-title3     fi-hel1     online          normal       10     maxiops              

`

	buf := new(bytes.Buffer)
	err := lc.HandleOutput(buf, storages)

	assert.Nil(t, err)
	assert.Equal(t, expected, buf.String())
}
