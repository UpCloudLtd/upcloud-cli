package partneraccount

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/testutils"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

var testPartnerAccount = upcloud.PartnerAccount{
	Username:   "someuser",
	FirstName:  "Some",
	LastName:   "User",
	Country:    "FIN",
	State:      "Some state",
	Email:      "some.user@somecompany.com",
	Phone:      "+358.91234567",
	Company:    "Some company",
	Address:    "Some street",
	PostalCode: "00100",
	City:       "Some city",
	VATNumber:  "FI12345678",
}

func TestListCommand(t *testing.T) {
	mService := smock.Service{}
	mService.On("GetPartnerAccounts").Return([]upcloud.PartnerAccount{testPartnerAccount}, nil)

	conf := config.New()
	conf.Viper().Set(config.KeyOutput, config.ValueOutputJSON)

	command := commands.BuildCommand(ListCommand(), nil, conf)

	output, err := mockexecute.MockExecute(command, &mService, conf)
	assert.NoError(t, err)
	testutils.AssertOutputIsList(t, output)
}
