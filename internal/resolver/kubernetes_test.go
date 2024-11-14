package resolver_test

import (
	"context"
	"errors"
	"testing"

	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockClusters = []upcloud.KubernetesCluster{
	{Name: "asd", UUID: "abcdef"},
	{Name: "asd", UUID: "abcghi"},
	{Name: "qwe", UUID: "jklmno"},
}

func TestKubernetesResolution(t *testing.T) {
	t.Run("UUID", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetKubernetesClusters", mock.Anything).Return(mockClusters, nil)
		res := resolver.CachingKubernetes{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, db := range mockClusters {
			resolved := argResolver(db.UUID)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, db.UUID, value)
		}

		// Make sure caching works, eg. we didn't call GetKubernetesClusters more than once
		mService.AssertNumberOfCalls(t, "GetKubernetesClusters", 1)
	})

	t.Run("Name", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetKubernetesClusters", mock.Anything).Return(mockClusters, nil)
		res := resolver.CachingKubernetes{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		db := mockClusters[2]
		resolved := argResolver(db.Name)
		value, err := resolved.GetOnly()
		assert.NoError(t, err)
		assert.Equal(t, db.UUID, value)
		// Make sure caching works, eg. we didn't call GetKubernetesClusters more than once
		mService.AssertNumberOfCalls(t, "GetKubernetesClusters", 1)
	})

	t.Run("Failures", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetKubernetesClusters", mock.Anything).Return(mockClusters, nil)

		res := resolver.CachingKubernetes{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		// Ambiguous Name
		resolved := argResolver("asd")
		value, err := resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError("asd"))
		assert.Equal(t, "", value)

		// Not found
		resolved = argResolver("not-found")
		value, err = resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("not-found"))
		assert.Equal(t, "", value)

		// Make sure caching works, eg. we didn't call GetKubernetesClusters more than once
		mService.AssertNumberOfCalls(t, "GetKubernetesClusters", 1)
	})
}

func TestFailingKubernetesResolution(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetKubernetesClusters", mock.Anything).Return(nil, errors.New("MOCKERROR"))
	res := resolver.CachingKubernetes{}
	argResolver, err := res.Get(context.TODO(), mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
