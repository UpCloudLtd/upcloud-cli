package partneraccount

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type createCommand struct {
	*commands.BaseCommand
	params request.CreatePartnerAccountRequest
}

func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a new account that will be linked to partner's existing invoicing",
			"upctl partner account create --username newuser --password superSecret123",
			`upctl partner account create --username newuser --password superSecret123 --first-name New --last-name User --company "Example Ltd" --country FIN --phone +358.91111111 --email new.user@gmail.com`,
		),
	}
}

func (s *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	s.params.ContactDetails = &request.CreatePartnerAccountContactDetails{}
	cReqDesc := " Required when other contact details are given."

	fs.StringVar(&s.params.Username, "username", "", "Account username.")
	fs.StringVar(&s.params.Password, "password", "", "Account password.")
	fs.StringVar(&s.params.ContactDetails.FirstName, "first-name", "", "Contact first name."+cReqDesc)
	fs.StringVar(&s.params.ContactDetails.LastName, "last-name", "", "Contact last name."+cReqDesc)
	fs.StringVar(&s.params.ContactDetails.Company, "company", "", "Contact company name.")
	fs.StringVar(&s.params.ContactDetails.Address, "address", "", "Contact street address.")
	fs.StringVar(&s.params.ContactDetails.PostalCode, "postal-code", "", "Contact postal/zip code.")
	fs.StringVar(&s.params.ContactDetails.City, "city", "", "Contact city.")
	fs.StringVar(&s.params.ContactDetails.State, "state", "", "Contact state. Required when other contact details are given and country is 'USA'.")
	fs.StringVar(&s.params.ContactDetails.Country, "country", "", "Contact ISO 3166-1 three character country code."+cReqDesc)
	fs.StringVar(&s.params.ContactDetails.Phone, "phone", "", "Contact phone number in international format, country code and national part separated by a period."+cReqDesc)
	fs.StringVar(&s.params.ContactDetails.Email, "email", "", "Contact email address."+cReqDesc)
	fs.StringVar(&s.params.ContactDetails.VATNumber, "vat-number", "", "Contact VAT number.")

	s.AddFlags(fs)
	commands.Must(s.Cobra().MarkFlagRequired("username"))
	commands.Must(s.Cobra().MarkFlagRequired("password"))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("username", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("password", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("first-name", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("last-name", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("company", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("address", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("postal-code", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("city", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("state", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("country", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("phone", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("email", cobra.NoFileCompletions))
	commands.Must(s.Cobra().RegisterFlagCompletionFunc("vat-number", cobra.NoFileCompletions))
}

func (s *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	if (*s.params.ContactDetails == request.CreatePartnerAccountContactDetails{}) {
		s.params.ContactDetails = nil
	} else {
		cd := s.params.ContactDetails
		if cd.FirstName == "" || cd.LastName == "" || cd.Country == "" || cd.Phone == "" || cd.Email == "" {
			return nil, fmt.Errorf(`when contact details are given, the following flags are required: "first-name", "last-name", "country", "phone", "email"`)
		}
		if cd.Country == "USA" && cd.State == "" {
			return nil, fmt.Errorf(`when contact country is "USA", flag "state" is also required`)
		}
	}

	msg := fmt.Sprintf("Creating account %s", s.params.Username)
	exec.PushProgressStarted(msg)

	_, err := exec.All().CreatePartnerAccount(exec.Context(), &s.params)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.None{}, nil
}
