package dokku

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
)

//go:embed config/**
var dokkuChartFS embed.FS

func DeployDokkuCommand() commands.Command {
	return &deployDokkuCommand{
		BaseCommand: commands.New(
			"dokku",
			"Deploy a Dokku Builder stack",
			"upctl stack deploy dokku --zone <zone-name> --name <project-name> --github-pat <github-personal-access-token> --github-user <github-username>",
			"upctl stack deploy dokku --zone pl-waw1 --name my-dokku-project --github-pat ghp_Uiej1N1fA1W... --github-user dokkumaster",
			"upctl stack deploy dokku --name my-new-project --zone es-mad1 --github-pat ghp_Uiej1N1fA1W... --github-user dokkumaster",
		),
	}
}

type deployDokkuCommand struct {
	*commands.BaseCommand
	zone             string
	name             string
	githubPAT        string
	githubUser       string
	certManagerEmail string
	globalDomain     string
	numNodes         int
	sshPath          string
	sshPubPath       string
	githubPackageUrl string
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Handle error as you prefer, e.g., log it or return a default value
		return ""
	}
	return homeDir
}

func (s *deployDokkuCommand) InitCommand() {
	defaultSSH := filepath.Join(getHomeDir(), ".ssh", "id_rsa")
	defaultSSHPub := filepath.Join(getHomeDir(), ".ssh", "id_rsa.pub")

	fs := &pflag.FlagSet{}
	fs.StringVar(&s.zone, "zone", s.zone, "Zone (location) for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Dokku project name")
	fs.StringVar(&s.githubPAT, "github-pat", s.githubPAT, "GitHub Personal Access Token. Used to allow Dokku to push your app images to your GitHub Container Registry. Make sure it has write:packages and read:packages permissions")
	fs.StringVar(&s.githubUser, "github-user", s.githubUser, "Used to allow Dokku to push your app images to your GitHub Container Registry")
	fs.StringVar(&s.certManagerEmail, "cert-manager-email", "ops@example.com", "Email for TLS cert registration (default: ops@example.com)")
	fs.StringVar(&s.globalDomain, "global-domain", s.globalDomain, "Example: example.com. If you do not have a domain name leave this empty and it will get the value of the ingress nginx load balancer automatically. Example: lb-0a39e6584…")
	fs.IntVar(&s.numNodes, "num-nodes", 3, "Number of nodes in the Dokku cluster (default: 3)")
	fs.StringVar(&s.sshPath, "ssh-path", defaultSSH, "Path to your private SSH key (default: ~/.ssh/id_rsa). Needed to be able to ‘git push dokku@<host>:<app>’ when deploying apps with git push")
	fs.StringVar(&s.sshPubPath, "ssh-path-pub", defaultSSHPub, "Path to your public SSH key (default: ~/.ssh/id_rsa.pub)")
	fs.StringVar(&s.githubPackageUrl, "github-package-url", "ghcr.io", "Container registry hostname (default: ghcr.io)")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))
	commands.Must(s.Cobra().MarkFlagRequired("github-pat"))
	commands.Must(s.Cobra().MarkFlagRequired("github-user"))
}

func (s *deployDokkuCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Create a tmp dir for this deployment
	configDir, err := os.MkdirTemp("", fmt.Sprintf("dokku-%s-%s", s.name, s.zone))
	if err != nil {
		return nil, fmt.Errorf("failed to make temp dir for deployment: %w", err)
	}

	// unpack the dokku charts and config files into that temp dir
	if err := stack.ExtractFolder(dokkuChartFS, configDir); err != nil {
		return nil, fmt.Errorf("failed to extract dokku charts and configuration files: %w", err)
	}

	if err = s.deploy(exec, configDir); err != nil {
		return nil, fmt.Errorf("failed to deploy dokku stack: %w , if the issue persist contact us at: %s", err, stack.SupportEmail)
	}

	return output.Raw([]byte("Command executed successfully")), nil
}
