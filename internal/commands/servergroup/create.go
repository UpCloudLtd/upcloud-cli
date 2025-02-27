package servergroup

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CreateCommand creates the "servergroup create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a server group",
			`upctl server-group create \
				--title my-server-group \
				--anti-affinity-policy yes \
				--server 1fdfda29-ead1-4855-b71f-a432179800ab \
				--server my-server`,
			`upctl server-group create \
				--title my-server-group \
				--anti-affinity-policy yes \
				--label env=dev`,
			`upctl server-group create \
				--title my-server-group \
				--anti-affinity-policy strict \
				--label env=dev \
				--label owner=operations`,
		),
	}
}

type createParams struct {
	request.CreateServerGroupRequest

	antiAffinityPolicy string
	labels             []string
	servers            []string
}

var defaultCreateParams = &createParams{
	CreateServerGroupRequest: request.CreateServerGroupRequest{},
	antiAffinityPolicy:       string(upcloud.ServerGroupAntiAffinityPolicyBestEffort),
}

func (p *createParams) processParams(exec commands.Executor) error {
	p.AntiAffinityPolicy = upcloud.ServerGroupAntiAffinityPolicy(p.antiAffinityPolicy)

	if len(p.labels) > 0 {
		labelSlice, err := labels.StringsToUpCloudLabelSlice(p.labels)
		if err != nil {
			return err
		}

		p.Labels = labelSlice
	}

	if len(p.servers) > 0 {
		servers, err := stringsToServerUUIDSlice(exec, p.servers)
		if err != nil {
			return err
		}
		p.Members = servers
	}

	return nil
}

type createCommand struct {
	*commands.BaseCommand
	params createParams
}

// InitCommand implements Command.InitCommand
func (c *createCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	c.params = createParams{CreateServerGroupRequest: request.CreateServerGroupRequest{}}

	fs.StringVar(&c.params.Title, "title", defaultCreateParams.Title, "Server group title.")
	aaPolicies := []string{"yes", "strict", "no"}
	fs.StringVar(&c.params.antiAffinityPolicy, "anti-affinity-policy", defaultCreateParams.antiAffinityPolicy, "Anti-affinity policy. Valid values are "+namedargs.ValidValuesHelp(aaPolicies...)+". Will take effect upon server start.")
	fs.StringArrayVar(&c.params.labels, "label", defaultCreateParams.labels, "Labels to describe the server group in `key=value` format, multiple can be declared.\nUsage: --label env=dev\n\n--label owner=operations")
	fs.StringArrayVar(&c.params.servers, "server", defaultCreateParams.servers, "Servers to be added to the server group, multiple can be declared.\nUsage: --server my-server\n\n--server 00333d1b-3a4a-4b75-820a-4a56d70395dd")

	c.AddFlags(fs)

	commands.Must(c.Cobra().MarkFlagRequired("title"))
	commands.Must(c.Cobra().RegisterFlagCompletionFunc("title", cobra.NoFileCompletions))
	commands.Must(c.Cobra().RegisterFlagCompletionFunc("anti-affinity-policy", cobra.FixedCompletions(aaPolicies, cobra.ShellCompDirectiveNoFileComp)))
	commands.Must(c.Cobra().RegisterFlagCompletionFunc("label", cobra.NoFileCompletions))

	// Deprecating servergroup in favour of server-group
	// TODO: Remove this in the future
	commands.SetSubcommandDeprecationHelp(c, []string{"servergroup"})
}

func (c *createCommand) InitCommandWithConfig(cfg *config.Config) {
	commands.Must(c.Cobra().RegisterFlagCompletionFunc("server", namedargs.CompletionFunc(completion.Server{}, cfg)))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Deprecating servergroup in favour of server-group
	// TODO: Remove this in the future
	commands.SetSubcommandExecutionDeprecationMessage(c, []string{"servergroup"}, "server-group")

	svc := exec.All()

	if err := c.params.processParams(exec); err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Creating server group %s", c.params.Title)
	exec.PushProgressStarted(msg)

	r := c.params.CreateServerGroupRequest
	res, err := svc.CreateServerGroup(exec.Context(), &r)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.MarshaledWithHumanDetails{Value: res, Details: []output.DetailRow{
		{Title: "UUID", Value: res.UUID, Colour: ui.DefaultUUUIDColours},
	}}, nil
}
