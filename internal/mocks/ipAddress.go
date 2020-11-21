package mocks

import (
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
  "github.com/stretchr/testify/mock"
)

type MockIpAddressService struct {
  mock.Mock
}

func (m *MockIpAddressService) GetIPAddresses() (*upcloud.IPAddresses, error) {
  args := m.Called()
  return args[0].(*upcloud.IPAddresses), args.Error(1)
}
func (m *MockIpAddressService) GetIPAddressDetails(r *request.GetIPAddressDetailsRequest) (*upcloud.IPAddress, error) {
  args := m.Called(r)
  return args[0].(*upcloud.IPAddress), args.Error(1)
}
func (m *MockIpAddressService) AssignIPAddress(r *request.AssignIPAddressRequest) (*upcloud.IPAddress, error) {
  args := m.Called(r)
  return args[0].(*upcloud.IPAddress), args.Error(1)
}
func (m *MockIpAddressService) ModifyIPAddress(r *request.ModifyIPAddressRequest) (*upcloud.IPAddress, error) {
  args := m.Called(r)
  return args[0].(*upcloud.IPAddress), args.Error(1)
}
func (m *MockIpAddressService) ReleaseIPAddress(r *request.ReleaseIPAddressRequest) error {
  args := m.Called(r)
  return args.Error(0)
}
