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

func (m MockStorageService) GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	return nil, nil
}
func (m MockStorageService) GetStorageDetails(r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error) {
	return nil, nil
}
func (m MockStorageService) CreateStorage(r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	return nil, nil
}
func (m MockStorageService) ModifyStorage(r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error) {
	return nil, nil
}
func (m MockStorageService) AttachStorage(r *request.AttachStorageRequest) (*upcloud.ServerDetails, error) {
	return nil, nil
}
func (m MockStorageService) DetachStorage(r *request.DetachStorageRequest) (*upcloud.ServerDetails, error) {
	return nil, nil
}
func (m MockStorageService) CloneStorage(r *request.CloneStorageRequest) (*upcloud.StorageDetails, error) {
	return nil, nil
}
func (m MockStorageService) TemplatizeStorage(r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error) {
	return nil, nil
}
func (m MockStorageService) WaitForStorageState(r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error) {
	return nil, nil
}
func (m MockStorageService) LoadCDROM(r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error) {
	return nil, nil
}
func (m MockStorageService) EjectCDROM(r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error) {
	return nil, nil
}
func (m MockStorageService) CreateBackup(r *request.CreateBackupRequest) (*upcloud.StorageDetails, error) {
	return nil, nil
}
func (m MockStorageService) RestoreBackup(r *request.RestoreBackupRequest) error { return nil }
func (m MockStorageService) CreateStorageImport(r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
	return nil, nil
}
func (m MockStorageService) GetStorageImportDetails(r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error) {
	return nil, nil
}
func (m MockStorageService) WaitForStorageImportCompletion(r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error) {
	return nil, nil
}
func (m MockStorageService) DeleteStorage(r *request.DeleteStorageRequest) error { return nil }
