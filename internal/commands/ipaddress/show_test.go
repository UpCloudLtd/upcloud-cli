package ipaddress

import (
	"bytes"
	"github.com/UpCloudLtd/cli/internal/commands"
	"github.com/UpCloudLtd/cli/internal/config"
	smock "github.com/UpCloudLtd/cli/internal/mock"
	"github.com/UpCloudLtd/cli/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"

	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestShowCommand(t *testing.T) {
	ipAddress := upcloud.IPAddress{
		Address:    "94.237.117.150",
		Access:     "public",
		Family:     "IPv4",
		PartOfPlan: upcloud.FromBool(true),
		PTRRecord:  "94-237-117-150.fi-hel1.upcloud.host",
		ServerUUID: "005ab220-7ff6-42c9-8615-e4c02eb4104e",
		MAC:        "ee:1b:db:ca:6b:80",
		Floating:   upcloud.FromBool(false),
		Zone:       "fi-hel1",
	}

	expected := `  
  Address:      94.237.117.150                       
  Access:       public                               
  Family:       IPv4                                 
  Part of Plan: yes                                  
  PTR Record:   94-237-117-150.fi-hel1.upcloud.host  
  Server UUID:  005ab220-7ff6-42c9-8615-e4c02eb4104e 
  MAC:          ee:1b:db:ca:6b:80                    
  Floating:     no                                   
  Zone:         fi-hel1                              

`

	svc := &smock.Service{}
	conf := config.New()

	svc.On("GetIPAddressDetails",
		&request.GetIPAddressDetailsRequest{Address: ipAddress.Address},
	).Return(&ipAddress, nil)
	conf.Viper().Set(config.KeyOutput, config.ValueOutputHuman)

	command := commands.BuildCommand(ShowCommand(), nil, conf)
	out, err := command.(commands.Command).Execute(commands.NewExecutor(conf, svc), ipAddress.Address)
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf, out)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())

	svc.AssertNumberOfCalls(t, "GetIPAddressDetails", 1)

}
