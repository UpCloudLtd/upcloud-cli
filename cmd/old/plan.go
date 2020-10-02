package old

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Server plans",
}

var listPlansCmd = &cobra.Command{
	Use:   "list",
	Short: "List available plans",
	RunE: func(cmd *cobra.Command, args []string) error {

		plans, err := apiService.GetPlans()
		if err != nil {
			return err
		}

		if jsonOutput {
			plansJSON, err := json.MarshalIndent(plans.Plans, "", "  ")
			if err != nil {
				return errors.Wrap(err, "JSON serialization failed")
			}

			fmt.Printf("%s\n", plansJSON)
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Plan", "Cores", "Memory (MB)", "Storage (GB)", "Storage tier", "Traffic (out)"})

		for _, plan := range plans.Plans {
			table.Append([]string{
				plan.Name,
				fmt.Sprintf("%d", plan.CoreNumber),
				fmt.Sprintf("%d", plan.MemoryAmount),
				fmt.Sprintf("%d", plan.StorageSize),
				plan.StorageTier,
				fmt.Sprintf("%d", plan.PublicTrafficOut),
			})
		}
		table.Render()

		return nil
	},
}

func init() {
	serverCmd.AddCommand(planCmd)

	planCmd.AddCommand(listPlansCmd)
}
