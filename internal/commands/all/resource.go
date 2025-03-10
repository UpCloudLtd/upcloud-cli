package all

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/network"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/networkpeering"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/router"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
)

type Resource struct {
	Name string
	Type string
	UUID string
}

func getMatches[T any](exec commands.Executor, resolutionProvider resolver.CachingResolutionProvider[T], name string) ([]T, error) {
	resolve, err := resolutionProvider.Get(exec.Context(), exec.All())
	if err != nil {
		return nil, err
	}

	resolved := resolve(name)
	uuids, err := resolved.GetAll()
	if err != nil {
		if !errors.Is(err, resolver.NotFoundError(name)) {
			return nil, err
		}
	}

	var matches []T
	for _, uuid := range uuids {
		val, err := resolutionProvider.GetCached(uuid)
		if err != nil {
			return nil, err
		}
		matches = append(matches, val)
	}

	return matches, nil
}

func listResources(exec commands.Executor, name string) ([]Resource, error) {
	var resources []Resource

	if name == "" {
		name = "*"
	}

	// Use the same order as in hub.upcloud.com

	networks, err := getMatches(exec, &resolver.CachingNetwork{}, name)
	if err != nil {
		return nil, err
	}
	for _, network := range networks {
		resources = append(resources, Resource{
			Name: network.Name,
			Type: "network",
			UUID: network.UUID,
		})
	}

	peerings, err := getMatches(exec, &resolver.CachingNetworkPeering{}, name)
	if err != nil {
		return nil, err
	}
	for _, peering := range peerings {
		resources = append(resources, Resource{
			Name: peering.Name,
			Type: "network-peering",
			UUID: peering.UUID,
		})
	}

	routers, err := getMatches(exec, &resolver.CachingRouter{}, name)
	if err != nil {
		return nil, err
	}
	for _, router := range routers {
		if router.Type == "service" {
			continue
		}

		resources = append(resources, Resource{
			Name: router.Name,
			Type: "router",
			UUID: router.UUID,
		})
	}

	return resources, nil
}

func deleteResource(exec commands.Executor, resource Resource) (err error) {
	switch resource.Type {
	case "network":
		_, err = network.Delete(exec, resource.UUID)
	case "network-peering":
		_, err = networkpeering.Delete(exec, resource.UUID, true)
	case "router":
		_, err = router.Delete(exec, resource.UUID)
	}
	return
}

type deleteResult struct {
	Worker   int
	Resource Resource
	Error    error
}

func deleteResources(exec commands.Executor, resources []Resource, workerCount int) error {
	cfg := progress.GetDefaultOutputConfig()
	buf := bytes.NewBuffer(nil)
	cfg.Target = buf
	delExec := exec.WithProgress(progress.NewProgress(cfg))

	deleteQueue := make(chan Resource, len(resources))
	for _, resource := range resources {
		deleteQueue <- resource
	}

	workerQueue := make(chan int, workerCount)
	for n := 0; n < workerCount; n++ {
		workerQueue <- n
	}

	returnChan := make(chan deleteResult)
	deleted := make([]Resource, 0, len(resources))
	for {
		select {
		case i := <-deleteQueue:
			go func(r Resource) {
				workerID := <-workerQueue
				exec.PushProgressUpdate(messages.Update{
					Key:     r.UUID,
					Message: fmt.Sprintf("Deleting %s %s", r.Type, r.Name),
					Status:  messages.MessageStatusStarted,
				})
				defer func() {
					workerQueue <- workerID
				}()
				err := deleteResource(delExec, r)
				returnChan <- deleteResult{
					Worker:   workerID,
					Error:    err,
					Resource: r,
				}
			}(i)
		case res := <-returnChan:
			// Requeue failed deletes after 5 seconds
			if res.Error != nil {
				exec.PushProgressUpdate(messages.Update{
					Key:     res.Resource.UUID,
					Message: fmt.Sprintf("Waiting 5 seconds before retrying to delete %s %s", res.Resource.Type, res.Resource.Name),
				})
				go func(r Resource) {
					time.Sleep(5 * time.Second)
					deleteQueue <- r
				}(res.Resource)
			} else {
				exec.PushProgressSuccess(res.Resource.UUID)
				deleted = append(deleted, res.Resource)
			}

			if len(deleted) >= len(resources) {
				return nil
			}
		}
	}
}
