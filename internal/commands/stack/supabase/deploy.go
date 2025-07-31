package supabase

import (
	"embed"
	"fmt"
	"os"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack/core"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
)

//go:embed charts/supabase/*
//go:embed charts/supabase/templates/*
//go:embed charts/supabase/templates/*/*
var ChartFS embed.FS

func DeploySupabaseCommand() commands.Command {
	return &deploySupabaseCommand{
		BaseCommand: commands.New(
			"supabase",
			"Deploy a Supabase stack",
			"upctl stack deploy supabase <project-name>",
			"upctl stack deploy supabase my-new-project",
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
	msg := fmt.Sprintf("Creating supabase stack %v in zone %v", s.name, s.zone)
	exec.PushProgressStarted(msg)

	// Create a tmp dir for this deployment
	chartDir, err := os.MkdirTemp("", fmt.Sprintf("supabase-%s-%s", s.name, s.zone))
	if err != nil {
		return nil, fmt.Errorf("failed to make temp dir for deployment: %w", err)
	}

	// unpack the supabase charts into that temp dir
	if err := core.ExtractFolder(ChartFS, chartDir); err != nil {
		return nil, fmt.Errorf("failed to extract supabase chart: %w", err)
	}

	// Command implementation for deploying a Supabase stack
	config, err := s.deploy(exec, chartDir)
	if err != nil {
		fmt.Printf("Error deploying Supabase stack: %+v\n", err)
		return commands.HandleError(exec, msg, err)
	}

	// clean up at the end
	//defer os.RemoveAll(chartDir)

	exec.PushProgressSuccess(msg)

	return output.Raw(summaryOutput(config)), nil
}
