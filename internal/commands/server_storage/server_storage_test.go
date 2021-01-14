package server_storage

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/mock"
)

type MockServerService struct {
	mock.Mock
}

func (m *MockServerService) GetServerConfigurations() (*upcloud.ServerConfigurations, error) {
	args := m.Called()
	return args[0].(*upcloud.ServerConfigurations), args.Error(1)
}
func (m *MockServerService) GetServers() (*upcloud.Servers, error) {
	args := m.Called()
	return args[0].(*upcloud.Servers), args.Error(1)
}
func (m *MockServerService) GetServerDetails(r *request.GetServerDetailsRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) CreateServer(r *request.CreateServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) WaitForServerState(r *request.WaitForServerStateRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) StartServer(r *request.StartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) StopServer(r *request.StopServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) RestartServer(r *request.RestartServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) ModifyServer(r *request.ModifyServerRequest) (*upcloud.ServerDetails, error) {
	args := m.Called(r)
	return args[0].(*upcloud.ServerDetails), args.Error(1)
}
func (m *MockServerService) DeleteServer(r *request.DeleteServerRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
func (m *MockServerService) DeleteServerAndStorages(r *request.DeleteServerAndStoragesRequest) error {
	args := m.Called(r)
	return args.Error(0)
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

var (
	Title1 = "mock-storage-title1"
	Title2 = "mock-storage-title2"
	Title3 = "mock-storage-title3"
	Uuid1  = "0127dfd6-3884-4079-a948-3a8881df1a7a"
	Uuid2  = "012bde1d-f0e7-4bb2-9f4a-74e1f2b49c07"
	Uuid3  = "012c61a6-b8f0-48c2-a63a-b4bf7d26a655"
	Uuid4  = "012c61a6-er4g-mf2t-b63a-b4be4326a655"
)

const mockResponse = "mock-response"
const mockRequest = "mock-request"

type MockHandler struct{}

func (s MockHandler) Handle(requests []interface{}) (interface{}, error) {
	for _, r := range requests {
		if r != mockRequest {
			return nil, fmt.Errorf("upexpected request %q", r)
		}
	}
	return mockResponse, nil
}
