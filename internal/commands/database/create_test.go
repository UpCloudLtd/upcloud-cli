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
)

var (
	Title1             = "mock-storage-title1"
	Title2             = "mock-storage-title2"
	UUID1              = "0127dfd6-3884-4079-a948-3a8881df1a7a"
	UUID3              = "012c61a6-b8f0-48c2-a63a-b4bf7d26a655"
	PrivateNetworkUUID = "03b5b0a0-ad4c-4817-9632-dafdb3ace5d9"
	MockPrivateIPv4    = "10.0.0.1"
	MockPrivateIPv6    = "fd42:42::1"
	MockPublicIPv4     = "192.0.2.0"
	MockPublicIPv6     = "2001:DB8::1"
)

func TestCreateCommand(t *testing.T) {
	databaseDetailsMaint := upcloud.ManagedDatabase{
		UUID: UUID1,
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
				"--title=db-test",
				"--zone=fi-hel1",
				"--host-name-prefix=testdb",
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
			name: "pg with label and protection parameters",
			args: []string{
				"--title=full-test",
				"--zone=fi-hel1",
				"--host-name-prefix=fulldb",
				"--plan=4x4xCPU-8GB-200GB",
				"--type=pg",
				"--label=env=test,app=database",
				"--enable-termination-protection",
			},
			createDatabaseReq: request.CreateManagedDatabaseRequest{
				Title:                 "full-test",
				Zone:                  "fi-hel1",
				HostNamePrefix:        "fulldb",
				Plan:                  "4x4xCPU-8GB-200GB",
				Type:                  upcloud.ManagedDatabaseServiceTypePostgreSQL,
				Labels:                []upcloud.Label{{Key: "env", Value: "test"}, {Key: "app", Value: "database"}},
				TerminationProtection: boolPtr(true),
			},
		},
		{
			name: "opensearch with properties parameters",
			args: []string{
				"--title=full-test",
				"--zone=fi-hel1",
				"--host-name-prefix=fulldb",
				"--plan=4x4xCPU-8GB-200GB",
				"--type=opensearch",
				"--property=saml={\"enabled\":true}", // without quotes
				"--property=openid=\"{\"client_id\":\"test_client_id\"}\"", // with quotes
				"--property=ism_enabled=true",
				"--property=custom_domain=custom.upcloud.com",
			},
			createDatabaseReq: request.CreateManagedDatabaseRequest{
				Title:          "full-test",
				Zone:           "fi-hel1",
				HostNamePrefix: "fulldb",
				Plan:           "4x4xCPU-8GB-200GB",
				Type:           upcloud.ManagedDatabaseServiceTypeOpenSearch,
				Properties: request.ManagedDatabasePropertiesRequest{
					"ism_enabled":   true,
					"custom_domain": "custom.upcloud.com",
					"saml":          map[string]interface{}{"enabled": true},
					"openid":        map[string]interface{}{"client_id": "test_client_id"},
				},
			},
		},
		{
			name: "mysql default database type",
			args: []string{
				"--title=mysql-test",
				"--zone=fi-hel1",
				"--host-name-prefix=mysqldb",
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
			name: "missing required title parameter",
			args: []string{
				"--zone=fi-hel1",
				"--host-name-prefix=testdb",
			},
			error: "required flag(s) \"title\" not set",
		},
		{
			name: "missing required zone parameter",
			args: []string{
				"--title=db-test",
				"--host-name-prefix=testdb",
			},
			error: "required flag(s) \"zone\" not set",
		},
		{
			name: "missing required host-name-prefix parameter",
			args: []string{
				"--title=db-test",
				"--zone=fi-hel1",
			},
			error: "required flag(s) \"host-name-prefix\" not set",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			testCmd := CreateCommand()
			mService := new(smock.Service)

			storage.CachedStorages = nil
			createDatabaseReq := test.createDatabaseReq
			mService.On("CreateManagedDatabase", &createDatabaseReq).Return(&databaseDetailsMaint, nil)

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
