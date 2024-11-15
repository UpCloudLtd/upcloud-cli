package resolver_test

import (
	"context"
	"errors"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"

	smock "github.com/UpCloudLtd/upcloud-cli/v3/internal/mock"
	"github.com/UpCloudLtd/upcloud-cli/v3/internal/resolver"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestServerGroupResolution(t *testing.T) {
	ServerGroup1 := upcloud.ServerGroup{
		Title: "serverGroup-1-title",
		UUID:  "1fdfda29-ead1-4855-b71f-1e33eb2ca9da",
	}

	ServerGroup2 := upcloud.ServerGroup{
		Title: "serverGroup-2-title",
		UUID:  "f77a5b25-84af-4f52-bc40-581930091fab",
	}

	ServerGroup3 := upcloud.ServerGroup{
		Title: "serverGroup-3-title",
		UUID:  "f0131b8f-ffe0-4271-83a8-c75b99e168cc",
	}

	ServerGroup4 := upcloud.ServerGroup{
		Title: ServerGroup1.Title,
		UUID:  "e5b3a855-cd8a-45b6-8cef-c7c860a0221d",
	}

	ServerGroup5 := upcloud.ServerGroup{
		Title: "serverGroup-5-title",
		UUID:  "39bc2725-213d-46c8-8b25-49990c6966af",
	}

	allServerGroups := upcloud.ServerGroups{
		ServerGroup1,
		ServerGroup2,
		ServerGroup3,
		ServerGroup4,
		ServerGroup5,
	}
	unambiguousServerGroups := []upcloud.ServerGroup{
		ServerGroup2,
		ServerGroup3,
	}

	t.Run("resolve uuid", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServerGroups", &request.GetServerGroupsRequest{}).Return(allServerGroups, nil)
		res := resolver.CachingServerGroup{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, srv := range allServerGroups {
			resolved := argResolver(srv.UUID)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, srv.UUID, value)
		}
		// make sure caching works, eg. we didn't call GetServerGroups more than once
		mService.AssertNumberOfCalls(t, "GetServerGroups", 1)
	})

	t.Run("resolve title", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServerGroups", &request.GetServerGroupsRequest{}).Return(allServerGroups, nil)
		res := resolver.CachingServerGroup{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, srv := range unambiguousServerGroups {
			resolved := argResolver(srv.Title)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, srv.UUID, value)
		}
		// make sure caching works, eg. we didn't call GetServerGroups more than once
		mService.AssertNumberOfCalls(t, "GetServerGroups", 1)
	})

	t.Run("failure situations", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetServerGroups", &request.GetServerGroupsRequest{}).Return(allServerGroups, nil)

		res := resolver.CachingServerGroup{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		// ambiguous title
		resolved := argResolver(ServerGroup1.Title)
		value, err := resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(ServerGroup1.Title))
		assert.Equal(t, "", value)

		// not found
		resolved = argResolver("notfound")
		value, err = resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.NotFoundError("notfound"))
		assert.Equal(t, "", value)

		// make sure caching works, eg. we didn't call GetServerGroups more than once
		mService.AssertNumberOfCalls(t, "GetServerGroups", 1)
	})
}

func TestFailingServerGroupResolution(t *testing.T) {
	mService := &smock.Service{}
	mService.On("GetServerGroups", &request.GetServerGroupsRequest{}).Return(upcloud.ServerGroups{}, errors.New("MOCKERROR"))
	res := resolver.CachingServerGroup{}
	argResolver, err := res.Get(context.TODO(), mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
