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
	upcloudv9 "github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const shownDatabaseUUID = "9a8effcb-80e6-4a63-a7e5-066a6d093c14"

func TestGetVersion(t *testing.T) {
	pg, valkey := "pg", "valkey"
	for _, test := range []struct {
		name     string
		expected string
		db       *upcloudv9.DatabaseServiceInformationResponse
	}{
		{
			name:     "nil database",
			expected: "",
			db:       nil,
		}, {
			name:     "nil metadata",
			expected: "",
			db:       &upcloudv9.DatabaseServiceInformationResponse{Metadata: nil},
		}, {
			name:     "pg",
			expected: "15",
			db: &upcloudv9.DatabaseServiceInformationResponse{
				Type:     &pg,
				Metadata: &map[string]any{"pg_version": "15"},
			},
		}, {
			name:     "valkey",
			expected: "8.1.5",
			db: &upcloudv9.DatabaseServiceInformationResponse{
				Type:     &valkey,
				Metadata: &map[string]any{"valkey_version": "8.1.5"},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			version := getVersion(test.db)
			assert.Equal(t, test.expected, version)
		})
	}
}

func TestShowCommand(t *testing.T) {
	text.DisableColors()
	id := shownDatabaseUUID
	database := `{
		"uuid": "` + id + `", "title": "my-pg", "name": "my-pg-prefix", "type": "pg", "zone": "fi-hel1", "state": "running",
		"plan": "rdb.standard.2x-2CPU-8GB-120GB-regular", "termination_protection": true,
		"plan_components": {"compute": {"name": "rdb.standard.2CPU-8GB", "node_count": 2, "cpu": 2, "memory_gb": 8}, "storage": {"total_gib": 120}, "backups": {"name": "regular"}},
		"metadata": {"pg_version": "16"},
		"maintenance": {
			"dow": "sunday",
			"time": "05:00:00",
			"pending_updates": [{"deadline": "", "description": "platform update", "start_after": "2026-09-30T10:00:00Z", "start_at": ""}]
		},
		"service_uri_params": {"dbname": "defaultdb", "host": "my-pg.example.com", "port": 11550, "user": "upadmin"},
		"labels": [{"key": "env", "value": "dev"}],
		"node_states": [{"name": "my-pg-1", "role": "master", "state": "running"}],
		"components": [{"component": "pg", "host": "my-pg.example.com", "port": 11550, "route": "dynamic", "usage": "primary"}]
	}`

	api := newFakeAPI(t, map[string]fakeResponse{
		"GET /1.3/database/" + id: {status: http.StatusOK, body: database},
	})
	conf := config.New()
	c := commands.BuildCommand(ShowCommand(), nil, conf)
	c.Cobra().SetArgs([]string{id})

	out, err := mockexecute.MockExecuteWithV9(c, &smock.Service{}, api.client, conf)
	require.NoError(t, err)

	for _, pattern := range []string{
		`UUID:\s+` + id,
		`Type:\s+PostgreSQL`,
		`Version:\s+16`,
		`Plan:\s+rdb.standard.2x-2CPU-8GB-120GB-regular`,
		`Termination protection:\s+yes`,
		`Plan components:`,
		`Compute:\s+rdb.standard.2CPU-8GB`,
		`Node count:\s+2`,
		`Storage per node \(GiB\):\s+120`,
		`Backups:\s+regular`,
		`Weekday:\s+sunday`,
		`Port:\s+11550`,
		`env\s+dev`,
		`my-pg-1\s+master\s+running`,
	} {
		assert.Regexp(t, pattern, out)
	}
}

func TestShowCommand_NotFound(t *testing.T) {
	id := shownDatabaseUUID
	api := newFakeAPI(t, map[string]fakeResponse{
		"GET /1.3/database/" + id: {status: http.StatusNotFound, contentType: "application/json", body: `{"error": {"error_code": "SERVICE_NOT_FOUND", "error_message": "Service not found"}}`},
	})
	conf := config.New()
	c := commands.BuildCommand(ShowCommand(), nil, conf)
	c.Cobra().SetArgs([]string{id})

	_, err := mockexecute.MockExecuteWithV9(c, &smock.Service{}, api.client, conf)
	assert.EqualError(t, err, "Service not found (type=SERVICE_NOT_FOUND, status=404)")
	assert.True(t, apierror.IsNotFound(err))
}

func TestShowCommand_MySQL(t *testing.T) {
	id := shownDatabaseUUID
	api := newFakeAPI(t, map[string]fakeResponse{
		"GET /1.3/database/" + id: {status: http.StatusOK, body: `{"uuid":"` + id + `","type":"mysql","metadata":{"mysql_version":"8"},"state":"running"}`},
	})
	conf := config.New()
	c := commands.BuildCommand(ShowCommand(), nil, conf)
	c.Cobra().SetArgs([]string{id})

	out, err := mockexecute.MockExecuteWithV9(c, &smock.Service{}, api.client, conf)
	require.NoError(t, err)
	assert.Regexp(t, `Type:\s+MySQL`, out)
	assert.Regexp(t, `Version:\s+8`, out)
}

func TestShowCommand_MachineReadableOutputPreservesLegacyResponse(t *testing.T) {
	id := shownDatabaseUUID
	response := `{"uuid":"` + id + `","title":"my-pg","type":"pg","service_uri_params":{"port":11550},"plan_components":{"compute":{"name":"rdb.standard.2CPU-8GB","node_count":2}}}`

	for _, outputFormat := range []string{config.ValueOutputJSON, config.ValueOutputYAML} {
		t.Run(outputFormat, func(t *testing.T) {
			api := newFakeAPI(t, map[string]fakeResponse{"GET /1.3/database/" + id: {status: http.StatusOK, body: response}})
			conf := config.New()
			conf.Viper().Set(config.KeyOutput, outputFormat)
			c := commands.BuildCommand(ShowCommand(), nil, conf)
			c.Cobra().SetArgs([]string{id})

			out, err := mockexecute.MockExecuteWithV9(c, &smock.Service{}, api.client, conf)
			require.NoError(t, err)
			assert.NotContains(t, out, "plan_components")
			assert.NotContains(t, out, "rdb.standard.2CPU-8GB")
			if outputFormat == config.ValueOutputJSON {
				var got map[string]any
				require.NoError(t, json.Unmarshal([]byte(out), &got))
				assert.Equal(t, id, got["uuid"])
				assert.NotContains(t, got, "plan_components")
				assert.Equal(t, "11550", got["service_uri_params"].(map[string]any)["port"])
			}
		})
	}
}

func TestShowCommand_OpenSearchAccessControl(t *testing.T) {
	text.DisableColors()
	id := shownDatabaseUUID
	api := newFakeAPI(t, map[string]fakeResponse{
		"GET /1.3/database/" + id:                     {status: http.StatusOK, body: `{"uuid": "` + id + `", "type": "opensearch", "state": "running"}`},
		"GET /1.3/database/" + id + "/access-control": {status: http.StatusOK, body: `{"access_control": true, "extended_access_control": false}`},
	})
	conf := config.New()
	c := commands.BuildCommand(ShowCommand(), nil, conf)
	c.Cobra().SetArgs([]string{id})

	out, err := mockexecute.MockExecuteWithV9(c, &smock.Service{}, api.client, conf)
	require.NoError(t, err)
	assert.Regexp(t, `Type:\s+OpenSearch`, out)
	assert.Regexp(t, `Access control:\s+yes`, out)
	assert.Regexp(t, `Extended access control:\s+no`, out)
}
