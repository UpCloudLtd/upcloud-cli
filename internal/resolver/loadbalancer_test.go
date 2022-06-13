package resolver_test

import (
	"errors"
	"testing"

	smock "github.com/UpCloudLtd/upcloud-cli/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockLoadBalancers = []upcloud.LoadBalancer{
	{Name: "asd", UUID: "abcdef"},
	{Name: "asd", UUID: "abcghi"},
	{Name: "qwe", UUID: "jklmno"},
}

func TestLoadBalancerResolution(t *testing.T) {
	t.Run("UUID", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetLoadBalancers", mock.Anything).Return(mockLoadBalancers, nil)
		res := resolver.CachingLoadBalancer{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)
		for _, db := range mockLoadBalancers {
			resolved, err := argResolver(db.UUID)
			assert.NoError(t, err)
			assert.Equal(t, db.UUID, resolved)
		}

		// Make sure caching works, eg. we didn't call GetLoadBalancers more than once
		mService.AssertNumberOfCalls(t, "GetLoadBalancers", 1)
	})

	t.Run("Name", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetLoadBalancers", mock.Anything).Return(mockLoadBalancers, nil)
		res := resolver.CachingLoadBalancer{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)

		db := mockLoadBalancers[2]
		resolved, err := argResolver(db.Name)
		assert.NoError(t, err)
		assert.Equal(t, db.UUID, resolved)
		// Make sure caching works, eg. we didn't call GetLoadBalancers more than once
		mService.AssertNumberOfCalls(t, "GetLoadBalancers", 1)
	})

	t.Run("Failures", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetLoadBalancers", mock.Anything).Return(mockLoadBalancers, nil)

		res := resolver.CachingLoadBalancer{}
		argResolver, err := res.Get(mService)
		assert.NoError(t, err)
		var resolved string

		// Ambigous Name
		resolved, err = argResolver("asd")
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError("asd"))
		assert.Equal(t, "", resolved)

		// Not found
		resolved, err = argResolver("not-found")
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("not-found"))
		assert.Equal(t, "", resolved)

		// Make sure caching works, eg. we didn't call GetLoadBalancers more than once
		mService.AssertNumberOfCalls(t, "GetLoadBalancers", 1)
	})
}

func TestFailingLoadBalancerResolution(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetLoadBalancers", mock.Anything).Return(nil, errors.New("MOCKERROR"))
	res := resolver.CachingLoadBalancer{}
	argResolver, err := res.Get(mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
