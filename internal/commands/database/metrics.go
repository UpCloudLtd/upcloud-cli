package database

import (
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/format"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type metricsCommand struct {
	*commands.BaseCommand
	resolver.CachingDatabase
	completion.Database

	period     string
	latestOnly config.OptionalBoolean
}

// MetricsCommand creates the "database metrics" command
func MetricsCommand() commands.Command {
	return &metricsCommand{
		BaseCommand: commands.New(
			"metrics",
			"Show database metrics",
			"upctl database metrics 0497728e-76ef-41d0-997f-fa9449eb71bc --period day",
			"upctl database metrics my_database --period month",
		),
	}
}

// InitCommand implements Command.InitCommand
func (c *metricsCommand) InitCommand() {
	periods := []string{"hour", "day", "week", "month", "year"}

	flags := &pflag.FlagSet{}
	flags.StringVar(&c.period, "period", "day", "Period for the metrics. Valid values are "+namedargs.ValidValuesHelp(periods...)+".")
	config.AddToggleFlag(flags, &c.latestOnly, "latest-only", false, "Only output the latest data point. Only affects the human readable output.")
	c.AddFlags(flags)

	commands.Must(c.Cobra().MarkFlagRequired("period"))
	commands.Must(c.Cobra().RegisterFlagCompletionFunc("period", cobra.FixedCompletions(periods, cobra.ShellCompDirectiveNoFileComp)))
}

func float64MetricsToTable(metrics upcloud.ManagedDatabaseMetricsChartFloat64, latestOnly bool) output.CombinedSection {
	columns := make([]output.TableColumn, len(metrics.Columns)+1)
	columns[0] = output.TableColumn{Key: "time", Header: "Time"}
	for i, column := range metrics.Columns {
		columns[i+1] = output.TableColumn{Key: column.Label, Header: column.Label, Format: format.RoundedNumber(2)}
	}

	N := len(metrics.Rows)
	if latestOnly && N > 0 {
		N = 1
	}

	rows := make([]output.TableRow, N)
	for i, row := range metrics.Rows[len(metrics.Rows)-N:] {
		rows[i] = make(output.TableRow, len(row)+1)
		rows[i][0] = metrics.Timestamps[len(metrics.Timestamps)-N+i].Format("2006-01-02 15:04:05")
		for j, value := range row {
			rows[i][j+1] = value
		}
	}

	return output.CombinedSection{
		Title: metrics.Title,
		Contents: output.Table{
			Columns: columns,
			Rows:    rows,
		},
	}
}

// Execute implements commands.MultipleArgumentCommand
func (c *metricsCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()
	metrics, err := svc.GetManagedDatabaseMetrics(exec.Context(), &request.GetManagedDatabaseMetricsRequest{UUID: uuid, Period: upcloud.ManagedDatabaseMetricPeriod(c.period)})
	if err != nil {
		return nil, err
	}
	if metrics == nil {
		panic("GetManagedDatabaseMetrics returned nil metrics")
	}

	latestOnly := c.latestOnly.Value()

	// For JSON and YAML output, passthrough API response
	return output.MarshaledWithHumanOutput{
		Value: metrics,
		Output: output.Combined{
			float64MetricsToTable(metrics.CPUUsage, latestOnly),
			float64MetricsToTable(metrics.MemoryUsage, latestOnly),
			float64MetricsToTable(metrics.DiskUsage, latestOnly),
		},
	}, nil
}
