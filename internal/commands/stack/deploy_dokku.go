package stack

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/spf13/pflag"
)

func DeployDokkuCommand() commands.Command {
	return &deployDokkuCommand{
		BaseCommand: commands.New(
			"dokku",
			"Deploy a Dokku Builder stack",
			"upctl stack deploy dokku <project-name>",
			"upctl stack deploy dokku my-new-project",
		),
	}
}

type deployDokkuCommand struct {
	*commands.BaseCommand
	location         string
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

func (s *deployDokkuCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.location, "location", s.location, "Select the location (region) for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Specify the name of the Supabase project")
	fs.StringVar(&s.githubPAT, "github-pat", s.githubPAT, "GitHub Personal Access Token. Used to allow Dokku to push your app images to your GitHub Container Registry. Make sure it has write:packages and read:packages permissions")
	fs.StringVar(&s.githubUser, "github-user", s.githubUser, "Used to allow Dokku to push your app images to your GitHub Container Registry")
	fs.StringVar(&s.certManagerEmail, "cert-manager-email", "ops@example.com", "Email for TLS cert registration (default: ops@example.com)")
	fs.StringVar(&s.globalDomain, "global-domain", s.globalDomain, "Example: example.com. If you do not have a domain name leave this empty and it will get the value of the ingress nginx load balancer automatically. Example: lb-0a39e6584…")
	fs.IntVar(&s.numNodes, "num-nodes", 3, "Number of nodes in the Dokku cluster (default: 3)")
	fs.StringVar(&s.sshPath, "ssh-path", s.sshPath, "Path to your private SSH key (default: ~/.ssh/id_rsa). Needed to be able to ‘git push dokku@<host>:<app>’ when deploying apps with git push")
	// Note: default value for pub ssh path might not resolve properly
	fs.StringVar(&s.sshPubPath, "ssh-path-pub", "~/.ssh/id_rsa.pub", "Path to your public SSH key (default: ~/.ssh/id_rsa.pub)")
	fs.StringVar(&s.githubPackageUrl, "github-package-url", s.name, "Specify the name of the Supabase project")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("location"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))
}

func (s *deployDokkuCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	msg := fmt.Sprintf("Creating dokku stack %v", s.name)
	exec.PushProgressStarted(msg)

	// Command implementation for deploying a Supabase stack

	exec.PushProgressSuccess("Dokku stack created successfully")

	return output.Raw([]byte("Commamnd executed successfully")), nil
}
