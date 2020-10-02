package old

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var defaultOSTemplates map[string]string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "List, show & control servers",
}

var listServersCmd = &cobra.Command{
	Use:   "list",
	Short: "List current servers",
	RunE: func(cmd *cobra.Command, args []string) error {

		servers, err := apiService.GetServers()
		if err != nil {
			return err
		}

		if jsonOutput {
			serversJSON, err := json.MarshalIndent(servers.Servers, "", "  ")
			if err != nil {
				return errors.Wrap(err, "JSON serialization failed")
			}

			fmt.Printf("%s\n", serversJSON)
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"UUID", "Hostname", "Plan", "Zone", "State"})

		for _, server := range servers.Servers {
			plan := server.Plan
			if plan == "custom" {
				memory := server.MemoryAmount / 1024
				plan = fmt.Sprintf("Custom (%dxCPU, %dGB)", server.CoreNumber, memory)
			}
			table.Append([]string{server.UUID, server.Hostname, plan, server.Zone, server.State})
		}
		table.Render()

		return nil
	},
}

var showServerCmd = &cobra.Command{
	Use:   "show UUID",
	Short: "Show server information",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Must specify a single server")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := apiService.GetServerDetails(&request.GetServerDetailsRequest{UUID: args[0]})
		if err != nil {
			return errors.Wrap(err, "Fetching server information failed")
		}

		if jsonOutput {
			serverJSON, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
				return errors.Wrap(err, "JSON serialization failed")
			}

			fmt.Printf("%s\n", serverJSON)
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"", server.Title})

		table.Append([]string{"Title", server.Title})
		table.Append([]string{"Hostname", server.Hostname})
		table.Append([]string{"UUID", server.UUID})

		plan := server.Plan
		if plan == "custom" {
			memory := server.MemoryAmount / 1024
			plan = fmt.Sprintf("Custom (%dxCPU, %dGB)", server.CoreNumber, memory)
		}
		table.Append([]string{"Plan", plan})
		table.Append([]string{"Zone", server.Zone})
		table.Append([]string{"Tags", strings.Join(server.Tags, ", ")})
		table.Append([]string{"State", server.State})
		table.Append([]string{"Firewall", server.Firewall})
		table.Append([]string{"Metadata", server.Metadata.String()})
		table.Render()

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"", "REMOTE ACCESS"})
		table.Append([]string{"Enabled", server.RemoteAccessEnabled.String()})
		if server.RemoteAccessEnabled.Bool() {
			table.Append([]string{"Type", server.RemoteAccessType})
			table.Append([]string{"Host", server.RemoteAccessHost})
			table.Append([]string{"Port", fmt.Sprintf("%d", server.RemoteAccessPort)})
			table.Append([]string{"Password", server.RemoteAccessPassword})
		}
		table.Render()
		fmt.Println()

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Storage", "UUID", "Size (GB)", "Type"})
		for _, storage := range server.StorageDevices {
			table.Append([]string{storage.Title, storage.UUID, fmt.Sprintf("%d", storage.Size), storage.Type})
		}
		table.Render()
		fmt.Println()

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Interface", "Type", "IP addresses", "Family", "MAC"})
		for _, iface := range server.Networking.Interfaces {
			addresses := ""
			var family string
			for _, addr := range iface.IPAddresses {
				if addresses != "" {
					addresses += ", " + addr.Address
				} else {
					addresses = addr.Address
					family = addr.Family
				}
				if addr.Floating.Bool() {
					addresses += " (floating)"
				}
			}
			table.Append([]string{fmt.Sprintf("%d", iface.Index), iface.Type, addresses, family, iface.MAC})
		}

		table.Render()

		return nil
	},
}

var createServerCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("hostname") {
			return errors.New("Hostname must be specified")
		}

		coreMemoryError := "Both memory and cores must be specified, if other one is"
		if cmd.Flags().Changed("cores") && !cmd.Flags().Changed("memory") {
			return errors.New(coreMemoryError)
		}
		if !cmd.Flags().Changed("cores") && cmd.Flags().Changed("memory") {
			return errors.New(coreMemoryError)
		}

		keys, _ := cmd.Flags().GetString("ssh-keys-file")
		if keys != "" {
			if !fileExists(keys) {
				return errors.New(fmt.Sprintf("SSH keys file %s does not exist", keys))
			}
		}
		userData, _ := cmd.Flags().GetString("userdata-file")
		if userData != "" {
			if !fileExists(userData) {
				return errors.New(fmt.Sprintf("Userdata file %s does not exist", userData))
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		req := request.CreateServerRequest{}

		req.Zone, _ = cmd.Flags().GetString("zone")
		req.PasswordDelivery, _ = cmd.Flags().GetString("password-delivery")

		req.Hostname, _ = cmd.Flags().GetString("hostname")
		if cmd.Flags().Changed("title") {
			req.Title, _ = cmd.Flags().GetString("title")
		} else {
			req.Title = req.Hostname
		}

		storageSize, _ := cmd.Flags().GetInt("storage-size")
		storageTier, _ := cmd.Flags().GetString("storage-tier")
		storageAction := "clone"
		if cmd.Flags().Changed("cores") && cmd.Flags().Changed("memory") {
			req.CoreNumber, _ = cmd.Flags().GetInt("cores")
			req.MemoryAmount, _ = cmd.Flags().GetInt("memory")
		} else {
			req.Plan, _ = cmd.Flags().GetString("plan")
			plans, err := apiService.GetPlans()
			if err == nil {
				for _, plan := range plans.Plans {
					if strings.EqualFold(req.Plan, plan.Name) {
						storageSize = maxInt(storageSize, plan.StorageSize)
						storageTier = plan.StorageTier
						break
					}
				}
			}
		}

		storageID, _ := cmd.Flags().GetString("storage")
		if val, ok := defaultOSTemplates[storageID]; ok {
			storageID = val
		}

		storage := request.CreateServerStorageDevice{
			Action:  storageAction,
			Storage: storageID,
			Title:   fmt.Sprintf("%s disk 1", req.Hostname),
			Size:    storageSize,
			Tier:    storageTier,
		}

		req.StorageDevices = request.CreateServerStorageDeviceSlice{storage}

		if cmd.Flags().Changed("host") {
			req.Host, _ = cmd.Flags().GetInt("host")
		}

		loginUserDefined := false
		loginUser := request.LoginUser{
			CreatePassword: "yes",
			Username:       "root",
		}
		if cmd.Flags().Changed("login-user") {
			loginUser.Username, _ = cmd.Flags().GetString("login-user")
			loginUserDefined = true
		}
		if cmd.Flags().Changed("ssh-keys-file") {
			keysFile, _ := cmd.Flags().GetString("ssh-keys-file")
			keys, err := readFileLines(keysFile)
			if err != nil {
				return errors.Wrap(err, "Unable to read SSH keys from a file")
			}

			loginUser.SSHKeys = keys
			loginUser.CreatePassword = "no"
			loginUserDefined = true
		}
		if loginUserDefined {
			req.LoginUser = &loginUser
		}

		if cmd.Flags().Changed("userdata-file") {
			userDataFile, _ := cmd.Flags().GetString("userdata-file")
			userData, err := readFile(userDataFile)
			if err != nil {
				return errors.Wrap(err, "Unable to read userdata file")
			}
			req.UserData = userData
		}

		metadata, _ := cmd.Flags().GetBool("metadata")
		req.Metadata = upcloud.FromBool(metadata)
		// TODO firewall after library has better support

		interfaces := request.CreateServerInterfaceSlice{}
		if ok, _ := cmd.Flags().GetBool("public-ipv4"); ok {
			interfaces = append(interfaces, createServerInterface("public", "IPv4", ""))
		}
		if ok, _ := cmd.Flags().GetBool("utility"); ok {
			interfaces = append(interfaces, createServerInterface("utility", "IPv4", ""))
		}
		if ok, _ := cmd.Flags().GetBool("public-ipv6"); ok {
			interfaces = append(interfaces, createServerInterface("public", "IPv6", ""))
		}
		if len(interfaces) > 0 {
			req.Networking = &request.CreateServerNetworking{Interfaces: interfaces}
		}

		details, err := apiService.CreateServer(&req)
		if err != nil {
			fmt.Println(errors.Wrap(err, "Server creation failed"))
			return nil
		}

		if jsonOutput {
			return printServerJSON(details)
		}

		if verbose {
			fmt.Printf("Server %s (%s) created successfully in zone %s.\n", details.Hostname, details.UUID, details.Zone)
		}

		return nil
	},
}

