package interfaces

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type Storage interface {
	GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error)
	GetStorageDetails(r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error)
	CreateStorage(r *request.CreateStorageRequest) (*upcloud.StorageDetails, error)
	ModifyStorage(r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error)
	AttachStorage(r *request.AttachStorageRequest) (*upcloud.ServerDetails, error)
	DetachStorage(r *request.DetachStorageRequest) (*upcloud.ServerDetails, error)
	CloneStorage(r *request.CloneStorageRequest) (*upcloud.StorageDetails, error)
	TemplatizeStorage(r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error)
	WaitForStorageState(r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error)
	LoadCDROM(r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error)
	EjectCDROM(r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error)
	CreateBackup(r *request.CreateBackupRequest) (*upcloud.StorageDetails, error)
	RestoreBackup(r *request.RestoreBackupRequest) error
	CreateStorageImport(r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error)
	GetStorageImportDetails(r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error)
	WaitForStorageImportCompletion(r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error)
	DeleteStorage(*request.DeleteStorageRequest) error
}
