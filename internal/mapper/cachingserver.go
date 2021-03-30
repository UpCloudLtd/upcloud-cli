package mapper

import (
	"fmt"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

// CachingServer maps server by title, hostname or uuid to uuids. It caches the results from backend so servers are only queried once.
func CachingServer(svc service.Server) (Argument, error) {
	servers, err := svc.GetServers()
	if err != nil {
		return nil, err
	}
	return func(arg string) (uuid string, err error) {
		rv := ""
		for _, server := range servers.Servers {
			if server.Title == arg || server.Hostname == arg || server.UUID == arg {
				if rv != "" {
					return "", fmt.Errorf("'%v' is ambiguous, found multiple servers matching", arg)
				}
				rv = server.UUID
			}
		}
		if rv != "" {
			return rv, nil
		}
		return "", fmt.Errorf("no server found matching '%v'", arg)
	}, nil
}
