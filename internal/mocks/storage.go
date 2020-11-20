package mocks

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/mock"
)

var (
	Title1 = "mock-storage-title1"
	Title2 = "mock-storage-title2"
	Title3 = "mock-storage-title3"
	Uuid1  = "mock-uuid-1"
	Uuid2  = "mock-uuid-2"
	Uuid3  = "mock-uuid-3"
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

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
	args := m.Called(r)
	return args[0].(*upcloud.Storages), args.Error(1)
}
func (m *MockStorageService) GetStorageDetails(r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func (m *MockStorageService) CreateStorage(r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func (m *MockStorageService) ModifyStorage(r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func (m *MockStorageService) AttachStorage(r *request.AttachStorageRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockStorageService) DetachStorage(r *request.DetachStorageRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockStorageService) CloneStorage(r *request.CloneStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func (m *MockStorageService) TemplatizeStorage(r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func (m *MockStorageService) WaitForStorageState(r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func (m *MockStorageService) LoadCDROM(r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockStorageService) EjectCDROM(r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockStorageService) CreateBackup(r *request.CreateBackupRequest) (*upcloud.StorageDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageDetails), args.Error(1)
}

func (m *MockStorageService) RestoreBackup(r *request.RestoreBackupRequest) error {
	return m.Called(r).Error(0)
}
func (m *MockStorageService) CreateStorageImport(r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}
func (m *MockStorageService) GetStorageImportDetails(r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}
func (m *MockStorageService) WaitForStorageImportCompletion(r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}
func (m *MockStorageService) DeleteStorage(r *request.DeleteStorageRequest) error {
	return m.Called(r).Error(0)
}
