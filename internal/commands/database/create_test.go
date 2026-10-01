package database

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/apierror"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const createdDatabaseUUID = "0927dfd6-3884-4079-a948-3a8881df1a7a"

const planCatalog = `{"service_types": [
	{"type": "pg", "componentised": true, "node_counts": [1, 2, 3], "zones": ["fi-hel1"], "compute_shapes": [
		{"compute": "rdb.standard.2CPU-8GB", "family": "standard", "cpu": 2, "memory_gb": 8, "dynamic_storage_supported": true,
			"node_counts": [1, 2, 3], "backups": ["regular", "extended"],
			"storage": {"dynamic_max_multiplier": 4, "step_gib": 10, "total_cap_gib": 10240, "options": [{"base_gib": 10, "max_gib": 50}, {"base_gib": 80, "max_gib": 400}]}},
		{"compute": "rdb.development.1CPU-1GB", "family": "development", "cpu": 1, "memory_gb": 1, "dynamic_storage_supported": true,
			"node_counts": [1], "backups": ["mini"],
			"storage": {"dynamic_max_multiplier": 4, "step_gib": 10, "total_cap_gib": null, "options": [{"base_gib": 10, "max_gib": 50}]}}
	]},
	{"type": "mysql", "componentised": true, "node_counts": [1, 2, 3], "zones": ["fi-hel1"], "compute_shapes": [
		{"compute": "rdb.standard.2CPU-8GB", "family": "standard", "cpu": 2, "memory_gb": 8, "dynamic_storage_supported": true,
			"node_counts": [1, 2, 3], "backups": ["regular", "extended"],
			"storage": {"dynamic_max_multiplier": 4, "step_gib": 10, "total_cap_gib": 10240, "options": [{"base_gib": 80, "max_gib": 400}]}}
	]}
]}`

const createdDatabase = `{"uuid": "` + createdDatabaseUUID + `", "title": "db-test", "type": "pg", "state": "rebuilding", "plan": "rdb.standard.2x-2CPU-8GB-100GB-regular", "plan_components": {"compute": {"name": "rdb.standard.2CPU-8GB"}}}`

// serviceType uses a list type for version, as the live API does.
const serviceType = `{"name": "pg", "properties": {
	"version": {"type": ["string", "null"], "title": "Version"},
	"numeric_string": {"type": "string", "title": "Numeric string"},
	"max_connections": {"type": "integer", "title": "Max connections"}
}}`

func newCreateAPI(t *testing.T, create fakeResponse) *fakeAPI {
	responses := map[string]fakeResponse{
		"GET /1.3/database/plans": {status: http.StatusOK, body: planCatalog},
		"POST /1.3/database":      create,
		"GET /1.3/database/" + createdDatabaseUUID: {sequence: []fakeResponse{
			{status: http.StatusOK, body: `{"uuid": "` + createdDatabaseUUID + `", "state": "rebuilding"}`},
			{status: http.StatusOK, body: `{"uuid": "` + createdDatabaseUUID + `", "state": "running"}`},
		}},
	}
	for _, dbType := range []string{"mysql", "pg", "opensearch", "valkey"} {
		responses["GET /1.3/database/service-types/"+dbType] = fakeResponse{status: http.StatusOK, body: serviceType}
	}
	return newFakeAPI(t, responses)
}

func runCreate(t *testing.T, api *fakeAPI, conf *config.Config, args ...string) (string, error) {
	t.Helper()
	c := commands.BuildCommand(CreateCommand(), nil, conf)
	c.Cobra().SetArgs(args)
	return mockexecute.MockExecuteWithV9(c, &smock.Service{}, api.client, conf)
}

var requiredCreateArgs = []string{"--title", "db-test", "--zone", "fi-hel1", "--hostname-prefix", "testdb"}

