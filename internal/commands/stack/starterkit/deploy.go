package starterkit

import (
	"fmt"
	"os"

	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/database"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/kubernetes"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/stack"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/output"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

func DeployStarterKitCommand() commands.Command {
	return &deployStarterKitCommand{
		BaseCommand: commands.New(
			"starter-kit",
			"Deploy a Starter Kit stack",
			"upctl stack deploy starter-kit --name <project-name> --zone <zone-name>",
			"upctl stack deploy starter-kit --name my-new-project --zone es-mad1",
		),
	}
}

type deployStarterKitCommand struct {
	*commands.BaseCommand
	zone string
	name string
}

func (s *deployStarterKitCommand) InitCommand() {
	fs := &pflag.FlagSet{}
	fs.StringVar(&s.zone, "zone", s.zone, "Zone for the stack deployment")
	fs.StringVar(&s.name, "name", s.name, "Starter Kit stack name")
	s.AddFlags(fs)

	commands.Must(s.Cobra().MarkFlagRequired("zone"))
	commands.Must(s.Cobra().MarkFlagRequired("name"))
}

func (s *deployStarterKitCommand) ExecuteWithoutArguments(exec commands.Executor) (output.Output, error) {
	// Command implementation for deploying a Starter Kit stack
	summary, err := s.deploy(exec)
	if err != nil {
		return nil, fmt.Errorf("deploying Starter Kit stack: %w", err)
	}

	return output.Raw(summary), nil
}

func (s *deployStarterKitCommand) deploy(exec commands.Executor) (string, error) {
	// Generate configuration
	config := GenerateStarterKitConfig(s.name, s.zone)

	// Validate configuration
	if err := config.Validate(exec); err != nil {
		return "", err
	}

	// Create a tmp dir for this deployment
	projectDir, err := os.MkdirTemp("", fmt.Sprintf("starter-kit-%s-%s", s.name, s.zone))
	if err != nil {
		return "", fmt.Errorf("failed to make temp dir for deployment: %w", err)
	}

	msg := "Creating network for starter kit"
	exec.PushProgressStarted(msg)

	network, err := stack.CreateNetwork(exec, config.NetworkName, config.Zone, stack.StackTypeStarterKit)
	if err != nil {
		return "", fmt.Errorf("failed to create network: %w for starter kit deployment", err)
	}

	if network == nil {
		return "", fmt.Errorf("created network %s is nil", config.NetworkName)
	}
	if len(network.IPNetworks) == 0 {
		return "", fmt.Errorf("created network %s has no IP networks", config.NetworkName)
	}

	// Create the router
	router, err := exec.All().CreateRouter(exec.Context(), &request.CreateRouterRequest{
		Name: config.RouterName,
		Labels: []upcloud.Label{
			{Key: "stacks.upcloud.com/stack", Value: string(stack.StackTypeStarterKit)},
			{Key: "stacks.upcloud.com/created-by", Value: "upctl"},
			{Key: "stacks.upcloud.com/version", Value: string(stack.VersionV0_1_0_0)},
			{Key: "stacks.upcloud.com/name", Value: config.RouterName},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create router: %w for starter kit deployment", err)
	}

	// Attach the router to the network
	err = exec.All().AttachNetworkRouter(exec.Context(), &request.AttachNetworkRouterRequest{
		NetworkUUID: network.UUID,
		RouterUUID:  router.UUID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to attach router to network: %w", err)
	}

	exec.PushProgressSuccess(msg)

	// Deploy kubernetes, database and object storage in parallel
	eg, ctx := errgroup.WithContext(exec.Context())
	var cluster *upcloud.KubernetesCluster
	var db *upcloud.ManagedDatabase
	var objStorage *upcloud.ManagedObjectStorage
	var kubeconfigPath string

	// Create kubernetes cluster
	eg.Go(func() error {
		cluster, kubeconfigPath, err = createKubernetes(ctx, exec, config, network, projectDir)
		if err != nil {
			return fmt.Errorf("failed to create kubernetes cluster: %w", err)
		}
		return nil
	})

	// Create the database
	eg.Go(func() error {
		db, err = createDatabase(ctx, exec, config, network)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		return nil
	})

	// Create object storage
	eg.Go(func() error {
		objStorage, err = createObjectStorage(ctx, exec, config, network)
		if err != nil {
			return fmt.Errorf("failed to create object storage: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return "", err
	}

	// Wait for cluster, db and obj storage before getting summary information
	// Wait for database to be ready
	msg = fmt.Sprintf("Waiting for database with UUID: %s to be ready ", db.UUID)
	exec.PushProgressStarted(msg)

	database.WaitForManagedDatabaseState(db.UUID, upcloud.ManagedDatabaseStateRunning, exec, msg)
	db, err = exec.All().GetManagedDatabase(exec.Context(), &request.GetManagedDatabaseRequest{
		UUID: db.UUID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get database info: %w with UUID %s", err, db.UUID)
	}

	exec.PushProgressSuccess(msg)

	// Wait for kubernetes cluster to be ready
	msg = fmt.Sprintf("Waiting for kubernetes cluster with UUID: %s to be ready ", cluster.UUID)
	exec.PushProgressStarted(msg)

	kubernetes.WaitForClusterState(cluster.UUID, upcloud.KubernetesClusterStateRunning, exec, msg)
	cluster, err = exec.All().GetKubernetesCluster(exec.Context(), &request.GetKubernetesClusterRequest{
		UUID: cluster.UUID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get kubernetes cluster info: %w with UUID %s", err, cluster.UUID)
	}

	exec.PushProgressSuccess(msg)

	// Wait for object storage to be ready
	msg = fmt.Sprintf("Waiting for object storage with UUID: %s to be ready ", objStorage.UUID)
	exec.PushProgressStarted(msg)
	stack.WaitForManagedObjectStorageState(objStorage.UUID, upcloud.ManagedObjectStorageOperationalStateRunning, exec, msg)
	objStorage, err = exec.All().GetManagedObjectStorage(exec.Context(), &request.GetManagedObjectStorageRequest{
		UUID: objStorage.UUID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get object storage info: %w with UUID %s", err, objStorage.UUID)
	}

	exec.PushProgressSuccess(msg)

	// Create summary and return
	return buildSummary(cluster, kubeconfigPath, network, router, db, objStorage), nil
}
