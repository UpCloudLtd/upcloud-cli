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

func TestIPAddressResolution(t *testing.T) {
	ipAddress1 := upcloud.IPAddress{
		Address:    "94.237.117.151",
		Access:     "public",
		Family:     "IPv4",
		PartOfPlan: upcloud.FromBool(true),
		PTRRecord:  "94-237-117-151.fi-hel1.upcloud.host",
		ServerUUID: "005ab220-7ff6-42c9-8615-e4c02eb4104b",
		MAC:        "ee:1b:db:ca:6b:80",
		Floating:   upcloud.FromBool(false),
		Zone:       "fi-hel1",
	}
	ipAddress2 := upcloud.IPAddress{
		Address:    "94.237.117.152",
		Access:     "public",
		Family:     "IPv4",
		PartOfPlan: upcloud.FromBool(true),
		PTRRecord:  "94-237-117-152.fi-hel1.upcloud.host",
		ServerUUID: "005ab220-7ff6-42c9-8615-e4c02eb4104c",
		MAC:        "ee:1b:db:ca:6b:81",
		Floating:   upcloud.FromBool(false),
		Zone:       "fi-hel1",
	}
	ipAddress3 := upcloud.IPAddress{
		Address:    "94.237.117.153",
		Access:     "public",
		Family:     "IPv4",
		PartOfPlan: upcloud.FromBool(true),
		PTRRecord:  "94-237-117-153.fi-hel1.upcloud.host",
		ServerUUID: "005ab220-7ff6-42c9-8615-e4c02eb4104d",
		MAC:        "ee:1b:db:ca:6b:82",
		Floating:   upcloud.FromBool(false),
		Zone:       "fi-hel1",
	}
	ipAddress4 := upcloud.IPAddress{
		Address:    "94.237.117.154",
		Access:     "public",
		Family:     "IPv4",
		PartOfPlan: upcloud.FromBool(true),
		PTRRecord:  "94-237-117-154.fi-hel1.upcloud.host",
		ServerUUID: "005ab220-7ff6-42c9-8615-e4c02eb4104e",
		MAC:        "ee:1b:db:ca:6b:83",
		Floating:   upcloud.FromBool(false),
		Zone:       "fi-hel1",
	}
	ipAddress5 := upcloud.IPAddress{
		Address:    "94.237.117.156",
		Access:     "public",
		Family:     "IPv4",
		PartOfPlan: upcloud.FromBool(true),
		PTRRecord:  "94-237-117-154.fi-hel1.upcloud.host", // same PTR as 4
		ServerUUID: "005ab220-7ff6-42c9-8615-e4c02eb4104e",
		MAC:        "ee:1b:db:ca:6b:85",
		Floating:   upcloud.FromBool(false),
		Zone:       "fi-hel1",
	}

	addresses := &upcloud.IPAddresses{IPAddresses: []upcloud.IPAddress{
		ipAddress1, ipAddress2, ipAddress3, ipAddress4, ipAddress5,
	}}
	unambiguousAddresses := []upcloud.IPAddress{ipAddress1, ipAddress2, ipAddress3}

	t.Run("resolve address", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetIPAddresses").Return(addresses, nil)
		res := resolver.CachingIPAddress{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, network := range unambiguousAddresses {
			resolved := argResolver(network.Address)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, network.Address, value)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetIPAddresses", 1)
	})

	t.Run("resolve ptr records", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetIPAddresses").Return(addresses, nil)
		res := resolver.CachingIPAddress{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)
		for _, network := range unambiguousAddresses {
			resolved := argResolver(network.PTRRecord)
			value, err := resolved.GetOnly()
			assert.NoError(t, err)
			assert.Equal(t, network.Address, value)
		}
		// make sure caching works, eg. we didn't call GetServers more than once
		mService.AssertNumberOfCalls(t, "GetIPAddresses", 1)
	})

	t.Run("failure situations", func(t *testing.T) {
		mService := &smock.Service{}
		mService.On("GetIPAddresses").Return(addresses, nil)
		res := resolver.CachingIPAddress{}
		argResolver, err := res.Get(context.TODO(), mService)
		assert.NoError(t, err)

		// ambiguous ptr record
		resolved := argResolver(ipAddress4.PTRRecord)
		value, err := resolved.GetOnly()
		if !assert.Error(t, err) {
			t.FailNow()
		}
		assert.ErrorIs(t, err, resolver.AmbiguousResolutionError(ipAddress4.PTRRecord))
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
		mService.AssertNumberOfCalls(t, "GetIPAddresses", 1)
	})
}

func TestFailingIPAddressResolution(t *testing.T) {
	mService := &smock.Service{}
	var nilResponse *upcloud.IPAddresses
	mService.On("GetIPAddresses").Return(nilResponse, errors.New("MOCKERROR"))
	res := resolver.CachingIPAddress{}
	argResolver, err := res.Get(context.TODO(), mService)
	assert.EqualError(t, err, "MOCKERROR")
	assert.Nil(t, argResolver)
}
