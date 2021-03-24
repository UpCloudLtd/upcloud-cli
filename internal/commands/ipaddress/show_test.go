package ipaddress

import (
	"bytes"
	"testing"

	smock "github.com/UpCloudLtd/cli/internal/mock"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestShowCommand(t *testing.T) {

	account := upcloud.IPAddress{
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

	buf := new(bytes.Buffer)
	command := ShowCommand(&smock.MockService{})
	err := command.HandleOutput(buf, &account)

	assert.Nil(t, err)
	assert.Equal(t, expected, buf.String())
}
