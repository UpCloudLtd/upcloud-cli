package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "List, show & control servers",
}

var listServersCmd = &cobra.Command{
	Use: "list",
	Short: "List current servers",
	RunE: func(cmd *cobra.Command, args []string) error {

		servers, err := apiService.GetServers()
		if err != nil {
			return err
		}

		if jsonOutput {
			serversJson, err := json.MarshalIndent(servers.Servers, "", "  ")
			if err != nil {
				return errors.Wrap(err, "JSON serialization failed.")
			}

			fmt.Printf("%s\n", serversJson)
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
	Use:   "show",
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
			serverJson, err := json.MarshalIndent(server, "", "  ")
			if err != nil {
				return errors.Wrap(err, "JSON serialization failed.")
			}

			fmt.Printf("%s\n", serverJson)
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

		table = tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Storage", "UUID", "Size (GB)", "Type"})
		for _, storage := range server.StorageDevices {
			table.Append([]string{storage.Title, storage.UUID, fmt.Sprintf("%d", storage.Size), storage.Type})
		}
		table.Render()

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
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO
		return nil
	},
}

var startServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO
		return nil
	},
}

var stopServerCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO
		return nil
	},
}

var restartServerCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO
		return nil
	},
}


func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(listServersCmd)
	serverCmd.AddCommand(showServerCmd)
}