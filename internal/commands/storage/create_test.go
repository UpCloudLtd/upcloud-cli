package storage

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/cli/internal/mocks"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type CreateTestMock struct {
	mocks.MockStorageService
}

func (s CreateTestMock) CreateStorage(r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	return &upcloud.StorageDetails{
		Storage:     upcloud.Storage{UUID: mocks.Uuid1},
		BackupRule:  nil,
		BackupUUIDs: upcloud.BackupUUIDSlice{mocks.Uuid2},
		ServerUUIDs: upcloud.ServerUUIDSlice{mocks.Uuid3},
	}, nil
}

func TestCreateStorage(t *testing.T) {

	for _, testcase := range []struct {
		name   string
		args   []string
		testFn func(res *upcloud.StorageDetails, e error)
	}{
		{
			name: "Storage with given title found",
			args: []string{"--title", mocks.Title1, "--size", "1234", "--tier", "test-tier", "--zone", "fi-hel1"},
			testFn: func(res *upcloud.StorageDetails, e error) {
				assert.Equal(t, res.UUID, mocks.Uuid1)
				assert.Nil(t, e)
			},
		},
		{
			name: "Storage with given title does not exist",
			args: []string{"--asdf", "something"},
			testFn: func(res *upcloud.StorageDetails, e error) {
				assert.Nil(t, res)
				assert.Equal(t, "unknown flag: --asdf", e.Error())
			},
		},
		{
			name: "When no argument given default parameters are used",
			args: []string{},
			testFn: func(res *upcloud.StorageDetails, e error) {
				assert.Equal(t, res.UUID, mocks.Uuid1)
				assert.Nil(t, e)
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			cc := commands.BuildCommand(CreateCommand(CreateTestMock{}), nil, config.New(viper.New()))

			res, err := cc.MakeExecuteCommand()(testcase.args)
			var result []*upcloud.StorageDetails
			if res != nil {
				result = res.([]*upcloud.StorageDetails)
				testcase.testFn(result[0], err)
			} else {
				testcase.testFn(nil, err)
			}
		})
	}
}
