package resolver_test

import (
	"context"
	"errors"
	"testing"

	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestServerResolution(t *testing.T) {
	Server1 := upcloud.Server{
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

	Server2 := upcloud.Server{
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

	Server3 := upcloud.Server{
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

	Server4 := upcloud.Server{
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

	Server5 := upcloud.Server{
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

	allServers := &upcloud.Servers{
		Servers: []upcloud.Server{
			Server1,
			Server2,
			Server3,
			Server4,
			Server5,
		},
	}
	unambiguousServers := []upcloud.Server{
		Server2,
		Server3,
	}

	t.Run("resolve uuid", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServers").Return(allServers, nil)
		res := resolver.CachingServer{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, srv := range allServers.Servers {
			resolved := argResolver(srv.UUID)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, srv.UUID, value)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetServers", 1)
	})

	t.Run("resolve hostname", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServers").Return(allServers, nil)
		res := resolver.CachingServer{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, srv := range unambiguousServers {
			resolved := argResolver(srv.Hostname)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, srv.UUID, value)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetServers", 1)
	})

	t.Run("resolve title", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServers").Return(allServers, nil)
		res := resolver.CachingServer{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, srv := range unambiguousServers {
			resolved := argResolver(srv.Title)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, srv.UUID, value)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetServers", 1)
	})

	t.Run("resolve title glob", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServers").Return(allServers, nil)
		res := resolver.CachingServer{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		resolved := argResolver("server-*-title")
		value, err := resolved.GetAll()
		assert.NoError(t, err)
		assert.Len(t, value, 5)

		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetServers", 1)
	})

	t.Run("failure situations", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServers").Return(allServers, nil)

		res := resolver.CachingServer{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		// ambiguous hostname
		resolved := argResolver(Server4.Hostname)
		value, err := resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(Server4.Hostname))
		assert.Equal(t, "", value)

		// ambiguous title
		resolved = argResolver(Server1.Title)
		value, err = resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(Server1.Title))
		assert.Equal(t, "", value)

		// not found
		resolved = argResolver("notfound")
		value, err = resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("notfound"))
		assert.Equal(t, "", value)

		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetServers", 1)
	})
}

func TestFailingServerResolution(t *testing.T) {
	mService := &smock.Service{}
	var nilResponse *upcloud.Servers
	mService.On("GetServers").Return(nilResponse, errors.New("MOCKERROR"))
	res := resolver.CachingServer{}
	argResolver, err := res.Get(context.TODO(), mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