func TestCreateCommand_Request(t *testing.T) {
	for _, test := range []struct {
		name       string
		args       []string
		expected   string
		components bool
	}{
		{
			name: "no plan flag uses the component defaults",
			expected: `{"hostname_prefix": "testdb", "title": "db-test", "type": "mysql", "zone": "fi-hel1",
				"plan_compute": "rdb.standard.2CPU-8GB", "plan_node_count": 2, "plan_storage_gib": 100, "plan_backups": "regular"}`,
			components: true,
		},
		{
			name: "explicit legacy plan with labels, termination protection and string property",
			args: []string{"--type", "pg", "--plan", "4x4xCPU-8GB-200GB", "--label", "env=test,app=database", "--enable-termination-protection", "--property", "version=13", "--property", "max_connections=200"},
			expected: `{"hostname_prefix": "testdb", "plan": "4x4xCPU-8GB-200GB", "title": "db-test", "type": "pg", "zone": "fi-hel1",
				"labels": [{"key": "env", "value": "test"}, {"key": "app", "value": "database"}],
				"termination_protection": true, "properties": {"version": "13", "max_connections": 200}}`,
		},
		{
			name: "opensearch with typed properties",
			args: []string{
				"--type", "opensearch", "--plan", "1x2xCPU-4GB-80GB-1D",
				"--property", `saml={"enabled":true}`,
				"--property", `openid="{"client_id":"test_client_id"}"`,
				"--property", "ism_enabled=true",
				"--property", "custom_domain=custom.upcloud.com",
				"--property", "numeric_string=123",
			},
			expected: `{"hostname_prefix": "testdb", "plan": "1x2xCPU-4GB-80GB-1D", "title": "db-test", "type": "opensearch", "zone": "fi-hel1",
				"properties": {"saml": {"enabled": true}, "openid": {"client_id": "test_client_id"}, "ism_enabled": true, "custom_domain": "custom.upcloud.com", "numeric_string": "123"}}`,
		},
		{
			name: "maintenance and networks",
			args: []string{
				"--maintenance-dow", "monday", "--maintenance-time", "02:00:00",
				"--network", "name=net-1,family=IPv4,type=private,uuid=030e83d2-d413-4d19-b1c9-af05cdb60c1f",
			},
			expected: `{"hostname_prefix": "testdb", "title": "db-test", "type": "mysql", "zone": "fi-hel1",
				"plan_compute": "rdb.standard.2CPU-8GB", "plan_node_count": 2, "plan_storage_gib": 100, "plan_backups": "regular",
				"maintenance": {"dow": "monday", "time": "02:00:00"},
				"networks": [{"family": "IPv4", "name": "net-1", "type": "private", "uuid": "030e83d2-d413-4d19-b1c9-af05cdb60c1f"}]}`,
			components: true,
		},
		{
			name: "component selection uses defaults for omitted components",
			args: []string{"--type", "pg", "--plan-storage-gib", "120"},
			expected: `{"hostname_prefix": "testdb", "title": "db-test", "type": "pg", "zone": "fi-hel1",
				"plan_compute": "rdb.standard.2CPU-8GB", "plan_node_count": 2, "plan_storage_gib": 120, "plan_backups": "regular"}`,
			components: true,
		},
		{
			name: "component selection with every component",
			args: []string{"--type", "mysql", "--plan-compute", "rdb.standard.2CPU-8GB", "--plan-node-count", "3", "--plan-storage-gib", "400", "--plan-backups", "extended"},
			expected: `{"hostname_prefix": "testdb", "title": "db-test", "type": "mysql", "zone": "fi-hel1",
				"plan_compute": "rdb.standard.2CPU-8GB", "plan_node_count": 3, "plan_storage_gib": 400, "plan_backups": "extended"}`,
			components: true,
		},
		{
			name: "component selection on the base storage of a shape",
			args: []string{"--type", "pg", "--plan-compute", "rdb.development.1CPU-1GB", "--plan-node-count", "1", "--plan-storage-gib", "10", "--plan-backups", "mini"},
			expected: `{"hostname_prefix": "testdb", "title": "db-test", "type": "pg", "zone": "fi-hel1",
				"plan_compute": "rdb.development.1CPU-1GB", "plan_node_count": 1, "plan_storage_gib": 10, "plan_backups": "mini"}`,
			components: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			api := newCreateAPI(t, fakeResponse{status: http.StatusCreated, body: createdDatabase})

			_, err := runCreate(t, api, config.New(), append(requiredCreateArgs, test.args...)...)
			require.NoError(t, err)

			creates := api.requestsTo(http.MethodPost, "/1.3/database")
			require.Len(t, creates, 1)
			assert.JSONEq(t, test.expected, string(creates[0].Body))

			if test.components {
				var body map[string]any
				require.NoError(t, json.Unmarshal(creates[0].Body, &body))
				assert.NotContains(t, body, "plan", "component selection must not send an implicit plan")
				assert.NotContains(t, body, "additional_disk_space_gib", "plan_storage_gib is exclusive with additional_disk_space_gib")
			}
		})
	}
}

