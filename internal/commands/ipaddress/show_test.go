package ipaddress

import (
	"bytes"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/gemalto/flume"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
)

func TestShowCommand(t *testing.T) {
	text.DisableColors()
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
	out, err := command.(commands.MultipleArgumentCommand).Execute(commands.NewExecutor(conf, svc, flume.New("test")), ipAddress.Address)
	assert.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = output.Render(buf, conf.Output(), out)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())

	svc.AssertNumberOfCalls(t, "GetIPAddressDetails", 1)
}
