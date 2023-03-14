package resolver_test

import (
	"context"
	"errors"
	"testing"

	smock "github.com/UpCloudLtd/upcloud-cli/v2/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v2/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/stretchr/testify/assert"
)

var Network1 = upcloud.Network{
	Name: "network-1",
	UUID: "28e15cf5-8817-42ab-b017-970666be96ec",
	Type: upcloud.NetworkTypeUtility,
	Zone: "fi-hel1",
}

var Network2 = upcloud.Network{
	Name: "network-2",
	UUID: "f9f5ad16-a63a-4670-8449-c01d1e97281e",
	Type: upcloud.NetworkTypePrivate,
	Zone: "fi-hel1",
}

var Network3 = upcloud.Network{
	Name: "network-3",
	UUID: "e157ce0a-eeb0-49fc-9f2c-a05c3ac57066",
	Type: upcloud.NetworkTypeUtility,
	Zone: "uk-lon1",
}

var Network4 = upcloud.Network{
	Name: Network1.Name,
	UUID: "b3e49768-f13a-42c3-bea7-4e2471657f2f",
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
			resolved, err := argResolver(network.UUID)
			assert.NoError(t, err)
			assert.Equal(t, network.UUID, resolved)
		}
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
			resolved, err := argResolver(network.Name)
			assert.NoError(t, err)
			assert.Equal(t, network.UUID, resolved)
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

		// ambigous name
		resolved, err := argResolver(Network1.Name)
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(Network1.Name))
		assert.Equal(t, "", resolved)

		// not found
		resolved, err = argResolver("notfounf")
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("notfounf"))
		assert.Equal(t, "", resolved)

		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetNetworks", 1)
	})
}

func TestCachingNetwork_GetCached(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetNetworks").Return(networks, nil)
	res := resolver.CachingNetwork{}

	// should fail before cahe initialized
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
