package resolver_test

import (
	"errors"
	smock "github.com/UpCloudLtd/cli/internal/mock"
	"github.com/UpCloudLtd/cli/internal/resolver"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerResolution(t *testing.T) {
	var Server1 = upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server-1-hostname",
		License:      0,
		MemoryAmount: 1024,
		Plan:         "server-1-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-1-title",
		UUID:         "1fdfda29-ead1-4855-b71f-1e33eb2ca9de",
		Zone:         "fi-hel1",
	}

	var Server2 = upcloud.Server{
		CoreNumber:   1,
		Hostname:     "server-2-hostname",
		License:      0,
		MemoryAmount: 1024,
		Plan:         "server-2-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-2-title",
		UUID:         "f77a5b25-84af-4f52-bc40-581930091fad",
		Zone:         "fi-hel1",
	}

	var Server3 = upcloud.Server{
		CoreNumber:   2,
		Hostname:     "server-3-hostname",
		License:      0,
		MemoryAmount: 4096,
		Plan:         "server-3-plan",
		Progress:     0,
		State:        "stopped",
		Tags:         nil,
		Title:        "server-3-title",
		UUID:         "f0131b8f-ffe0-4271-83a8-c75b99e168c3",
		Zone:         "hu-bud1",
	}

	var Server4 = upcloud.Server{
		CoreNumber:   4,
		Hostname:     "server-4-hostname",
		License:      0,
		MemoryAmount: 5120,
		Plan:         "server-4-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        Server1.Title,
		UUID:         "e5b3a855-cd8a-45b6-8cef-c7c860a02217",
		Zone:         "uk-lon1",
	}

	var Server5 = upcloud.Server{
		CoreNumber:   4,
		Hostname:     Server4.Hostname,
		License:      0,
		MemoryAmount: 5120,
		Plan:         "server-5-plan",
		Progress:     0,
		State:        "started",
		Tags:         nil,
		Title:        "server-5-title",
		UUID:         "39bc2725-213d-46c8-8b25-49990c6966a7",
		Zone:         "uk-lon1",
	}

	var allServers = &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
			Server2,
			Server3,
			Server4,
			Server5,
		},
	}
	var unambiguousServers = []upcloud.Server{
		Server2,
		Server3,
	}

	t.Run("resolve uuid", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServers").Return(allServers, nil)
		res := resolver.CachingServer{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)
		for _, srv := range allServers.Servers {
			resolved, err := argResolver(srv.UUID)
			assert.NoError(t, err)
			assert.Equal(t, srv.UUID, resolved)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetServers", 1)
	})

	t.Run("resolve hostname", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServers").Return(allServers, nil)
		res := resolver.CachingServer{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)
		for _, srv := range unambiguousServers {
			resolved, err := argResolver(srv.Hostname)
			assert.NoError(t, err)
			assert.Equal(t, srv.UUID, resolved)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetServers", 1)
	})

	t.Run("resolve title", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServers").Return(allServers, nil)
		res := resolver.CachingServer{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)
		for _, srv := range unambiguousServers {
			resolved, err := argResolver(srv.Title)
			assert.NoError(t, err)
			assert.Equal(t, srv.UUID, resolved)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetServers", 1)
	})

	t.Run("failure situations", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServers").Return(allServers, nil)

		res := resolver.CachingServer{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)

		// ambigous hostname
		resolved, err := argResolver(Server4.Hostname)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(Server4.Hostname))
		assert.Equal(t, "", resolved)

		// ambigous title
		resolved, err = argResolver(Server1.Title)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(Server1.Title))
		assert.Equal(t, "", resolved)

		// not found
		resolved, err = argResolver("notfounf")
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("notfounf"))
		assert.Equal(t, "", resolved)

		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetServers", 1)
	})
}

func TestFailingServerResolution(t *testing.T) {
	mService := &smock.Service{}
	var nilResponse *upcloud.Servers
	mService.On("GetServers").Return(nilResponse, errors.New("MOCKERROR"))
	res := resolver.CachingServer{}
	argResolver, err := res.Get(mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
