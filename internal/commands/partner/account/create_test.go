package partneraccount

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommand(t *testing.T) {
	targetMethod := "CreatePartnerAccount"
	password := "superSecret123"

	generateArgs := func(namePhoneEmail bool, country string, state string, fullDetails bool) (args []string) {
		args = append(args,
			"--username", testPartnerAccount.Username,
			"--password", password,
		)
		if namePhoneEmail {
			args = append(args,
				"--first-name", testPartnerAccount.FirstName,
				"--last-name", testPartnerAccount.LastName,
				"--phone", testPartnerAccount.Phone,
				"--email", testPartnerAccount.Email,
			)
		}
		if country != "" {
			args = append(args, "--country", country)
		}
		if state != "" {
			args = append(args, "--state", state)
		} else if fullDetails {
			args = append(args,
				"--state", testPartnerAccount.State,
				"--company", testPartnerAccount.Company,
				"--address", testPartnerAccount.Address,
				"--postal-code", testPartnerAccount.PostalCode,
				"--city", testPartnerAccount.City,
				"--vat-number", testPartnerAccount.VATNumber,
			)
		}
		return
	}

	generateRequest := func(country string, state string, fullDetails bool) (req request.CreatePartnerAccountRequest) {
		req.Username = testPartnerAccount.Username
		req.Password = password
		if country != "" {
			cd := &request.CreatePartnerAccountContactDetails{}
			cd.FirstName = testPartnerAccount.FirstName
			cd.LastName = testPartnerAccount.LastName
			cd.Country = country
			cd.Phone = testPartnerAccount.Phone
			cd.Email = testPartnerAccount.Email
			if state != "" {
				cd.State = state
			} else if fullDetails {
				cd.State = testPartnerAccount.State
				cd.Company = testPartnerAccount.Company
				cd.Address = testPartnerAccount.Address
				cd.PostalCode = testPartnerAccount.PostalCode
				cd.City = testPartnerAccount.City
				cd.VATNumber = testPartnerAccount.VATNumber
			}
			req.ContactDetails = cd
		}
		return
	}

	for _, test := range []struct {
		name    string
		args    []string
		request request.CreatePartnerAccountRequest
		error   string
	}{
		{
			name:  "no args",
			args:  []string{},
			error: `required flag(s) "username", "password" not set`,
		},
		{
			name:    "no contact details",
			args:    generateArgs(false, "", "", false),
			request: generateRequest("", "", false),
		},
		{
			name:  "partial contact details",
			args:  generateArgs(true, "", "", false),
			error: `when contact details are given, the following flags are required: "first-name", "last-name", "country", "phone", "email"`,
		},
		{
			name:    "minimal contact details",
			args:    generateArgs(true, testPartnerAccount.Country, "", false),
			request: generateRequest(testPartnerAccount.Country, "", false),
		},
		{
			name:  "USA without a state",
			args:  generateArgs(true, "USA", "", false),
			error: `when contact country is "USA", flag "state" is also required`,
		},
		{
			name:    "USA with a state",
			args:    generateArgs(true, "USA", "Florida", false),
			request: generateRequest("USA", "Florida", false),
		},
		{
			name:    "full contact details",
			args:    generateArgs(true, testPartnerAccount.Country, "", true),
			request: generateRequest(testPartnerAccount.Country, "", true),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			mService := smock.Service{}
			req := test.request
			mService.On(targetMethod, &req).Return(&testPartnerAccount, nil)

			conf := config.New()
			command := commands.BuildCommand(CreateCommand(), nil, conf)
			command.Cobra().SetArgs(test.args)
			_, err := mockexecute.MockExecute(command, &mService, conf)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
			} else {
				assert.NoError(t, err)
				mService.AssertNumberOfCalls(t, targetMethod, 1)
			}
		})
	}
}
