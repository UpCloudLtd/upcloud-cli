package resolver_test

import (
	"context"
	"errors"
	"testing"

	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud"
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
			resolved, err := argResolver(db.UUID)
			assert.NoError(t, err)
			assert.Equal(t, db.UUID, resolved)
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
		resolved, err := argResolver(db.Name)
		assert.NoError(t, err)
		assert.Equal(t, db.UUID, resolved)
		// Make sure caching works, eg. we didn't call GetKubernetesClusters more than once
		mService.AssertNumberOfCalls(t, "GetKubernetesClusters", 1)
	})

	t.Run("Failures", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetKubernetesClusters", mock.Anything).Return(mockClusters, nil)

		res := resolver.CachingKubernetes{}
		argResolver, err := res.Get(context.TODO(), mService)
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