var startServerCmd = &cobra.Command{
	Use:   "start UUID",
	Short: "Start server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Must specify a single server")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		req := &request.StartServerRequest{UUID: args[0]}
		wantedHost, _ := cmd.Flags().GetInt("host")
		if wantedHost > 0 {
			req.Host = wantedHost
		}

		details, err := apiService.StartServer(req)
		if err != nil {
			fmt.Println(errors.Wrap(err, "Starting the server failed"))
			return nil
		}

		if jsonOutput {
			return printServerJSON(details)
		}

		if verbose {
			successfulStart := fmt.Sprintf("Server %s successfully started", args[0])
			if wantedHost > 0 {
				fmt.Printf("%s on host %d.\n", successfulStart, details.Host)
			} else {
				fmt.Printf("%s.\n", successfulStart)
			}
		}

		return nil
	},
}

var stopServerCmd = &cobra.Command{
	Use:   "stop UUID",
	Short: "Stop server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Must specify a single server")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		hardStop, _ := cmd.Flags().GetBool("hard")
		softTimeout, _ := cmd.Flags().GetInt("timeout")

		req := &request.StopServerRequest{UUID: args[0]}
		if hardStop {
			req.StopType = request.ServerStopTypeHard
		} else if softTimeout > 0 {
			req.Timeout = time.Duration(softTimeout) * time.Second
		}

		details, err := apiService.StopServer(req)
		if err != nil {
			fmt.Println(errors.Wrap(err, "Server stop failed"))
			return nil
		}

		if jsonOutput {
			return printServerJSON(details)
		}

		if verbose {
			fmt.Printf("Server %s successfully stopped.\n", args[0])
		}

		return nil
	},
}

var restartServerCmd = &cobra.Command{
	Use:   "restart UUID",
	Short: "Restart server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Must specify a single server")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		req := &request.RestartServerRequest{UUID: args[0]}
		wantedHost, _ := cmd.Flags().GetInt("host")
		if wantedHost > 0 {
			req.Host = wantedHost
		}

		hardStop, _ := cmd.Flags().GetBool("hard")
		softTimeout, _ := cmd.Flags().GetInt("timeout")

		if hardStop {
			req.StopType = request.ServerStopTypeHard
		} else if softTimeout > 0 {
			req.Timeout = time.Duration(softTimeout) * time.Second
		}

		details, err := apiService.RestartServer(req)
		if err != nil {
			fmt.Println(errors.Wrap(err, "Server restart failed"))
			return nil
		}

		if jsonOutput {
			return printServerJSON(details)
		}

		if verbose {
			fmt.Printf("Server %s successfully restarted.\n", args[0])
		}

		return nil
	},
}

var deleteServerCmd = &cobra.Command{
	Use:   "delete UUID",
	Short: "Delete server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Must specify a single server")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		confirmation, _ := cmd.Flags().GetBool("confirm")
		if !confirmation {
			return errors.New("Deletion must be confirmed with --confirm")
		}

		deleteStorages, _ := cmd.Flags().GetBool("delete-attached-storages")

		var err error
		var successMsg string
		if deleteStorages {
			err = apiService.DeleteServerAndStorages(&request.DeleteServerAndStoragesRequest{UUID: args[0]})
			successMsg = fmt.Sprintf("Server %s and its storages successfully deleted.", args[0])
		} else {
			err = apiService.DeleteServer(&request.DeleteServerRequest{UUID: args[0]})
			successMsg = fmt.Sprintf("Server %s successfully deleted.", args[0])
		}

		if err != nil {
			fmt.Println(errors.Wrap(err, "Server deletion failed"))
			return nil
		}

		if verbose {
			fmt.Println(successMsg)
		}

		return nil
	},
}

