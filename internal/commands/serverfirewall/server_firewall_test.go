package serverfirewall

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/mock"
)

type MockFirewallRuleService struct {
	mock.Mock
}

func (m *MockFirewallRuleService) GetFirewallRule() (*upcloud.FirewallRules, error) {
	args := m.Called(r)
	return args[0].(*upcloud.FirewallRules), args.Error(1)
}

func (m *MockFirewallRuleService) CreateFirewallRule(r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRules, error) {
	args := m.Called(r)
	return args[0].(*upcloud.FirewallRules), args.Error(1)
}

func (m *MockFirewallRuleService) DeleteFirewallRule(r *request.DeleteFirewallRuleRequest) (*upcloud.FirewallRules, error) {
	args := m.Called(r)
	return args.Error(0)
}

var (
	Title1 = "mock-storage-title1"
	UUID1  = "0127dfd6-3884-4079-a948-3a8881df1a7a"
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
