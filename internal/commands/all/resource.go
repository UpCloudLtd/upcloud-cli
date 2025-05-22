package all

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"sort"
	"sync"
	"time"

	"github.com/UpCloudLtd/progress"
	"github.com/UpCloudLtd/progress/messages"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/database"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/kubernetes"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/loadbalancer"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/network"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/networkpeering"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/objectstorage"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/router"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/server"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/servergroup"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/commands/storage"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
)

const (
	includeHelp = "Include resources matching the given name. If defined multiple times, resource is included if it matches any of the given names. `*` matches all resources."
	excludeHelp = "Exclude resources matching the given name. If defined multiple times, resource is included if it matches any of the given names."

	typeKubernetes     = "kubernetes-cluster"
	typeLoadBalancer   = "load-balancer"
	typeNetwork        = "network"
	typeNetworkPeering = "network-peering"
	typeRouter         = "router"
	typeObjectStorage  = "object-storage"
	typeDatabase       = "database"
	typeServer         = "server"
	typeServerGroup    = "server-group"
	typeStorage        = "storage"
	typeTag            = "tag"
)

type Resource struct {
	Name  string
	Type  string
	UUID  string
	State string
}

func (r Resource) Key() string {
	key := r.UUID
	if key == "" {
		key = r.Name
	}
	return key
}

func setAdd(a, b []string) []string {
	for _, i := range b {
		if !slices.Contains(a, i) {
			a = append(a, i)
		}
	}
	return a
}

func setRemove(a, b []string) []string {
	var res []string
	for _, i := range a {
		if !slices.Contains(b, i) {
			res = append(res, i)
		}
	}
	return res
}

