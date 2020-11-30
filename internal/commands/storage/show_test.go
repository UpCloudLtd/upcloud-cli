package storage

import (
	"bytes"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestStorageHumanOutput(t *testing.T) {
	s := upcloud.Storage{
		Access:     "private",
		License:    0,
		PartOfPlan: "",
		Size:       10,
		State:      "online",
		Tier:       "maxiops",
		Title:      "test-1",
		Type:       "normal",
		UUID:       "01101f27-196f-47e9-a055-4e2e8bb3b419",
		Zone:       "fi-hel1",
		Origin:     "",
		Created:    time.Time{},
	}

	sd := &upcloud.StorageDetails{
		Storage: s,
		BackupRule: &upcloud.BackupRule{
			Interval:  "daily",
			Time:      "0400",
			Retention: 7,
		},
		ServerUUIDs: []string{"0077fa3d-32db-4b09-9f5f-30d9e9afb565"},
		BackupUUIDs: []string{"012580a1-32a1-466e-a323-689ca16f2d43"},
	}

	sid := &upcloud.StorageImportDetails{
		ClientContentLength: 1,
		ClientContentType:   "abc",
		ErrorCode:           "ghi",
		ErrorMessage:        "jkl",
		MD5Sum:              "mno",
		ReadBytes:           2,
		SHA256Sum:           "pqr",
		Source:              "directUpload",
		State:               "prepared",
		UUID:                "07a6c9a3-300e-4d0e-b935-624f3dbdff3f",
		WrittenBytes:        3,
	}

	storages := []upcloud.Storage{
		{
			PartOfPlan: "yes",
			UUID:       "012580a1-32a1-466e-a323-689ca16f2d43",
			Size:       20,
			Title:      "Storage for server1.example.com",
			Type:       "disk",
			Created:    time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC),
		},
	}

	servers := []upcloud.Server{
		{
			CoreNumber:   0,
			Hostname:     "server1.example.com",
			License:      0,
			MemoryAmount: 2048,
			State:        "started",
			Plan:         "1xCPU-2GB",
			Title:        "server1.example.com",
			UUID:         "0077fa3d-32db-4b09-9f5f-30d9e9afb565",
			Zone:         "fi-hel1",
			Tags: []string{
				"DEV",
				"Ubuntu",
			},
		},
	}

	expected := `  
  Common:
    UUID:       01101f27-196f-47e9-a055-4e2e8bb3b419 
    Title:      test-1                               
    Zone:       fi-hel1                              
    State:      online                               
    Size (GiB): 10                                   
    Type:       normal                               
    Tier:       maxiops                              
    Licence:    0                                    
    Created:                                         
    Origin:                                          
  
  Servers:
    
     UUID                                   Title                 Hostname              State   
    ────────────────────────────────────── ───────────────────── ───────────────────── ─────────
     0077fa3d-32db-4b09-9f5f-30d9e9afb565   server1.example.com   server1.example.com   started 
  
  Backup Rule:
    Interval:  daily 
    Time:      0400  
    Retention: 7     
  
  Available Backups:
    
     UUID                                   Title                             Created              
    ────────────────────────────────────── ───────────────────────────────── ──────────────────────
     012580a1-32a1-466e-a323-689ca16f2d43   Storage for server1.example.com   2020-01-01T00:00:00Z 
  
  Import:
    State:           prepared     
    Source:          directUpload 
    Content Length:  1B           
    Read:            2B           
    Written:         3B           
    SHA256 Checksum: pqr          
    Error:           ghi          
                     jkl          
    Content Type:    abc          
    Created:                      
    Completed:                    
`

	buf := new(bytes.Buffer)
	err := ShowCommand(nil, nil).HandleOutput(buf, &commandResponseHolder{sd, sid, servers, storages})

	assert.Nil(t, err)
	assert.Equal(t, expected, buf.String())
}
