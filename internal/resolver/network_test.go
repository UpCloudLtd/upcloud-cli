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

var Network1 = upcloud.Network{
	Name: "network-1",
	UUID: "03e15cf5-8817-42ab-b017-970666be96ec",
	Type: upcloud.NetworkTypeUtility,
	Zone: "fi-hel1",
}

var Network2 = upcloud.Network{
	Name: "network-2",
	UUID: "03f5ad16-a63a-4670-8449-c01d1e97281e",
	Type: upcloud.NetworkTypePrivate,
	Zone: "fi-hel1",
}

var Network3 = upcloud.Network{
	Name: "network-3",
	UUID: "0357ce0a-eeb0-49fc-9f2c-a05c3ac57066",
	Type: upcloud.NetworkTypeUtility,
	Zone: "uk-lon1",
}

var Network4 = upcloud.Network{
	Name: Network1.Name,
	UUID: "03e49768-f13a-42c3-bea7-4e2471657f2f",
	Type: upcloud.NetworkTypePublic,
	Zone: "uk-lon1",
}

var (
	networks            = &upcloud.Networks{Networks: []upcloud.Network{Network1, Network2, Network3, Network4}}
	unambiguousNetworks = []upcloud.Network{Network2, Network3}
)

func TestNetworkResolution(t *testing.T) {
	t.Run("resolve uuid", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetNetworks").Return(networks, nil)
		res := resolver.CachingNetwork{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, network := range networks.Networks {
			resolved := argResolver(network.UUID)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, network.UUID, value)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetNetworks", 1)
	})

	t.Run("resolve uuid prefix", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetNetworks").Return(networks, nil)
		res := resolver.CachingNetwork{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		resolved := argResolver("035")
		value, err := resolved.GetOnly()
		assert.NoError(t, err)
		assert.Equal(t, Network3.UUID, value)
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetNetworks", 1)
	})

	t.Run("resolve name", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetNetworks").Return(networks, nil)
		res := resolver.CachingNetwork{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, network := range unambiguousNetworks {
			resolved := argResolver(network.Name)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, network.UUID, value)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetNetworks", 1)
	})

	t.Run("failure situations", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetNetworks").Return(networks, nil)
		res := resolver.CachingNetwork{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		// ambiguous name
		resolved := argResolver(Network1.Name)
		value, err := resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(Network1.Name))
		assert.Equal(t, "", value)

		// ambiguous UUID prefix
		resolved = argResolver("03")
		value, err = resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError("03"))
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
		mService.AssertNumberOfCalls(t, "GetNetworks", 1)
	})
}

func TestCachingNetwork_GetCached(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetNetworks").Return(networks, nil)
	res := resolver.CachingNetwork{}

	// should fail before cache initialized
	cached, err := res.GetCached(Network1.UUID)
	assert.Error(t, err)
	assert.Equal(t, upcloud.Network{}, cached)

	// get resolver to init the cache.. TODO: is this the best way?
	_, err = res.Get(context.TODO(), mService)
	assert.NoError(t, err)
	for _, network := range networks.Networks {
		cached, err := res.GetCached(network.UUID)
		assert.NoError(t, err)
		assert.Equal(t, network, cached)
	}

	// try not found
	cached, err = res.GetCached("dslkfjsdkfj")
	assert.ErrorIs(t, err, resolver.NotFoundError("dslkfjsdkfj"))
	assert.Equal(t, upcloud.Network{}, cached)
}

func TestFailingNetworkResolution(t *testing.T) {
	mService := &smock.Service{}
	var nilResponse *upcloud.Networks
	mService.On("GetNetworks").Return(nilResponse, errors.New("MOCKERROR"))
	res := resolver.CachingNetwork{}
	argResolver, err := res.Get(context.TODO(), mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
