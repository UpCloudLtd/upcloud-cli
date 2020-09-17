package cmd

import (
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const version = "0.1"
const envPrefix = "UPCLOUD"

type userAgentTransport struct{}
func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", fmt.Sprintf("upctl/%s", version))
	return cleanhttp.DefaultTransport().RoundTrip(req)
}

var apiService *service.Service
var cfgFile string
var verbose, jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "up",
	Short: "UpCloud command line client",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of up, notify if a newer version is available.",
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Printf("UpCloud command client version: %s\n", version)
			// TODO check this by comparing version against Github releases.
			// https://api.github.com/repos/UpCloudLtd/cli/releases
			fmt.Println("Software is up to date.")
		} else {
			fmt.Println(version)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.upcloud-api.yml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Return output in JSON")
}

func initConfig() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}

	viper.SetConfigName(".upcloud-api")
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		configDir := path.Dir(cfgFile)
		if configDir != "." && configDir != dir {
			viper.AddConfigPath(configDir)
		}
	}

	viper.AddConfigPath(dir)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()


	viper.ReadInConfig()

	missingConfigDefault := "%s not defined. Either define it as '%s' in the config file ($HOME/.upcloud-api.yml), or as %s env variable.\n"

	username := viper.GetString("USERNAME")
	if username == "" {
		fmt.Printf(missingConfigDefault, "Username", "username", fmt.Sprintf("%s_%s", envPrefix, "USERNAME"))
		os.Exit(1)
	}

	password := viper.GetString("PASSWORD")
	if password == "" {
		fmt.Printf(missingConfigDefault, "Password", "password", fmt.Sprintf("%s_%s", envPrefix, "PASSWORD"))
		os.Exit(1)
	}

	apiService = service.New(client.NewWithHTTPClient(username, password, &http.Client{Transport: &userAgentTransport{}}))
}