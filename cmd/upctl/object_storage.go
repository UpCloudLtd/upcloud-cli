package cmd

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/spf13/cobra"
)

var ObjectStorageCmd = &cobra.Command{
	Use:   "object-storage",
	Short: "Manage object storage buckets",
	Long:  `Work with UpCloud Object Storage (S3-compatible) buckets.`,
}

func init() {
	rootCmd.AddCommand(ObjectStorageCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all object storage buckets",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := commands.NewExecutor(cmd).AllServices()
			res, err := svc.GetObjectStorages(cmd.Context())
			if err != nil {
				return err
			}

			if len(res.ObjectStorages) == 0 {
				fmt.Println("No object storage buckets found")
				return nil
			}

			fmt.Printf("%-36s %-30s %-10s %-8s %s\n", "UUID", "Name", "Zone", "Size(GB)", "State")
			fmt.Println("------------------------------------------------------------------------------------------------")
			for _, b := range res.ObjectStorages {
				fmt.Printf("%-36s %-30s %-10s %-8d %s\n",
					b.UUID, b.Name, b.Zone, b.ConfiguredSize, b.State)
			}
			return nil
		},
	}
	ObjectStorageCmd.AddCommand(listCmd)
}