func TestCreateCommand_LegacyPlanDeprecationWarning(t *testing.T) {
	for _, test := range []struct {
		name        string
		dbType      string
		wantWarning bool
	}{
		{name: "pg", dbType: "pg", wantWarning: true},
		{name: "mysql", dbType: "mysql", wantWarning: true},
		{name: "opensearch", dbType: "opensearch", wantWarning: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			api := newCreateAPI(t, fakeResponse{status: http.StatusCreated, body: createdDatabase})
			out, err := runCreate(t, api, config.New(), append(requiredCreateArgs, "--type", test.dbType, "--plan", "legacy-plan")...)
			require.NoError(t, err)

			const warning = "Deprecation Warning: --plan is deprecated for pg and mysql"
			if test.wantWarning {
				assert.Contains(t, out, warning)
			} else {
				assert.NotContains(t, out, warning)
			}
		})
	}
}

func TestCreateCommand_Errors(t *testing.T) {
	for _, test := range []struct {
		name  string
		args  []string
		error string
		// beforeAPI marks errors that must be reported without any API request.
		beforeAPI bool
	}{
		{
			name:      "plan and a component flag are mutually exclusive",
			args:      append(requiredCreateArgs, "--type", "pg", "--plan", "2x2xCPU-4GB-100GB", "--plan-compute", "rdb.standard.2CPU-8GB"),
			error:     "--plan cannot be combined with component plan flags (--plan-compute, --plan-node-count, --plan-storage-gib, or --plan-backups)",
			beforeAPI: true,
		},
		{
			name:      "invalid numeric value",
			args:      append(requiredCreateArgs, "--type", "pg", "--plan-node-count", "two"),
			error:     `--plan-node-count must be a positive integer, got "two"`,
			beforeAPI: true,
		},
		{
			name:      "non-positive node count",
			args:      append(requiredCreateArgs, "--type", "pg", "--plan-node-count", "0"),
			error:     `--plan-node-count must be a positive integer, got "0"`,
			beforeAPI: true,
		},
		{
			name:      "negative storage",
			args:      append(requiredCreateArgs, "--type", "pg", "--plan-storage-gib", "-10"),
			error:     `--plan-storage-gib must be a positive integer, got "-10"`,
			beforeAPI: true,
		},
		{
			name:      "component flags for an engine without component plans",
			args:      append(requiredCreateArgs, "--type", "opensearch", "--plan-compute", "rdb.standard.2CPU-8GB"),
			error:     `the --plan-* flags are only supported for pg and mysql, use --plan for database type "opensearch"`,
			beforeAPI: true,
		},
		{
			name:      "engine without component plans requires --plan",
			args:      append(requiredCreateArgs, "--type", "valkey"),
			error:     `--plan is required for database type "valkey", run "upctl database plans valkey" to list available plans`,
			beforeAPI: true,
		},
		{
			name:      "maintenance flags must be used together",
			args:      append(requiredCreateArgs, "--maintenance-dow", "monday"),
			error:     "if any flags in the group [maintenance-dow maintenance-time] are set they must all be set; missing [maintenance-time]",
			beforeAPI: true,
		},
		{
			name:  "invalid network uuid",
			args:  append(requiredCreateArgs, "--network", "name=net-1,family=IPv4,type=private,uuid=not-a-uuid"),
			error: `invalid networks: invalid network uuid "not-a-uuid"`,
		},
		{
			name:  "storage between the base options of a shape",
			args:  append(requiredCreateArgs, "--type", "pg", "--plan-storage-gib", "60"),
			error: "--plan-storage-gib 60 is not available for compute shape rdb.standard.2CPU-8GB, valid totals per node are 10-50, 80-400 GiB in steps of 10 GiB",
		},
		{
			name:  "storage not aligned to the storage step",
			args:  append(requiredCreateArgs, "--type", "pg", "--plan-storage-gib", "125"),
			error: "--plan-storage-gib 125 is not available for compute shape rdb.standard.2CPU-8GB",
		},
		{
			name:  "default node count not offered by the compute shape",
			args:  append(requiredCreateArgs, "--type", "pg", "--plan-compute", "rdb.development.1CPU-1GB"),
			error: "--plan-node-count 2 is not available for compute shape rdb.development.1CPU-1GB, valid node counts are 1",
		},
		{
			name:  "unknown compute shape",
			args:  append(requiredCreateArgs, "--type", "pg", "--plan-compute", "rdb.standard.99CPU-1GB"),
			error: `--plan-compute "rdb.standard.99CPU-1GB" is not available for pg, run "upctl database plans pg" to list valid values`,
		},
		{
			name:  "backup tier not offered by the compute shape",
			args:  append(requiredCreateArgs, "--type", "pg", "--plan-backups", "mini"),
			error: `--plan-backups "mini" is not available for compute shape rdb.standard.2CPU-8GB, valid backup tiers are regular, extended`,
		},
		{
			name:      "missing required title",
			args:      []string{"--zone", "fi-hel1", "--hostname-prefix", "testdb"},
			error:     `required flag(s) "title" not set`,
			beforeAPI: true,
		},
		{
			name:      "missing required zone",
			args:      []string{"--title", "db-test", "--hostname-prefix", "testdb"},
			error:     `required flag(s) "zone" not set`,
			beforeAPI: true,
		},
		{
			name:      "missing required hostname-prefix",
			args:      []string{"--title", "db-test", "--zone", "fi-hel1"},
			error:     `required flag(s) "hostname-prefix" not set`,
			beforeAPI: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			api := newCreateAPI(t, fakeResponse{status: http.StatusCreated, body: createdDatabase})

			_, err := runCreate(t, api, config.New(), test.args...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.error)
			assert.Empty(t, api.requestsTo(http.MethodPost, "/1.3/database"), "the database must not be created")
			if test.beforeAPI {
				assert.Empty(t, api.all(), "the error must be reported before calling the API")
			}
		})
	}
}

