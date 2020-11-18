package interfaces

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type Firewall interface {
	GetFirewallRules(r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error)
	GetFirewallRuleDetails(r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error)
	CreateFirewallRule(r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error)
	CreateFirewallRules(r *request.CreateFirewallRulesRequest) error
	DeleteFirewallRule(r *request.DeleteFirewallRuleRequest) error
}
