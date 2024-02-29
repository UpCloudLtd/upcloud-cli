package servergroup

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/completion"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/config"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/labels"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/namedargs"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/spf13/pflag"
)

type modifyCommand struct {
	*commands.BaseCommand
	params modifyParams
	resolver.CachingServerGroup
	completion.ServerGroup
}

// ModifyCommand creates the "servergroup modify" command
func ModifyCommand() commands.Command {
	return &modifyCommand{
		BaseCommand: commands.New(
			"modify",
			"Modify a server group",
			"upctl servergroup modify 8abc8009-4325-4b23-4321-b1232cd81231 --title your-server-group",
			"upctl servergroup modify my-server-group --anti-affinity-policy strict",
			`upctl servergroup modify my-server-group --server my-server-1 --server my-server-2 --server my-server-3-`,
			`upctl servergroup modify 8abc8009-4325-4b23-4321-b1232cd81231 --server 0bab98e5-b327-4ab8-ba16-738d4af7578b --server my-server-2`,
			`upctl servergroup modify my-server-group --label env=dev`,
		),
	}
}

type modifyParams struct {
	request.ModifyServerGroupRequest

	antiAffinityPolicy string
	labels             []string
	servers            []string
}

var defaultModifyParams = modifyParams{
	ModifyServerGroupRequest: request.ModifyServerGroupRequest{},
}

func (p *modifyParams) processParams(exec commands.Executor, uuid string) error {
	p.UUID = uuid

	if p.antiAffinityPolicy != "" {
		p.AntiAffinityPolicy = upcloud.ServerGroupAntiAffinityPolicy(p.antiAffinityPolicy)
	}

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
		p.Members = &servers
	}

	return nil
}

// InitCommand implements Command.InitCommand
func (c *modifyCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&c.params.antiAffinityPolicy, "anti-affinity-policy", defaultModifyParams.antiAffinityPolicy, "Anti-affinity policy. Valid values are `yes` (best effort), `strict` and `no`. Will take effect upon server start.")
	fs.StringArrayVar(&c.params.labels, "label", defaultModifyParams.labels, "Labels to describe the server in `key=value` format, multiple can be declared. If set, all the existing labels will be replaced with provided ones.\nUsage: --label env=dev\n\n--label owner=operations")
	fs.StringVar(&c.params.Title, "title", defaultModifyParams.Title, "New server group title.")
	fs.StringArrayVar(&c.params.servers, "server", defaultModifyParams.servers, "Servers that belong to the server group, multiple can be declared. If set, all the existing server entries will be replaced with provided ones.\nUsage: --server my-server\n\n--server 00333d1b-3a4a-4b75-820a-4a56d70395dd")

	c.AddFlags(fs)
}

// Execute implements commands.MultipleArgumentCommand
func (c *modifyCommand) Execute(exec commands.Executor, uuid string) (output.Output, error) {
	svc := exec.All()

	err := c.params.processParams(exec, uuid)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Modifying server group %s", uuid)
	exec.PushProgressStarted(msg)

	res, err := svc.ModifyServerGroup(exec.Context(), &c.params.ModifyServerGroupRequest)
	if err != nil {
		return commands.HandleError(exec, msg, err)
	}

	exec.PushProgressSuccess(msg)

	return output.OnlyMarshaled{Value: res}, nil
}

func (c *modifyCommand) InitCommandWithConfig(cfg *config.Config) {
	_ = c.Cobra().RegisterFlagCompletionFunc("server", namedargs.CompletionFunc(completion.Server{}, cfg))
}