func TestCreateCommand_Output(t *testing.T) {
	api := newCreateAPI(t, fakeResponse{status: http.StatusCreated, body: createdDatabase})
	out, err := runCreate(t, api, config.New(), requiredCreateArgs...)
	require.NoError(t, err)
	assert.Contains(t, out, createdDatabaseUUID)

	conf := config.New()
	conf.Viper().Set(config.KeyOutput, config.ValueOutputJSON)
	out, err = runCreate(t, api, conf, requiredCreateArgs...)
	require.NoError(t, err)
	var body map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &body))
	assert.Equal(t, createdDatabaseUUID, body["uuid"])
	assert.Equal(t, "rdb.standard.2x-2CPU-8GB-100GB-regular", body["plan"])
	assert.NotContains(t, body, "plan_components")
}

func TestCreateCommand_Wait(t *testing.T) {
	api := newCreateAPI(t, fakeResponse{status: http.StatusCreated, body: createdDatabase})

	out, err := runCreate(t, api, config.New(), append(requiredCreateArgs, "--wait")...)
	require.NoError(t, err)
	assert.Contains(t, out, createdDatabaseUUID)
	assert.Len(t, api.requestsTo(http.MethodGet, "/1.3/database/"+createdDatabaseUUID), 2, "polls until the database is running")
}

func TestCreateCommand_APIError(t *testing.T) {
	api := newCreateAPI(t, fakeResponse{
		status:      http.StatusUnprocessableEntity,
		contentType: "application/problem+json",
		body:        `{"type": "https://developers.upcloud.com/1.3/errors#ERROR_INVALID_PLAN", "title": "Plan is not available", "status": 422, "correlation_id": "abc"}`,
	})

	_, err := runCreate(t, api, config.New(), append(requiredCreateArgs, "--type", "pg", "--plan-storage-gib", "120")...)
	require.Error(t, err)
	assert.EqualError(t, err, "Plan is not available (type=https://developers.upcloud.com/1.3/errors#ERROR_INVALID_PLAN, status=422, correlation_id=abc)")
	var problem *apierror.Problem
	require.ErrorAs(t, err, &problem)
	assert.Equal(t, "INVALID_PLAN", problem.ErrorCode())
	assert.False(t, apierror.IsNotFound(err))
}
