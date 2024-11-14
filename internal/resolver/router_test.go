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

var Router1 = upcloud.Router{
	Name: "name-1",
	Type: "normal",
	UUID: "ffd3ab80-fe95-49c0-ab70-fbc987246c7a",
}

var Router2 = upcloud.Router{
	Name: "name-2",
	Type: "normal",
	UUID: "f14dd3e7-3dbb-4e3c-92b9-d1cf5178a13e",
}

var Router3 = upcloud.Router{
	Name: "name-3",
	Type: "normal",
	UUID: "ffd3ab80-fe95-49c0-ab70-fbc987246c99",
}

var Router4 = upcloud.Router{
	Name: "name-1",
	Type: "normal",
	UUID: "ffd3ab80-fe95-49c0-ab70-fbc987246c7b",
}

var allRouters = &upcloud.Routers{
	Routers: []upcloud.Router{
		Router1,
		Router2,
		Router3,
		Router4,
	},
}

var unambiguousRouters = []upcloud.Router{
	Router2,
	Router3,
}

func TestRouterResolution(t *testing.T) {
	t.Run("resolve uuid", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetRouters").Return(allRouters, nil)
		res := resolver.CachingRouter{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, router := range allRouters.Routers {
			resolved := argResolver(router.UUID)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, router.UUID, value)
		}
		// make sure caching works, eg. we didn't call GetRouters more than once
		mService.AssertNumberOfCalls(t, "GetRouters", 1)
	})

	t.Run("resolve hostname", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetRouters").Return(allRouters, nil)
		res := resolver.CachingRouter{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, srv := range unambiguousRouters {
			resolved := argResolver(srv.Name)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, srv.UUID, value)
		}
		// make sure caching works, eg. we didn't call GetRouters more than once
		mService.AssertNumberOfCalls(t, "GetRouters", 1)
	})

	t.Run("failure situations", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetRouters").Return(allRouters, nil)

		res := resolver.CachingRouter{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		// ambiguous name
		resolved := argResolver(Router1.Name)
		value, err := resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(Router1.Name))
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
		mService.AssertNumberOfCalls(t, "GetRouters", 1)
	})
}

func TestFailingRouterResolution(t *testing.T) {
	mService := &smock.Service{}
	var nilResponse *upcloud.Routers
	mService.On("GetRouters").Return(nilResponse, errors.New("MOCKERROR"))
	res := resolver.CachingRouter{}
	argResolver, err := res.Get(context.TODO(), mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}

func TestCachingRouter_GetCached(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetRouters").Return(allRouters, nil)
	res := resolver.CachingRouter{}

	// should fail before cache initialized
	cached, err := res.GetCached(Router1.UUID)
	assert.Error(t, err)
	assert.Equal(t, upcloud.Router{}, cached)

	// get resolver to init the cache.. TODO: is this the best way?
	_, err = res.Get(context.TODO(), mService)
	assert.NoError(t, err)
	for _, router := range allRouters.Routers {
		cached, err := res.GetCached(router.UUID)
		assert.NoError(t, err)
		assert.Equal(t, router, cached)
	}

	// try not found
	cached, err = res.GetCached("dslkfjsdkfj")
	assert.ErrorIs(t, err, resolver.NotFoundError("dslkfjsdkfj"))
	assert.Equal(t, upcloud.Router{}, cached)
}
