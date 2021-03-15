package serverfirewall_test

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/mock"
)

type MockFirewallRuleService struct {
	mock.Mock
}

func (m *MockFirewallRuleService) GetFirewallRules(r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error) {
	args := m.Called(r)
	return args[0].(*upcloud.FirewallRules), args.Error(1)
}

func (m *MockFirewallRuleService) GetFirewallRuleDetails(r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error) {
	args := m.Called(r)
	return args[0].(*upcloud.FirewallRule), args.Error(1)
}

func (m *MockFirewallRuleService) CreateFirewallRule(r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error) {
	args := m.Called(r)
	return args[0].(*upcloud.FirewallRule), args.Error(1)
}

func (m *MockFirewallRuleService) CreateFirewallRules(r *request.CreateFirewallRulesRequest) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockFirewallRuleService) DeleteFirewallRule(r *request.DeleteFirewallRuleRequest) error {
	args := m.Called(r)
	return args.Error(0)
}
