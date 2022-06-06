package cache

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	ErrNotFound = errors.New("not found")
)

type Cache struct {
	client *memcache.Client
}

func NewCache(hosts []string) *Cache {
	return &Cache{
		client: memcache.New(hosts...),
	}
}

func (c *Cache) Set(key string, value []byte) error {
	return c.client.Set(&memcache.Item{
		Key:   key,
		Value: value,
	})
}

func (c *Cache) Get(key string) ([]byte, error) {
	item, err := c.client.Get(key)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return nil, ErrNotFound
	}

	return item.Value, nil
}

func (c *Cache) Delete(key string) error {
	return c.client.Delete(key)
}
