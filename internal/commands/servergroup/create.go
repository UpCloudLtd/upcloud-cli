package servergroup

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v2/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/ui"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/spf13/pflag"
)

// CreateCommand creates the "servergroup create" command
func CreateCommand() commands.Command {
	return &createCommand{
		BaseCommand: commands.New(
			"create",
			"Create a server group",
			`upctl servergroup create \
				--name my-server-group \
				--anti-affinity-policy yes \
				--label env=dev`,
			`upctl servergroup create \
				--name my-server-group \
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

func (p *createParams) processParams() error {
	p.AntiAffinityPolicy = upcloud.ServerGroupAntiAffinityPolicy(p.antiAffinityPolicy)

	if len(p.labels) > 0 {
		labelSlice, err := labels.StringsToUpCloudLabelSlice(p.labels)
		if err != nil {
			return err
		}

		p.Labels = labelSlice
	}

	if len(p.servers) > 0 {
		servers := make(upcloud.ServerUUIDSlice, 0)

		for _, v := range p.servers {
			servers = append(servers, v)
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

	fs.StringVar(&c.params.Title, "name", "", "Server group name.")
	fs.StringVar(&c.params.antiAffinityPolicy, "anti-affinity-policy", defaultCreateParams.antiAffinityPolicy, "Anti-affinity policy. Valid values are `yes` (best effort), `strict` and `no`.")
	fs.StringArrayVar(&c.params.labels, "label", defaultCreateParams.labels, "Labels to describe the server group in `key=value` format, multiple can be declared.\nUsage: --label env=dev\n\n--label owner=operations")
	fs.StringArrayVar(&c.params.servers, "server", defaultCreateParams.servers, "Servers to be added to the server group, multiple can be declared.\nUsage: --server my-server\n\n--server aa39e313-d908-418a-a959-459699bdc83b")

	c.AddFlags(fs)

	_ = c.Cobra().MarkFlagRequired("name")
}

func (c *createCommand) InitCommandWithConfig(cfg *config.Config) {
	_ = c.Cobra().RegisterFlagCompletionFunc("servers", namedargs.CompletionFunc(completion.Server{}, cfg))
}

// ExecuteWithoutArguments implements commands.NoArgumentCommand
func (c *createCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	svc := exec.All()

	if err := c.params.processParams(); err != nil {
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
