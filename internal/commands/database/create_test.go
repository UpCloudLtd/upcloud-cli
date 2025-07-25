package database

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/storage"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCommand(t *testing.T) {
	databaseDetailsMaint := upcloud.ManagedDatabase{
		UUID: "0927dfd6-3884-4079-a948-3a8881df1a7a",
	}
	serviceType := upcloud.ManagedDatabaseType{
		Properties: map[string]upcloud.ManagedDatabaseServiceProperty{
			"version": {
				Type: []string{"string", "null"},
			},
			"numeric_string": {
				Type: "string",
			},
		},
	}

	for _, test := range []struct {
		name              string
		args              []string
		createDatabaseReq request.CreateManagedDatabaseRequest
		error             string
	}{
		{
			name: "minimum required parameters",
			args: []string{
				"--title", "db-test",
				"--zone", "fi-hel1",
				"--hostname-prefix", "testdb",
			},
			createDatabaseReq: request.CreateManagedDatabaseRequest{
				Title:          "db-test",
				Zone:           "fi-hel1",
				HostNamePrefix: "testdb",
				Plan:           "2x2xCPU-4GB-100GB",
				Type:           upcloud.ManagedDatabaseServiceTypeMySQL,
			},
		},
		{
			name: "pg with label, protection, and version parameters",
			args: []string{
				"--title", "full-test",
				"--zone", "fi-hel1",
				"--hostname-prefix", "fulldb",
				"--plan", "4x4xCPU-8GB-200GB",
				"--type", "pg",
				"--label", "env=test,app=database",
				"--enable-termination-protection",
				"--property", "version=13",
			},
			createDatabaseReq: request.CreateManagedDatabaseRequest{
				Title:                 "full-test",
				Zone:                  "fi-hel1",
				HostNamePrefix:        "fulldb",
				Plan:                  "4x4xCPU-8GB-200GB",
				Type:                  upcloud.ManagedDatabaseServiceTypePostgreSQL,
				Labels:                []upcloud.Label{{Key: "env", Value: "test"}, {Key: "app", Value: "database"}},
				TerminationProtection: boolPtr(true),
				Properties: request.ManagedDatabasePropertiesRequest{
					"version": "13",
				},
			},
		},
		{
			name: "opensearch with properties parameters",
			args: []string{
				"--title", "full-test",
				"--zone", "fi-hel1",
				"--hostname-prefix", "fulldb",
				"--plan", "4x4xCPU-8GB-200GB",
				"--type", "opensearch",
				"--property", "saml={\"enabled\":true}", // without quotes
				"--property", "openid=\"{\"client_id\":\"test_client_id\"}\"", // with quotes
				"--property", "ism_enabled=true",
				"--property", "custom_domain=custom.upcloud.com",
				"--property", "numeric_string=123",
			},
			createDatabaseReq: request.CreateManagedDatabaseRequest{
				Title:          "full-test",
				Zone:           "fi-hel1",
				HostNamePrefix: "fulldb",
				Plan:           "4x4xCPU-8GB-200GB",
				Type:           upcloud.ManagedDatabaseServiceTypeOpenSearch,
				Properties: request.ManagedDatabasePropertiesRequest{
					"ism_enabled":    true,
					"custom_domain":  "custom.upcloud.com",
					"saml":           map[string]interface{}{"enabled": true},
					"openid":         map[string]interface{}{"client_id": "test_client_id"},
					"numeric_string": "123",
				},
			},
		},
		{
			name: "mysql default database type",
			args: []string{
				"--title", "mysql-test",
				"--zone", "fi-hel1",
				"--hostname-prefix", "mysqldb",
			},
			createDatabaseReq: request.CreateManagedDatabaseRequest{
				Title:          "mysql-test",
				Zone:           "fi-hel1",
				HostNamePrefix: "mysqldb",
				Plan:           "2x2xCPU-4GB-100GB",
				Type:           upcloud.ManagedDatabaseServiceTypeMySQL,
			},
		},
		{
			name: "maintenance-dow and maintenance-time",
			args: []string{
				"--title", "mysql-test",
				"--zone", "fi-hel1",
				"--hostname-prefix", "mysqldb",
				"--maintenance-dow", "monday",
				"--maintenance-time", "02:00",
			},
			createDatabaseReq: request.CreateManagedDatabaseRequest{
				Title:          "mysql-test",
				Zone:           "fi-hel1",
				HostNamePrefix: "mysqldb",
				Plan:           "2x2xCPU-4GB-100GB",
				Type:           upcloud.ManagedDatabaseServiceTypeMySQL,
				Maintenance: request.ManagedDatabaseMaintenanceTimeRequest{
					DayOfWeek: "monday",
					Time:      "02:00",
				},
			},
		},
		{
			name: "missing required title parameter",
			args: []string{
				"--zone", "fi-hel1",
				"--hostname-prefix", "testdb",
			},
			error: "required flag(s) \"title\" not set",
		},
		{
			name: "missing required zone parameter",
			args: []string{
				"--title", "db-test",
				"--hostname-prefix", "testdb",
			},
			error: "required flag(s) \"zone\" not set",
		},
		{
			name: "missing required hostname-prefix parameter",
			args: []string{
				"--title", "db-test",
				"--zone", "fi-hel1",
			},
			error: "required flag(s) \"hostname-prefix\" not set",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			storage.CachedStorages = nil
			createDatabaseReq := test.createDatabaseReq
			mService.On("CreateManagedDatabase", &createDatabaseReq).Return(&databaseDetailsMaint, nil)
			mService.On("GetManagedDatabaseServiceType", mock.Anything).Return(&serviceType, nil)

			c := commands.BuildCommand(testCmd, nil, conf)

			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			if test.error != "" {
				if err == nil {
					t.Errorf("expected error '%v', got nil", test.error)
				} else {
					assert.Equal(t, test.error, err.Error())
				}
			} else {
				assert.NoError(t, err)
				// Fix: Check for CreateManagedDatabase method call instead of CreateDatabase
				mService.AssertNumberOfCalls(t, "CreateManagedDatabase", 1)
			}
		})
	}
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}
