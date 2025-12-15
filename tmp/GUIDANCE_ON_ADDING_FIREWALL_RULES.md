# Guidance on Implementing Firewall Rule Enable/Disable for Issue #244

**Issue:** [#244 - Enable / disable specific firewall rule via cli](https://github.com/UpCloudLtd/upcloud-cli/issues/244)

**Author:** Assistant
**Date:** 2025-12-13
**Status:** Implementation Guide

---

## Executive Summary

The UpCloud API does **not support modifying individual firewall rules**. Instead, the entire firewall ruleset must be replaced atomically using the `CreateFirewallRules` endpoint. This is already implemented in the `upcloud-go-api` library and used successfully by the Terraform provider.

**No new API features are required** - the CLI just needs to implement the proper workflow.

---

## How UpCloud Firewall Rules Work

### Key Concepts

1. **No Individual Rule Modification**: The UpCloud API does not have PATCH/PUT endpoints for individual firewall rules
2. **Atomic Replacement**: All rule operations (create, update, delete) use `CreateFirewallRules` which replaces the entire ruleset
3. **Position-Based Ordering**: Rules are evaluated in numerical order (position 1, 2, 3...) with a maximum of 1000 rules per server
4. **Action Field**: Rules can be "enabled" (`action: "accept"`) or "disabled" (`action: "drop"`)

### API Endpoints Available

From `upcloud-go-api/upcloud/service/firewall.go`:

```go
type Firewall interface {
    GetFirewallRules(ctx context.Context, r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error)
    GetFirewallRuleDetails(ctx context.Context, r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error)
    CreateFirewallRule(ctx context.Context, r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error)
    CreateFirewallRules(ctx context.Context, r *request.CreateFirewallRulesRequest) error
    DeleteFirewallRule(ctx context.Context, r *request.DeleteFirewallRuleRequest) error
}
```

**Important:** `CreateFirewallRules` (plural) replaces the entire ruleset atomically.

---

## How Terraform Provider Implements This

Reference: `/tmp/terraform-provider-upcloud/internal/service/firewall/firewall.go`

### Create Operation (lines 297-329)

```go
func (r *firewallRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data firewallRulesModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    // Build the complete ruleset
    apiFirewallRules, diags := buildFirewallRules(ctx, data)

    // Create request with ALL rules
    apiReq := request.CreateFirewallRulesRequest{
        ServerUUID:    data.ServerID.ValueString(),
        FirewallRules: apiFirewallRules,
    }

    // Replace entire ruleset
    err := r.client.CreateFirewallRules(ctx, &apiReq)
}
```

### Update Operation (lines 366-396)

**Critical Insight:** Update uses the exact same `CreateFirewallRules` call!

```go
func (r *firewallRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data firewallRulesModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    // Build the MODIFIED ruleset
    apiFirewallRules, diags := buildFirewallRules(ctx, data)

    // Use same CreateFirewallRules request
    apiReq := request.CreateFirewallRulesRequest{
        ServerUUID:    data.ServerID.ValueString(),
        FirewallRules: apiFirewallRules,
    }

    // REPLACE entire ruleset with modified version
    err := r.client.CreateFirewallRules(ctx, &apiReq)
}
```

### Delete Operation (lines 398-418)

```go
func (r *firewallRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Replace ruleset with empty list
    apiReq := request.CreateFirewallRulesRequest{
        ServerUUID:    data.ServerID.ValueString(),
        FirewallRules: nil,  // Empty = delete all rules
    }

    err := r.client.CreateFirewallRules(ctx, &apiReq)
}
```

**Pattern:** All operations use `CreateFirewallRules` to atomically replace the ruleset.

---

## Recommended Implementation for CLI Issue #244

### Proposed Command Syntax

```bash
# Enable a specific firewall rule (change action to "accept")
upctl server firewall rule enable <server-uuid> --position <rule-position>

# Disable a specific firewall rule (change action to "drop")
upctl server firewall rule disable <server-uuid> --position <rule-position>

# Alternative: single command with flag
upctl server firewall rule modify <server-uuid> --position <rule-position> --enable
upctl server firewall rule modify <server-uuid> --position <rule-position> --disable
```

### Implementation Steps

#### Step 1: Fetch Current Ruleset

```go
import (
    "context"
    "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
    "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
    "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
)

func modifyFirewallRule(ctx context.Context, svc *service.Service, serverUUID string, position int, enable bool) error {
    // 1. Get all current firewall rules
    rulesReq := &request.GetFirewallRulesRequest{
        ServerUUID: serverUUID,
    }

    currentRules, err := svc.GetFirewallRules(ctx, rulesReq)
    if err != nil {
        return fmt.Errorf("failed to fetch firewall rules: %w", err)
    }
```

#### Step 2: Find and Modify Target Rule

```go
    // 2. Find the rule at the specified position
    var targetRuleIndex = -1
    for i, rule := range currentRules.FirewallRules {
        if rule.Position == position {
            targetRuleIndex = i
            break
        }
    }

    if targetRuleIndex == -1 {
        return fmt.Errorf("firewall rule at position %d not found", position)
    }

    // 3. Modify the action field
    newAction := upcloud.FirewallRuleActionDrop
    if enable {
        newAction = upcloud.FirewallRuleActionAccept
    }

    currentRules.FirewallRules[targetRuleIndex].Action = newAction
```

#### Step 3: Replace Entire Ruleset

```go
    // 4. Replace the entire ruleset with the modified version
    replaceReq := &request.CreateFirewallRulesRequest{
        ServerUUID:    serverUUID,
        FirewallRules: currentRules.FirewallRules,
    }

    err = svc.CreateFirewallRules(ctx, replaceReq)
    if err != nil {
        return fmt.Errorf("failed to update firewall rules: %w", err)
    }

    return nil
}
```

### Complete Example Implementation

```go
package firewall

import (
    "context"
    "fmt"

    "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
    "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
    "github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service"
)

// EnableFirewallRule enables a specific firewall rule by setting action to "accept"
func EnableFirewallRule(ctx context.Context, svc *service.Service, serverUUID string, position int) error {
    return modifyFirewallRuleAction(ctx, svc, serverUUID, position, upcloud.FirewallRuleActionAccept)
}

// DisableFirewallRule disables a specific firewall rule by setting action to "drop"
func DisableFirewallRule(ctx context.Context, svc *service.Service, serverUUID string, position int) error {
    return modifyFirewallRuleAction(ctx, svc, serverUUID, position, upcloud.FirewallRuleActionDrop)
}

// modifyFirewallRuleAction changes the action of a specific firewall rule
func modifyFirewallRuleAction(ctx context.Context, svc *service.Service, serverUUID string, position int, action string) error {
    // Validate action
    if action != upcloud.FirewallRuleActionAccept && action != upcloud.FirewallRuleActionDrop {
        return fmt.Errorf("invalid action: %s (must be 'accept' or 'drop')", action)
    }

    // 1. Fetch current ruleset
    rulesReq := &request.GetFirewallRulesRequest{
        ServerUUID: serverUUID,
    }

    currentRules, err := svc.GetFirewallRules(ctx, rulesReq)
    if err != nil {
        return fmt.Errorf("failed to fetch firewall rules: %w", err)
    }

    // 2. Find the target rule by position
    ruleFound := false
    for i := range currentRules.FirewallRules {
        if currentRules.FirewallRules[i].Position == position {
            // 3. Modify the action
            currentRules.FirewallRules[i].Action = action
            ruleFound = true
            break
        }
    }

    if !ruleFound {
        return fmt.Errorf("firewall rule at position %d not found on server %s", position, serverUUID)
    }

    // 4. Replace entire ruleset atomically
    replaceReq := &request.CreateFirewallRulesRequest{
        ServerUUID:    serverUUID,
        FirewallRules: currentRules.FirewallRules,
    }

    err = svc.CreateFirewallRules(ctx, replaceReq)
    if err != nil {
        return fmt.Errorf("failed to update firewall rules: %w", err)
    }

    return nil
}
```

---

## Integration Points in upcloud-cli

### File Structure

Based on typical CLI structure, you'll likely need to modify/create:

```
internal/commands/server/
├── firewall/
│   ├── firewall.go          # Main firewall command group
│   ├── rule_enable.go       # New: enable rule command
│   ├── rule_disable.go      # New: disable rule command
│   └── rule_modify.go       # Alternative: single modify command
```

### Command Registration

Follow the existing pattern in `internal/commands/server/server.go`:

```go
import (
    "github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/server/firewall"
)

func BuildCommands() []commands.Command {
    return []commands.Command{
        // ... existing commands
        firewall.BaseFirewallCommand(),  // If you create a new subcommand group
    }
}
```

### Error Handling

Ensure proper error messages for common scenarios:

```go
if err != nil {
    switch {
    case strings.Contains(err.Error(), "not found"):
        return fmt.Errorf("server or firewall rule not found: %w", err)
    case strings.Contains(err.Error(), "UNAUTHORIZED"):
        return fmt.Errorf("authentication failed: check your credentials")
    case strings.Contains(err.Error(), "FORBIDDEN"):
        return fmt.Errorf("insufficient permissions to modify firewall rules")
    default:
        return fmt.Errorf("failed to modify firewall rule: %w", err)
    }
}
```

---

## Testing Recommendations

### Unit Tests

Test the rule modification logic in isolation:

```go
func TestModifyFirewallRuleAction(t *testing.T) {
    // Mock service
    mockSvc := &mockService{
        rules: &upcloud.FirewallRules{
            FirewallRules: []upcloud.FirewallRule{
                {Position: 1, Action: "accept", Direction: "in"},
                {Position: 2, Action: "drop", Direction: "out"},
            },
        },
    }

    // Test enabling a rule
    err := modifyFirewallRuleAction(context.Background(), mockSvc, "test-server", 2, "accept")
    assert.NoError(t, err)
    assert.Equal(t, "accept", mockSvc.lastReplacedRules[1].Action)
}
```

### Integration Tests

Test against actual UpCloud API (or use recorded fixtures):

1. Create a test server
2. Add firewall rules
3. Enable/disable specific rules
4. Verify the changes
5. Clean up

### Edge Cases to Test

1. Rule position doesn't exist
2. Server doesn't exist
3. Empty ruleset
4. Position 1 (first rule)
5. Last rule in the list
6. Server with maximum rules (1000)

---

## Performance Considerations

### Optimization Opportunities

1. **Caching:** If modifying multiple rules, fetch once and replace once
2. **Validation:** Validate position exists before making API call
3. **Dry Run:** Add `--dry-run` flag to show what would change without applying

### Example: Batch Modifications

```go
// For modifying multiple rules efficiently
func ModifyMultipleRules(ctx context.Context, svc *service.Service, serverUUID string, modifications map[int]string) error {
    // Fetch once
    currentRules, err := svc.GetFirewallRules(ctx, &request.GetFirewallRulesRequest{
        ServerUUID: serverUUID,
    })
    if err != nil {
        return err
    }

    // Apply all modifications
    for position, action := range modifications {
        for i := range currentRules.FirewallRules {
            if currentRules.FirewallRules[i].Position == position {
                currentRules.FirewallRules[i].Action = action
                break
            }
        }
    }

    // Replace once
    return svc.CreateFirewallRules(ctx, &request.CreateFirewallRulesRequest{
        ServerUUID:    serverUUID,
        FirewallRules: currentRules.FirewallRules,
    })
}
```

---

## Alternative Approaches

### Option 1: Direct Action Toggle

User specifies position, CLI automatically toggles between accept/drop:

```bash
upctl server firewall rule toggle <server-uuid> --position 10
```

### Option 2: Full Rule Update

Allow modifying any field, not just action:

```bash
upctl server firewall rule update <server-uuid> --position 10 \
    --action drop \
    --comment "Updated rule"
```

### Option 3: Rule Replace by Position

Replace entire rule definition at position:

```bash
upctl server firewall rule replace <server-uuid> --position 10 \
    --action accept \
    --direction in \
    --protocol tcp \
    --destination-port-start 80 \
    --destination-port-end 80
```

---

## User Experience Considerations

### Before/After Display

Show the user what changed:

```bash
$ upctl server firewall rule disable 00000000-0000-0000-0000-000000000000 --position 10

Firewall rule at position 10:
  Direction:      in
  Protocol:       tcp
  Port:           80
  Action:         accept → drop

✓ Successfully disabled firewall rule at position 10
```

### Interactive Confirmation

For safety, consider requiring confirmation:

```bash
$ upctl server firewall rule disable <server-uuid> --position 10

This will disable the following firewall rule:
  Position:       10
  Direction:      in
  Protocol:       tcp
  Destination:    0.0.0.0/0:80
  Current Action: accept
  New Action:     drop

This may affect network connectivity to your server.
Continue? [y/N]:
```

### List View with Status

Enhance the list command to show rule status:

```bash
$ upctl server firewall rule list <server-uuid>

POS  STATUS   DIR  PROTOCOL  SOURCE         DEST           ACTION  COMMENT
1    enabled  in   tcp       0.0.0.0/0      *:22           accept  SSH access
2    enabled  in   tcp       0.0.0.0/0      *:80           accept  HTTP
3    enabled  in   tcp       0.0.0.0/0      *:443          accept  HTTPS
10   disabled in   tcp       0.0.0.0/0      *:8080         drop    Dev port
```

---

## Related Issues and PRs

### Related CLI Issues

- [#243 - Enable firewall on server provision](https://github.com/UpCloudLtd/upcloud-cli/issues/243)
- [#100 - Add Automatic Upcloud DNS firewall rules](https://github.com/UpCloudLtd/upcloud-cli/issues/100)

### API Library PRs

- [#449 - firewall field type fix](https://github.com/UpCloudLtd/upcloud-go-api/issues/449) - Recently fixed, ensures billing API compatibility

---

## References

### Source Code

- **API Library:** `github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/service/firewall.go`
- **Terraform Provider:** `/tmp/terraform-provider-upcloud/internal/service/firewall/firewall.go`
  - Create operation: lines 297-329
  - Update operation: lines 366-396
  - Delete operation: lines 398-418

### Documentation

- [UpCloud API - Firewall Rules](https://developers.upcloud.com/1.2/11-firewall/)
- [UpCloud Firewall Guide](https://upcloud.com/docs/guides/managing-firewall/)
- [Terraform UpCloud Firewall Rules](https://registry.terraform.io/providers/UpCloudLtd/upcloud/latest/docs/resources/firewall_rules)

### API Constants

From `upcloud-go-api/v8/upcloud/firewall.go`:

```go
const (
    FirewallRuleActionAccept = "accept"
    FirewallRuleActionDrop   = "drop"
    FirewallRuleDirectionIn  = "in"
    FirewallRuleDirectionOut = "out"
    FirewallRuleProtocolTCP  = "tcp"
    FirewallRuleProtocolUDP  = "udp"
    FirewallRuleProtocolICMP = "icmp"
)
```

---

## Summary

**Key Takeaway:** Modifying individual firewall rules requires fetching the entire ruleset, modifying the target rule in memory, and replacing the entire ruleset atomically using `CreateFirewallRules`.

This pattern is proven and battle-tested by the Terraform provider. The implementation in the CLI should follow the same approach for consistency and reliability.

**Estimated Implementation Effort:** 2-4 hours for basic enable/disable commands with tests.

---

**Questions?** Review the Terraform provider implementation in `/tmp/terraform-provider-upcloud/internal/service/firewall/firewall.go` for a working reference implementation.
