package mocks

import (
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/spf13/viper"
)

var (
	Title1 = "mock-storage-title1"
	Title2 = "mock-storage-title2"
	Title3 = "mock-storage-title3"
	Uuid1  = "mock-uuid-1"
	Uuid2  = "mock-uuid-2"
	Uuid3  = "mock-uuid-3"
)

func GetBaseCommand() *commands.BaseCommand {
	v := viper.New()
	v.Set(config.ConfigKeyOutput, config.ConfigValueOutputJson)

	c := config.New(v)
	c.SetNamespace("testing")

	bc := commands.New("list", "dummy usage")
	bc.SetConfig(c)

	return bc
}

type MockStorageService struct{}

func (m MockStorageService) DeleteStorage(r *request.DeleteStorageRequest) error {
	return nil
}

func (m MockStorageService) GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	return nil, nil
}

func (m MockStorageService) CreateStorage(r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	return nil, nil
}
