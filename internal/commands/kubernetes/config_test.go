package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/mockexecute"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestConfigCommand(t *testing.T) {
	text.DisableColors()
	dir := t.TempDir()
	for _, tt := range []struct {
		name                 string
		output               string
		args                 []string
		wantErr              bool
		existingFileContents []byte
		expectedOutput       string
		expectedFileContents []byte
	}{
		{
			name:   "human-output",
			output: config.ValueOutputHuman,
			args: []string{
				exampleUUID(),
			},
			expectedOutput: `  
  Overview:
    Current context: human-output-admin@human-output 

  Clusters

     Server                 Insecure skip TLS verify 
    ────────────────────── ──────────────────────────
     https://human-output   false                    
    
  Contexts

     Cluster        Authinfo           
    ────────────── ────────────────────
     human-output   human-output-admin 
    
`,
		},
		{
			name:   "json-output",
			output: config.ValueOutputJSON,
			args: []string{
				exampleUUID(),
			},
			expectedOutput: string(yamlToJSON(exampleKubernetesKubeconfig("json-output"))),
		},
		{
			name:   "yaml-output",
			output: config.ValueOutputYAML,
			args: []string{
				exampleUUID(),
			},
			expectedOutput: `apiVersion: v1
clusters:
    - cluster:
        certificate-authority: RkFLRQ==
        certificate-authority-data: RkFLRQ==
        server: https://yaml-output
      name: yaml-output
contexts:
    - context:
        cluster: yaml-output
        user: yaml-output-admin
      name: yaml-output
current-context: yaml-output-admin@yaml-output
kind: Config
preferences: {}
users:
    - name: yaml-output
      user:
        client-certificate: RkFLRQ==
        client-certificate-data: RkFLRQ==
        client-key: RkFLRQ==
        client-key-data: RkFLRQ==
`,
		},
		{
			name:   "write-fail",
			output: config.ValueOutputYAML,
			args: []string{
				exampleUUID(),
				"--write",
				"",
			},
			expectedOutput: `Error: invalid write path
Usage:
  config <UUID/Name...> [flags]

Examples:
upctl kubernetes config 0fa980c4-0e4f-460b-9869-11b7bd62b831 --output human
upctl kubernetes config 0fa980c4-0e4f-460b-9869-11b7bd62b831 --output yaml --write $KUBECONFIG
upctl kubernetes config 0fa980c4-0e4f-460b-9869-11b7bd62b831 --output yaml --write ./my_kubeconfig.yaml

Flags:
      --write string   Absolute path for writing output. If the file exists, the config will be merged.
  -h, --help           help for config

`,
			wantErr: true,
		},
		{
			name:   "write-to-empty-file",
			output: config.ValueOutputYAML,
			args: []string{
				exampleUUID(),
				"--write",
				exampleFilename(dir, "write-to-empty-file"),
			},
			expectedOutput:       ``,
			expectedFileContents: exampleKubernetesKubeconfig("write-to-empty-file"),
		},
		{
			name:   "write-to-non-empty-file",
			output: config.ValueOutputYAML,
			args: []string{
				exampleUUID(),
				"--write",
				exampleFilename(dir, "write-to-non-empty-file"),
			},
			existingFileContents: exampleKubernetesKubeconfig("previous-config"),
			expectedOutput:       ``,
			expectedFileContents: exampleKubernetesKubeconfig("previous-config", "write-to-non-empty-file"),
		},
		{
			name:   "write-to-non-empty-file-with-override",
			output: config.ValueOutputYAML,
			args: []string{
				exampleUUID(),
				"--write",
				exampleFilename(dir, "write-to-non-empty-file-with-override"),
			},
			existingFileContents: exampleKubernetesKubeconfig("write-to-non-empty-file-with-override"),
			expectedOutput:       ``,
			expectedFileContents: exampleKubernetesKubeconfig("write-to-non-empty-file-with-override"),
		},
	} {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			mService := smock.Service{}
			mService.On("GetKubernetesClusters", mock.Anything).
				Return([]upcloud.KubernetesCluster{exampleKubernetesCluster(tt.name)}, nil)

			mService.On("GetKubernetesKubeconfig", exampleGetKubernetesKubeconfigRequest()).
				Return(string(exampleKubernetesKubeconfig(tt.name)), nil)

			filename := exampleFilename(dir, tt.name)
			file, err := os.Create(filename)
			if err != nil {
				t.Fatalf("Create(%s): %v", filename, err)
			}

			if tt.existingFileContents != nil {
				_, err := file.Write(tt.existingFileContents)
				if err != nil {
					t.Fatalf("Write: %v", err)
				}
			}

			err = file.Close()
			if err != nil {
				t.Fatalf("Close: %v", err)
			}

			conf := config.New()
			conf.Viper().Set(config.KeyOutput, tt.output)
			command := commands.BuildCommand(ConfigCommand(), nil, conf)

			// get resolver to initialize command cache
			_, err = command.(*configCommand).Get(context.TODO(), &mService)
			if err != nil {
				t.Fatal(err)
			}

			command.Cobra().SetArgs(tt.args)

			actualOutput, err := mockexecute.MockExecute(command, &mService, conf)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedOutput, actualOutput)

			if tt.expectedFileContents != nil {
				actualFileContents, err := os.ReadFile(filename)
				if err != nil {
					t.Fatalf("ReadFile(%s): %v", filename, err)
				}

				assert.Equal(t, tt.expectedFileContents, actualFileContents)
			}
		})
	}
}

func exampleGetKubernetesKubeconfigRequest() *request.GetKubernetesKubeconfigRequest {
	return &request.GetKubernetesKubeconfigRequest{UUID: exampleUUID()}
}

func exampleKubernetesKubeconfig(names ...string) []byte {
	apiConfig := api.NewConfig()
	for _, v := range names {
		apiConfig.Clusters[v] = &api.Cluster{
			Server:                   fmt.Sprintf("https://%s", v),
			InsecureSkipTLSVerify:    false,
			CertificateAuthority:     "RkFLRQ==",
			CertificateAuthorityData: []byte("FAKE"),
		}
		apiConfig.AuthInfos[v] = &api.AuthInfo{
			ClientCertificate:     "RkFLRQ==",
			ClientCertificateData: []byte("FAKE"),
			ClientKey:             "RkFLRQ==",
			ClientKeyData:         []byte("FAKE"),
		}
		apiConfig.Contexts[v] = &api.Context{
			LocationOfOrigin: "",
			Cluster:          v,
			AuthInfo:         fmt.Sprintf("%s-admin", v),
		}
		apiConfig.CurrentContext = fmt.Sprintf("%s-admin@%s", v, v)
	}

	b, _ := clientcmd.Write(*apiConfig)

	return b
}

func exampleFilename(dir, testName string) string {
	return fmt.Sprintf("%s%s", dir, testName)
}

func exampleKubernetesCluster(name string) upcloud.KubernetesCluster {
	return upcloud.KubernetesCluster{
		Name: name,
		UUID: exampleUUID(),
	}
}

func exampleUUID() string {
	return "0fa980c4-0e4f-460b-9869-11b7bd62b831"
}

func yamlToJSON(yamlIn []byte) []byte {
	if len(yamlIn) == 0 {
		return []byte{}
	}

	var jsonObj interface{}
	_ = yaml.Unmarshal(yamlIn, &jsonObj)

	out, _ := json.MarshalIndent(jsonObj, "", "  ")

	return append(out, []byte("\n")...)
}
