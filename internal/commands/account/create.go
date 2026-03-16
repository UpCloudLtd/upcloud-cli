package account

import (
	"fmt"
	"strings"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CreateCommand creates the "account create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a sub-account",
			"upctl account create --username newuser --password superSecret123",
			`upctl account create --username newuser --password superSecret123 \
	--first-name John --last-name Doe \
	--allow-gui \
	--ip-filter 1.2.3.4 \
	--permission server:00000000-0000-0000-0000-000000000001`,
		),
	}
}

type createCommand struct {
	*commands.BaseCommand
	username    string
	password    string
	firstName   string
	lastName    string
	phone       string
	email       string
	timezone    string
	currency    string
	allowAPI    bool
	allowGUI    bool
	ipFilters   []string
	permissions []string
}

// InitCommand implements Command.InitCommand
func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}

	fs.StringVar(&s.username, "username", "", "Sub-account username.")
	fs.StringVar(&s.password, "password", "", "Sub-account password. Minimum 12 characters with 1 lowercase, 1 uppercase, and 1 number.")
	fs.StringVar(&s.firstName, "first-name", "", "Sub-account first name.")
	fs.StringVar(&s.lastName, "last-name", "", "Sub-account last name.")
	fs.StringVar(&s.phone, "phone", "", "Phone number in international format (e.g. +358.91234567). Defaults to the main account value.")
	fs.StringVar(&s.email, "email", "", "Email address. Defaults to the main account value.")
	fs.StringVar(&s.timezone, "timezone", "", "Timezone. Defaults to the main account value.")
	fs.StringVar(&s.currency, "currency", "", "EUR/GBP/USD/SGD are the only accepted values. Defaults to the main account value.")
	fs.BoolVar(&s.allowAPI, "allow-api", true, "Allow API access for the sub-account.")
	fs.BoolVar(&s.allowGUI, "allow-gui", false, "Allow GUI (control panel) access for the sub-account.")
	fs.StringArrayVar(&s.ipFilters, "ip-filter", []string{}, "Restrict API/GUI access to this IP address. Can be specified multiple times.\n"+
		"Example: --ip-filter 1.2.3.4\n\n--ip-filter 5.6.7.8")
	fs.StringArrayVar(&s.permissions, "permission", []string{}, "Grant a permission to the sub-account in 'target_type:target_identifier' format. Can be specified multiple times.\n"+
		"Example: --permission server:00000000-0000-0000-0000-000000000001")

	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("username"))
	commands.Must(s.Cobra().MarkFlagRequired("password"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("username", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("password", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("first-name", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("last-name", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("phone", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("email", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("timezone", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("currency", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("ip-filter", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("permission", cobra.NoFileCompletions))
}

type parsedPermission struct {
	targetType       string
	targetIdentifier string
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Validate permission flags
	parsedPerms := make([]parsedPermission, 0, len(s.permissions))
	for _, perm := range s.permissions {
		parts := strings.SplitN(perm, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid permission format %q: expected 'target_type:target_identifier'", perm)
		}
		parsedPerms = append(parsedPerms, parsedPermission{targetType: parts[0], targetIdentifier: parts[1]})
	}

	// Validate currency value
	if s.currency != "" && s.currency != "EUR" && s.currency != "GBP" && s.currency != "USD" && s.currency != "SGD" {
		return nil, fmt.Errorf("invalid currency %q: only 'EUR', 'GBP', 'USD', and 'SGD' are accepted", s.currency)
	}

	// Fetch parent account details to use as defaults for unset fields
	parentAccount, err := exec.Account().GetAccount(exec.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to get parent account: %w", err)
	}

	parentDetails, err := exec.Account().GetAccountDetails(exec.Context(), &request.GetAccountDetailsRequest{Username: parentAccount.UserName})
	if err != nil {
		return nil, fmt.Errorf("failed to get parent account details: %w", err)
	}

	// Apply parent account values as defaults for fields not explicitly set
	phone := s.phone
	if phone == "" {
		phone = parentDetails.Phone
	}
	email := s.email
	if email == "" {
		email = parentDetails.Email
	}
	timezone := s.timezone
	if timezone == "" {
		timezone = parentDetails.Timezone
	}

	currency := s.currency
	if currency == "" {
		currency = parentDetails.Currency
	}

	allowAPI := upcloud.FromBool(s.allowAPI)
	allowGUI := upcloud.FromBool(s.allowGUI)

	// Build the IP filter list
	ipFilters := upcloud.AccountIPFilters{IPFilter: []string{}}
	for _, f := range s.ipFilters {
		ipFilters.IPFilter = append(ipFilters.IPFilter, f)
	}

	msg := fmt.Sprintf("Creating sub-account %s", s.username)
	exec.PushProgressStarted(msg)

	// Create the sub-account. The access control arrays (roles, network_access,
	// server_access, storage_access, tag_access) are sent as empty slices because
	// the API requires those fields to be present. Permissions are granted
	// separately below via the permissions endpoint.
	_, err = exec.All().CreateSubaccount(exec.Context(), &request.CreateSubaccountRequest{
		Subaccount: request.CreateSubaccount{
			Username:      s.username,
			Password:      s.password,
			FirstName:     s.firstName,
			LastName:      s.lastName,
			Phone:         phone,
			Email:         email,
			Timezone:      timezone,
			Language:      "en",
			Currency:      currency,
			AllowAPI:      allowAPI,
			AllowGUI:      allowGUI,
			Roles:         upcloud.AccountRoles{Role: []string{}},
			NetworkAccess: upcloud.AccountNetworkAccess{Network: []string{}},
			ServerAccess:  upcloud.AccountServerAccess{Server: []upcloud.AccountServer{}},
			StorageAccess: upcloud.AccountStorageAccess{Storage: []string{}},
			TagAccess:     upcloud.AccountTagAccess{Tag: []upcloud.AccountTag{}},

			IPFilters: ipFilters,
		},
	})
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	// Grant any permissions that were requested via --permission flags
	for _, perm := range parsedPerms {
		permMsg := fmt.Sprintf("Granting %s:%s permission to %s", perm.targetType, perm.targetIdentifier, s.username)
		exec.PushProgressStarted(permMsg)

		_, err = exec.All().GrantPermission(exec.Context(), &request.GrantPermissionRequest{
			Permission: upcloud.Permission{
				User:             s.username,
				TargetType:       upcloud.PermissionTarget(perm.targetType),
				TargetIdentifier: perm.targetIdentifier,
			},
		})
		if err != nil {
			return commands.HandleError(exec, permMsg, err)
		}

		exec.PushProgressSuccess(permMsg)
	}

	return output.None{}, nil
}
