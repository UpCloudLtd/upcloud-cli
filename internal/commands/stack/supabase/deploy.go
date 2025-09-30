package supabase

import (
	"embed"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
)

//go:embed charts/supabase/*
//go:embed charts/supabase/templates/*
//go:embed charts/supabase/templates/*/*
var chartFS embed.FS

func DeploySupabaseCommand() commands.Command {
	return &deploySupabaseCommand{
		BaseCommand: commands.New(
			"supabase",
			"Deploy a Supabase stack",
			"upctl stack deploy supabase --name <project-name> --zone <zone-name>",
			"upctl stack deploy supabase --name my-new-project --zone es-mad1",
		),
	}
}

type deploySupabaseCommand struct {
	*commands.BaseCommand
	zone       string
	name       string
	configPath string
}

func (s *deploySupabaseCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.zone, "zone", s.zone, "Zone for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Supabase stack name")
	fs.StringVar(&s.configPath, "configPath", s.configPath, "Optional path to a configuration file for the Supabase stack")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))

	// configPath is optional, but if provided, it should be a valid path
	if s.configPath != "" {
		if _, err := os.Stat(s.configPath); os.IsNotExist(err) {
			commands.Must(s.Cobra().MarkFlagRequired("configPath"))
		}
	}
}

func (s *deploySupabaseCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Create a tmp dir for this deployment
	chartDir, err := os.MkdirTemp("", fmt.Sprintf("supabase-%s-%s", s.name, s.zone))
	if err != nil {
		return nil, fmt.Errorf("failed to make temp dir for deployment: %w", err)
	}

	// unpack the supabase charts into that temp dir
	if err := stack.ExtractFolder(chartFS, chartDir); err != nil {
		return nil, fmt.Errorf("failed to extract supabase chart: %w", err)
	}

	// Command implementation for deploying a Supabase stack
	config, err := s.deploy(exec, chartDir)
	if err != nil {
		return nil, fmt.Errorf("deploying Supabase stack: %w", err)
	}

	// Build summary text
	summary := summaryOutput(config)

	// Save summary to file
	filename := fmt.Sprintf("supabase-%s-%s-%s.cfg", s.name, s.zone, time.Now().Format("20060102-150405"))
	if err := os.WriteFile(filename, []byte(summary), 0o600); err != nil {
		return nil, fmt.Errorf("failed to write summary file: %w", err)
	}

	// Add info about file to screen output
	summaryWithFileNotice := summary + fmt.Sprintf("\nConfiguration details also saved to: %s\n", filename)

	return output.Raw{
		Source: io.NopCloser(strings.NewReader(summaryWithFileNotice)),
	}, nil
}
