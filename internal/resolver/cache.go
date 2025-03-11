package resolver

import "fmt"

type Cache[T any] struct {
	cache map[string]T
}

func (c *Cache[T]) AddCached(uuid string, val T) {
	if c.cache == nil {
		c.cache = make(map[string]T)
	}
	c.cache[uuid] = val
}

func (c *Cache[T]) GetCached(uuid string) (T, error) {
	if c.cache == nil {
		return *new(T), CacheUninitializedError(fmt.Sprintf("%T", *new(T)))
	}

	if val, ok := c.cache[uuid]; ok {
		return val, nil
	}

	return *new(T), NotFoundError(uuid)
}