func createServerInterface(ifaceType, family, network string) request.CreateServerInterface {
	iface := request.CreateServerInterface{
		IPAddresses: request.CreateServerIPAddressSlice{},
		Type:        ifaceType,
	}
	ipAddr := request.CreateServerIPAddress{Family: family}

	if network != "" {
		iface.Network = network
	}
	// TODO add ipAddr.Address once it exists in API library

	iface.IPAddresses = append(iface.IPAddresses, ipAddr)
	return iface
}

func printServerJSON(details *upcloud.ServerDetails) error {
	serverJSON, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return errors.Wrap(err, "JSON serialization failed")
	}

	fmt.Printf("%s\n", serverJSON)
	return nil
}

func getDefaultOSTemplates() map[string]string {
	templates := make(map[string]string)

	templates["debian"] = "01000000-0000-4000-8000-000020050100"  // Buster
	templates["ubuntu"] = "01000000-0000-4000-8000-000030200200"  // Focal Forssa
	templates["centos"] = "01000000-0000-4000-8000-000050010400"  // 8.0
	templates["windows"] = "01000000-0000-4000-8000-000010070300" // 2019 Standard

	return templates
}

func init() {
	defaultOSTemplates = getDefaultOSTemplates()

	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(listServersCmd)
	serverCmd.AddCommand(createServerCmd)
	serverCmd.AddCommand(showServerCmd)
	serverCmd.AddCommand(startServerCmd)
	serverCmd.AddCommand(stopServerCmd)
	serverCmd.AddCommand(restartServerCmd)
	serverCmd.AddCommand(deleteServerCmd)

	createServerCmd.Flags().String("zone", "us-sjo1", "Zone where the server is created in")
	createServerCmd.Flags().String("hostname", "", "Server hostname")
	createServerCmd.Flags().String("title", "", "Server title (defaults to hostname)")
	createServerCmd.Flags().String("plan", "2xCPU-4GB", "Server plan (see: up server plan list)")
	createServerCmd.Flags().Int("cores", 0, "Server cores (defaults to plan if it is given)")
	createServerCmd.Flags().Int("memory", 0, "Server memory in MB (defaults to plan if it is given)")
	createServerCmd.Flags().Bool("metadata", false, "Enable server metadata")
	createServerCmd.Flags().String("password-delivery", "none", "Password delivery method (email, sms, none)")
	createServerCmd.Flags().String("storage", "ubuntu", "Storage template UUID or name (ubuntu, debian, centos, windows)")
	createServerCmd.Flags().String("storage-tier", "maxiops", "Storage tier (defaults to plan)")
	createServerCmd.Flags().Int("storage-size", 0, "Storage size in GB (defaults to plan if not larger)")
	createServerCmd.Flags().String("login-user", "root", "Username to create")
	createServerCmd.Flags().String("ssh-keys-file", "", "File that contains wanted SSH keys")
	createServerCmd.Flags().String("userdata-file", "", "File with userdata script to run in the server")
	createServerCmd.Flags().Bool("public-ipv4", false, "Add a public IPv4 network interface")
	createServerCmd.Flags().Bool("public-ipv6", false, "Add a public IPv6 network interface")
	createServerCmd.Flags().Bool("utility", false, "Add a utility network interface")

	hostDescription := "Private cloud: server host to start your server to (see host section)"
	startServerCmd.Flags().Int("host", 0, hostDescription)
	restartServerCmd.Flags().Int("host", 0, hostDescription)

	hardStop := "Hard stop. Shuts down the server immediately."
	stopTimeout := "Timeout for soft stop (seconds)."
	stopServerCmd.Flags().Bool("hard", false, hardStop)
	restartServerCmd.Flags().Bool("hard", false, hardStop)
	stopServerCmd.Flags().Int("timeout", 300, stopTimeout)
	restartServerCmd.Flags().Int("timeout", 300, stopTimeout)

	deleteServerCmd.Flags().Bool("delete-attached-storages", false, "Delete storages that are attached to the server")
	deleteServerCmd.Flags().Bool("confirm", false, "Confirm deletion")
}