func getMatches[T any](exec commands.Executor, resolutionProvider resolver.CachingResolutionProvider[T], include, exclude []string) ([]T, error) {
	resolve, err := resolutionProvider.Get(exec.Context(), exec.All())
	if err != nil {
		return nil, err
	}

	var uuids []string
	for _, i := range include {
		resolved := resolve(i)
		toAdd, err := resolved.GetAll()
		uuids = setAdd(uuids, toAdd)
		if err != nil {
			if !errors.Is(err, resolver.NotFoundError(i)) {
				return nil, err
			}
		}
	}

	for _, i := range exclude {
		resolved := resolve(i)
		toRemove, err := resolved.GetAll()
		uuids = setRemove(uuids, toRemove)
		if err != nil {
			if !errors.Is(err, resolver.NotFoundError(i)) {
				return nil, err
			}
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

type findResult struct {
	Resources []Resource
	Error     error
}

func findResources[T any](exec commands.Executor, wg *sync.WaitGroup, returnChan chan findResult, r resolver.CachingResolutionProvider[T], include, exclude []string) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		var resources []Resource
		matches, err := getMatches(exec, r, include, exclude)
		if err != nil {
			returnChan <- findResult{Resources: nil, Error: err}
			return
		}
		for _, match := range matches {
			resource, err := getResource(match)
			if err != nil {
				returnChan <- findResult{Resources: nil, Error: err}
				return
			}
			resources = append(resources, resource)
		}
		returnChan <- findResult{Resources: resources, Error: nil}
	}()
}

func listResources(exec commands.Executor, include, exclude []string) ([]Resource, error) {
	var resources []Resource
	returnChan := make(chan findResult, 11)

	var wg sync.WaitGroup

	findResources(exec, &wg, returnChan, &resolver.CachingKubernetes{}, include, exclude)
	findResources(exec, &wg, returnChan, &resolver.CachingLoadBalancer{}, include, exclude)
	findResources(exec, &wg, returnChan, &resolver.CachingNetwork{}, include, exclude)
	findResources(exec, &wg, returnChan, &resolver.CachingNetworkPeering{}, include, exclude)
	findResources(exec, &wg, returnChan, &resolver.CachingRouter{Type: "normal"}, include, exclude)
	findResources(exec, &wg, returnChan, &resolver.CachingObjectStorage{}, include, exclude)
	findResources(exec, &wg, returnChan, &resolver.CachingDatabase{}, include, exclude)
	findResources(exec, &wg, returnChan, &resolver.CachingServer{}, include, exclude)
	findResources(exec, &wg, returnChan, &resolver.CachingServerGroup{}, include, exclude)
	findResources(exec, &wg, returnChan, &resolver.CachingStorage{Access: "private"}, include, exclude)
	findResources(exec, &wg, returnChan, &cachingTag{}, include, exclude)

	wg.Wait()
	close(returnChan)

	for res := range returnChan {
		if res.Error != nil {
			return nil, res.Error
		}
		resources = append(resources, res.Resources...)
	}

	sort.Slice(resources, func(i, j int) bool {
		if resources[i].Type != resources[j].Type {
			return resources[i].Type < resources[j].Type
		}

		return resources[i].Name < resources[j].Name
	})

	return resources, nil
}

func getResource(val any) (Resource, error) {
	switch v := val.(type) {
	case upcloud.LoadBalancer:
		return Resource{
			Name: v.Name,
			Type: typeLoadBalancer,
			UUID: v.UUID,
		}, nil
	case upcloud.KubernetesCluster:
		return Resource{
			Name: v.Name,
			Type: typeKubernetes,
			UUID: v.UUID,
		}, nil
	case upcloud.Network:
		return Resource{
			Name: v.Name,
			Type: typeNetwork,
			UUID: v.UUID,
		}, nil
	case upcloud.NetworkPeering:
		return Resource{
			Name: v.Name,
			Type: typeNetworkPeering,
			UUID: v.UUID,
		}, nil
	case upcloud.Router:
		return Resource{
			Name: v.Name,
			Type: typeRouter,
			UUID: v.UUID,
		}, nil
	case upcloud.ManagedObjectStorage:
		return Resource{
			Name: v.Name,
			Type: typeObjectStorage,
			UUID: v.UUID,
		}, nil
	case upcloud.ManagedDatabase:
		return Resource{
			Name: v.Title,
			Type: typeDatabase,
			UUID: v.UUID,
		}, nil
	case upcloud.Server:
		return Resource{
			Name:  v.Title,
			Type:  typeServer,
			UUID:  v.UUID,
			State: v.State,
		}, nil
	case upcloud.ServerGroup:
		return Resource{
			Name: v.Title,
			Type: typeServerGroup,
			UUID: v.UUID,
		}, nil
	case upcloud.Storage:
		return Resource{
			Name: v.Title,
			Type: typeStorage,
			UUID: v.UUID,
		}, nil
	case upcloud.Tag:
		return Resource{
			Name: v.Name,
			Type: typeTag,
		}, nil
	}
	return Resource{}, fmt.Errorf("unsupported type %T", val)
}

func deleteResource(exec commands.Executor, resource Resource) (err error) {
	switch resource.Type {
	case typeKubernetes:
		_, err = kubernetes.Delete(exec, resource.UUID, true)
	case typeLoadBalancer:
		_, err = loadbalancer.Delete(exec, resource.UUID, true)
	case typeNetwork:
		_, err = network.Delete(exec, resource.UUID)
	case typeNetworkPeering:
		_, err = networkpeering.Delete(exec, resource.UUID, true)
	case typeRouter:
		_, err = router.Delete(exec, resource.UUID)
	case typeObjectStorage:
		_, err = objectstorage.Delete(exec, resource.UUID, true, true, true, true)
	case typeDatabase:
		_, err = database.Delete(exec, resource.UUID, true, true)
	case typeServer:
		_, err = server.Delete(exec, resource.UUID, resource.State, false, true)
	case typeServerGroup:
		_, err = servergroup.Delete(exec, resource.UUID)
	case typeStorage:
		_, err = storage.Delete(exec, resource.UUID, "delete")
	case typeTag:
		_, err = deleteTag(exec, resource.Name)
	}
	return
}

type deleteResult struct {
	Worker   int
	Resource Resource
	Error    error
}

func deleteResources(exec commands.Executor, resources []Resource, workerCount int) error {
	if len(resources) == 0 {
		return nil
	}

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
					Key:     r.Key(),
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
					Key:     res.Resource.Key(),
					Message: fmt.Sprintf("Waiting 5 seconds before retrying to delete %s %s", res.Resource.Type, res.Resource.Name),
				})
				go func(r Resource) {
					time.Sleep(5 * time.Second)
					deleteQueue <- r
				}(res.Resource)
			} else {
				exec.PushProgressSuccess(res.Resource.Key())
				deleted = append(deleted, res.Resource)
			}

			if len(deleted) >= len(resources) {
				return nil
			}
		}
	}
}
