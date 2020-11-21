package mocks

import (
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
  "github.com/stretchr/testify/mock"
)

type MockServerStorageService struct {
  mock.Mock
}

func(m *MockServerStorageService) GetServerConfigurations() (*upcloud.ServerConfigurations, error) {
  args := m.Called()
  return args[0].(*upcloud.ServerConfigurations), args.Error(1)
}
func(m *MockServerStorageService) GetServers() (*upcloud.Servers, error) {
  args := m.Called()
  return args[0].(*upcloud.Servers), args.Error(1)
}
func(m *MockServerStorageService) GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService) CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService) WaitForServerState(r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService) StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService) StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService) RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService) ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService) DeleteServer(r *request.DeleteServerRequest) error {
  args := m.Called(r)
  return args.Error(0)
}
func(m *MockServerStorageService) DeleteServerAndStorages(r *request.DeleteServerAndStoragesRequest) error {
  args := m.Called(r)
  return args.Error(0)
}

func(m *MockServerStorageService)GetStorages(r *request.GetStoragesRequest) (*upcloud.Storages, error) {
  args := m.Called(r)
  return args[0].(*upcloud.Storages), args.Error(1)
}
func(m *MockServerStorageService)GetStorageDetails(r *request.GetStorageDetailsRequest) (*upcloud.StorageDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func(m *MockServerStorageService)CreateStorage(r *request.CreateStorageRequest) (*upcloud.StorageDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func(m *MockServerStorageService)ModifyStorage(r *request.ModifyStorageRequest) (*upcloud.StorageDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func(m *MockServerStorageService)AttachStorage(r *request.AttachStorageRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService)DetachStorage(r *request.DetachStorageRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService)CloneStorage(r *request.CloneStorageRequest) (*upcloud.StorageDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func(m *MockServerStorageService)TemplatizeStorage(r *request.TemplatizeStorageRequest) (*upcloud.StorageDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func(m *MockServerStorageService)WaitForStorageState(r *request.WaitForStorageStateRequest) (*upcloud.StorageDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func(m *MockServerStorageService)LoadCDROM(r *request.LoadCDROMRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService)EjectCDROM(r *request.EjectCDROMRequest) (*upcloud.ServerDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func(m *MockServerStorageService)CreateBackup(r *request.CreateBackupRequest) (*upcloud.StorageDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageDetails), args.Error(1)
}
func(m *MockServerStorageService)RestoreBackup(r *request.RestoreBackupRequest) error {
  return m.Called(r).Error(0)
}
func(m *MockServerStorageService)CreateStorageImport(r *request.CreateStorageImportRequest) (*upcloud.StorageImportDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}
func(m *MockServerStorageService)GetStorageImportDetails(r *request.GetStorageImportDetailsRequest) (*upcloud.StorageImportDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}
func(m *MockServerStorageService)WaitForStorageImportCompletion(r *request.WaitForStorageImportCompletionRequest) (*upcloud.StorageImportDetails, error) {
  args := m.Called(r)
  return args[0].(*upcloud.StorageImportDetails), args.Error(1)
}
func(m *MockServerStorageService)DeleteStorage(r *request.DeleteStorageRequest) error {
  return m.Called(r).Error(0)
}
