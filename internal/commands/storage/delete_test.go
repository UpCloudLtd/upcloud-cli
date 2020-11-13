package storage

import (
  "github.com/UpCloudLtd/cli/internal/commands"
  "github.com/UpCloudLtd/cli/internal/config"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
  "github.com/spf13/viper"
  "github.com/stretchr/testify/assert"
  "testing"
)

var title1 = "mock-storage-title1"

type mockDeleteService struct{}

func (m mockDeleteService) DeleteStorage(r *request.DeleteStorageRequest) error {
	return nil
}

func (m mockDeleteService) GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	var storages []upcloud.Storage
	storages = append(storages, upcloud.Storage{
		UUID: title1,
	})

	return &upcloud.Storages{Storages: storages}, nil
}

func TestDeleteStorage(t *testing.T) {

	for _, testcase := range []struct {
		name   string
		titles []string
		result string
		err    string
		testFn func(e error)
	}{
		{
			name:   "Storage with given title found",
			titles: []string{"mock-storage-title1"},
			testFn: func(e error) { assert.Nil(t, e) },
		},
		{
			name:   "Storage with given title does not exist",
			titles: []string{"mock-storage-title2"},
			testFn: func(e error) {
				assert.Equal(t, "no storage with uuid, name or title \"mock-storage-title2\" was found", e.Error())
			},
		},
		{
			name:   "No title given",
			titles: []string{},
			testFn: func(e error) {
				assert.Equal(t, "server hostname, title or uuid is required", e.Error())
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			v := viper.New()
			v.Set(config.ConfigKeyOutput, config.ConfigValueOutputHuman)

			bc := commands.New("delete", "Delete a storage")
			bc.SetConfig(config.New(v))

			dc := deleteCommand{
				BaseCommand: bc,
				service:     mockDeleteService{},
			}

			res := dc.MakeExecuteCommand()(testcase.titles)

			testcase.testFn(res)
		})
	}
}
