# Add pricing information to server plans command

## Summary

This PR adds `--cost` and `--duration` flags to the `upctl server plans` command, allowing users to view pricing information for available server plans with flexible duration options.

**Closes**: #339 (partially - provides plan pricing foundation for billing summaries)

## Features

### `--cost` Flag
- Display pricing information fetched from UpCloud API
- Uses `GetPriceZones()` to retrieve current pricing
- Gracefully handles missing prices (returns 0)
- Works with all plan types including GPU plans

### `--duration` Flag
- Flexible duration parsing supporting multiple formats:
  - **Go duration format**: `1h`, `30m`, `24h`, `3600s`, etc.
  - **Friendly unit names**: `hour`, `day`, `week`, `month`, `year`
  - **Numeric + unit**: `3hours`, `10days`, `1week`, `15minutes`, `1month`, `2months`, `1year`, etc.
  - **Decimal periods**: `2.5hours`
- Default: `1h` (hourly pricing)
- Dynamic cost calculation for requested duration

## Examples

```bash
# View hourly pricing for all plans
upctl server plans --cost

# View 3-hour period costs
upctl server plans --cost --duration 3hours

# View 10-day period costs
upctl server plans --cost --duration 10days

# View monthly costs
upctl server plans --cost --duration 1month

# View annual costs
upctl server plans --cost --duration 1year

# Go duration format
upctl server plans --cost --duration 24h
```

## Sample Output

```
General purpose
Name               Cores Memory Storage size Storage tier Transfer out (GiB/month) Cost (per hour)
1xCPU-1GB          1     1024   25           maxiops      100                      0.0049
2xCPU-2GB          2     2048   50           maxiops      200                      0.0098
4xCPU-4GB          4     4096   100          maxiops      400                      0.0196
...

GPU
Name               Cores Memory Storage size Storage tier Transfer out (GiB/month) GPU model GPU amount Cost (per hour)
GPU-1xL40S         10    30720  250          maxiops      1000                     L40S       1          0.4900
...
```

## Technical Details

### Implementation Approach
- Uses reflection to dynamically map plan names to price zone fields
- Field naming: `ServerPlan` + plan name with "-" removed (e.g., "2xCPU-2GB" → `ServerPlan2xCPU2GB`)
- Price from API is hourly; duration used to calculate final cost
- Smart column headers based on duration selected

### Duration Parsing
- Supports Go's standard `time.ParseDuration` format
- Extends with friendly names and numeric unit patterns
- Case-insensitive parsing
- Comprehensive error messages for invalid formats

### Cost Calculation
```
final_cost = hourly_price * duration.Hours()
```

## Testing

✅ **Build**: Successful (no errors or warnings)
✅ **Unit Tests**: All pass (existing server command tests)
✅ **Duration Parsing**: 13/13 format combinations tested
✅ **Feature Tests**: Cost calculation, output formatting, all plan types

See `TEST_SUMMARY.md` for detailed test results.

## Backward Compatibility

✅ Fully backward compatible
- Both flags are optional
- Default behavior (no pricing) unchanged
- Existing output modes (JSON, human) unaffected
- All existing plan categorization and sorting preserved

## Related Issues

- **#339** - Support getting billing summary (enhancement)
  - This PR provides the foundation for plan pricing display
  - Future enhancement could extend to resource-level billing summaries as requested in #339

## Future Enhancements

This PR establishes the pricing infrastructure for potential future features:
- Resource-level billing summaries
- Multi-zone pricing comparison
- Cost breakdown by resource type (CPU, memory, storage)
- Billing filters and aggregation

## Files Changed

- `internal/commands/server/plan_list.go` (+203, -8)
  - Added `InitCommand()` for flag registration
  - Added `getPlanCost()` for price retrieval
  - Added `parseDuration()` for flexible duration parsing
  - Added `formatDurationHeader()` for smart column headers
  - Added `formatDuration()` for human-readable duration display
  - Updated `ExecuteWithoutArguments()` to fetch and display pricing
  - Updated `planSection()` to conditionally include cost column

- `go.mod`, `go.sum` - No new dependencies added
