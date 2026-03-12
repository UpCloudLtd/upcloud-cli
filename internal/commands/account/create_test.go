package account

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// parentAccount and parentDetails are the values returned by the mock for the
// parent account, used to verify default-inheritance behaviour.
var (
	parentAccount = upcloud.Account{
		UserName: "parent_user",
	}
	parentDetails = upcloud.AccountDetails{
		Phone:    "+358.91234567",
		Email:    "parent@example.com",
		Timezone: "Europe/Helsinki",
		Currency: "EUR",
	}
)

func setupCreateMocks(t *testing.T, mService *smock.Service, req *request.CreateSubaccountRequest) {
	t.Helper()
	mService.On("GetAccount").Return(&parentAccount, nil)
	mService.On("GetAccountDetails", &request.GetAccountDetailsRequest{Username: parentAccount.UserName}).
		Return(&parentDetails, nil)
	if req != nil {
		mService.On("CreateSubaccount", req).Return(&upcloud.AccountDetails{}, nil)
	}
}

func TestCreateCommand(t *testing.T) {
	emptyAccess := request.CreateSubaccount{
		Roles:         upcloud.AccountRoles{Role: []string{}},
		NetworkAccess: upcloud.AccountNetworkAccess{Network: []string{}},
		ServerAccess:  upcloud.AccountServerAccess{Server: []upcloud.AccountServer{}},
		StorageAccess: upcloud.AccountStorageAccess{Storage: []string{}},
		TagAccess:     upcloud.AccountTagAccess{Tag: []upcloud.AccountTag{}},
	}

	for _, test := range []struct {
		name        string
		args        []string
		req         *request.CreateSubaccountRequest
		permissions []request.GrantPermissionRequest
		errFn       assert.ErrorAssertionFunc
	}{
		{
			name: "minimal — required flags only, defaults inherited from parent",
			args: []string{
				"--username", "test_sub",
				"--password", "MyP@ssw0rd123",
			},
			req: &request.CreateSubaccountRequest{
				Subaccount: func() request.CreateSubaccount {
					s := emptyAccess
					s.Username = "test_sub"
					s.Password = "MyP@ssw0rd123"
					s.Phone = parentDetails.Phone
					s.Email = parentDetails.Email
					s.Timezone = parentDetails.Timezone
					s.Currency = parentDetails.Currency
					s.Language = "en"
					s.AllowAPI = upcloud.True
					s.AllowGUI = upcloud.False
					s.IPFilters = upcloud.AccountIPFilters{IPFilter: []string{}}
					return s
				}(),
			},
			errFn: assert.NoError,
		},
		{
			name: "all personal details provided",
			args: []string{
				"--username", "test_sub",
				"--password", "MyP@ssw0rd123",
				"--first-name", "John",
				"--last-name", "Doe",
				"--phone", "+1.5551234567",
				"--email", "john@example.com",
				"--timezone", "America/New_York",
			},
			req: &request.CreateSubaccountRequest{
				Subaccount: func() request.CreateSubaccount {
					s := emptyAccess
					s.Username = "test_sub"
					s.Password = "MyP@ssw0rd123"
					s.FirstName = "John"
					s.LastName = "Doe"
					s.Phone = "+1.5551234567"
					s.Email = "john@example.com"
					s.Timezone = "America/New_York"
					s.Currency = parentDetails.Currency
					s.Language = "en"
					s.AllowAPI = upcloud.True
					s.AllowGUI = upcloud.False
					s.IPFilters = upcloud.AccountIPFilters{IPFilter: []string{}}
					return s
				}(),
			},
			errFn: assert.NoError,
		},
		{
			name: "currency overridden",
			args: []string{
				"--username", "test_sub",
				"--password", "MyP@ssw0rd123",
				"--currency", "USD",
			},
			req: &request.CreateSubaccountRequest{
				Subaccount: func() request.CreateSubaccount {
					s := emptyAccess
					s.Username = "test_sub"
					s.Password = "MyP@ssw0rd123"
					s.Phone = parentDetails.Phone
					s.Email = parentDetails.Email
					s.Timezone = parentDetails.Timezone
					s.Currency = "USD"
					s.Language = "en"
					s.AllowAPI = upcloud.True
					s.AllowGUI = upcloud.False
					s.IPFilters = upcloud.AccountIPFilters{IPFilter: []string{}}
					return s
				}(),
			},
			errFn: assert.NoError,
		},
		{
			name: "allow-gui enabled, allow-api disabled",
			args: []string{
				"--username", "test_sub",
				"--password", "MyP@ssw0rd123",
				"--allow-gui",
				"--allow-api=false",
			},
			req: &request.CreateSubaccountRequest{
				Subaccount: func() request.CreateSubaccount {
					s := emptyAccess
					s.Username = "test_sub"
					s.Password = "MyP@ssw0rd123"
					s.Phone = parentDetails.Phone
					s.Email = parentDetails.Email
					s.Timezone = parentDetails.Timezone
					s.Currency = parentDetails.Currency
					s.Language = "en"
					s.AllowAPI = upcloud.False
					s.AllowGUI = upcloud.True
					s.IPFilters = upcloud.AccountIPFilters{IPFilter: []string{}}
					return s
				}(),
			},
			errFn: assert.NoError,
		},
		{
			name: "ip filters provided",
			args: []string{
				"--username", "test_sub",
				"--password", "MyP@ssw0rd123",
				"--ip-filter", "1.2.3.4",
				"--ip-filter", "5.6.7.8",
			},
			req: &request.CreateSubaccountRequest{
				Subaccount: func() request.CreateSubaccount {
					s := emptyAccess
					s.Username = "test_sub"
					s.Password = "MyP@ssw0rd123"
					s.Phone = parentDetails.Phone
					s.Email = parentDetails.Email
					s.Timezone = parentDetails.Timezone
					s.Currency = parentDetails.Currency
					s.Language = "en"
					s.AllowAPI = upcloud.True
					s.AllowGUI = upcloud.False
					s.IPFilters = upcloud.AccountIPFilters{IPFilter: []string{"1.2.3.4", "5.6.7.8"}}
					return s
				}(),
			},
			errFn: assert.NoError,
		},
		{
			name: "invalid currency rejected before API call",
			args: []string{
				"--username", "test_sub",
				"--password", "MyP@ssw0rd123",
				"--currency", "BTC",
			},
			req: nil, // CreateSubaccount must NOT be called
			errFn: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, `invalid currency "BTC"`)
			},
		},
		{
			name: "missing --username",
			args: []string{
				"--password", "MyP@ssw0rd123",
			},
			req: nil,
			errFn: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, `required flag(s) "username" not set`)
			},
		},
		{
			name: "missing --password",
			args: []string{
				"--username", "test_sub",
			},
			req: nil,
			errFn: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, `required flag(s) "password" not set`)
			},
		},
		{
			name: "single permission granted",
			args: []string{
				"--username", "test_sub",
				"--password", "MyP@ssw0rd123",
				"--permission", "server:*",
			},
			req: &request.CreateSubaccountRequest{
				Subaccount: func() request.CreateSubaccount {
					s := emptyAccess
					s.Username = "test_sub"
					s.Password = "MyP@ssw0rd123"
					s.Phone = parentDetails.Phone
					s.Email = parentDetails.Email
					s.Timezone = parentDetails.Timezone
					s.Currency = parentDetails.Currency
					s.Language = "en"
					s.AllowAPI = upcloud.True
					s.AllowGUI = upcloud.False
					s.IPFilters = upcloud.AccountIPFilters{IPFilter: []string{}}
					return s
				}(),
			},
			permissions: []request.GrantPermissionRequest{
				{
					Permission: upcloud.Permission{
						User:             "test_sub",
						TargetType:       upcloud.PermissionTarget("server"),
						TargetIdentifier: "*",
					},
				},
			},
			errFn: assert.NoError,
		},
		{
			name: "multiple permissions granted",
			args: []string{
				"--username", "test_sub",
				"--password", "MyP@ssw0rd123",
				"--permission", "server:*",
				"--permission", "storage:*",
				"--permission", "network:*",
			},
			req: &request.CreateSubaccountRequest{
				Subaccount: func() request.CreateSubaccount {
					s := emptyAccess
					s.Username = "test_sub"
					s.Password = "MyP@ssw0rd123"
					s.Phone = parentDetails.Phone
					s.Email = parentDetails.Email
					s.Timezone = parentDetails.Timezone
					s.Currency = parentDetails.Currency
					s.Language = "en"
					s.AllowAPI = upcloud.True
					s.AllowGUI = upcloud.False
					s.IPFilters = upcloud.AccountIPFilters{IPFilter: []string{}}
					return s
				}(),
			},
			permissions: []request.GrantPermissionRequest{
				{Permission: upcloud.Permission{User: "test_sub", TargetType: upcloud.PermissionTarget("server"), TargetIdentifier: "*"}},
				{Permission: upcloud.Permission{User: "test_sub", TargetType: upcloud.PermissionTarget("storage"), TargetIdentifier: "*"}},
				{Permission: upcloud.Permission{User: "test_sub", TargetType: upcloud.PermissionTarget("network"), TargetIdentifier: "*"}},
			},
			errFn: assert.NoError,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			conf := config.New()
			mService := new(smock.Service)

			setupCreateMocks(t, mService, test.req)
			if len(test.permissions) == 0 {
				mService.On("GrantPermission", mock.Anything).Return(&upcloud.Permission{}, nil).Maybe()
			} else {
				for i := range test.permissions {
					mService.On("GrantPermission", &test.permissions[i]).Return(&upcloud.Permission{}, nil).Once()
				}
			}

			c := commands.BuildCommand(CreateCommand(), nil, conf)
			c.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(c, mService, conf)

			test.errFn(t, err)
			if test.req != nil {
				mService.AssertCalled(t, "CreateSubaccount", test.req)
			} else {
				mService.AssertNotCalled(t, "CreateSubaccount", mock.Anything)
			}
			for i := range test.permissions {
				mService.AssertCalled(t, "GrantPermission", &test.permissions[i])
			}
		})
	}
}
