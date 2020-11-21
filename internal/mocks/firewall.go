package mocks

import (
  "github.com/UpCloudLtd/upcloud-go-api/upcloud"
  "github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
  "github.com/stretchr/testify/mock"
)

type MockFirewallService struct {
  mock.Mock
}

func(m *MockFirewallService) GetFirewallRules(r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error) {
  args := m.Called(r)
  return args[0].(*upcloud.FirewallRules), args.Error(1)
}
func(m *MockFirewallService) GetFirewallRuleDetails(r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error) {
  args := m.Called(r)
  return args[0].(*upcloud.FirewallRule), args.Error(1)
}
func(m *MockFirewallService) CreateFirewallRule(r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error) {
  args := m.Called(r)
  return args[0].(*upcloud.FirewallRule), args.Error(1)
}
func(m *MockFirewallService) CreateFirewallRules(r *request.CreateFirewallRulesRequest) error {
  args := m.Called(r)
  return args.Error(0)
}
func(m *MockFirewallService) DeleteFirewallRule(r *request.DeleteFirewallRuleRequest) error {
  args := m.Called(r)
  return args.Error(0)
}
