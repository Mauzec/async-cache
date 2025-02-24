package storage

import (
	"errors"
	"sync"
)

var ErrNotFoundByKey = errors.New("value not found")

type Cache interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

type CacheWithMetrics interface {
	Cache
	Metrics
}

type ToughAsyncCacheWithMetrics struct {
	st    map[string]string
	mu    sync.Mutex
	total int64
}

func NewCacheWithMetrics() CacheWithMetrics {
	return &ToughAsyncCacheWithMetrics{
		st: map[string]string{},
		mu: sync.Mutex{},
	}
}

func (c *ToughAsyncCacheWithMetrics) Get(key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.st[key]
	if !ok {
		return "", ErrNotFoundByKey
	}

	return v, nil
}

func (c *ToughAsyncCacheWithMetrics) Set(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.st[key] = value
	c.total++
	return nil
}

func (c *ToughAsyncCacheWithMetrics) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.st, key)
	c.total--
	return nil
}

func (c *ToughAsyncCacheWithMetrics) TotalAmount() int64 {
	return c.total
}

type ToughtAsyncCache struct {
	data map[string]string
	mu   sync.Mutex
}

func NewToughtAsyncCache() Cache {
	return &ToughtAsyncCache{
		data: make(map[string]string),
		mu:   sync.Mutex{},
	}
}

func (c *ToughtAsyncCache) Get(key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.data[key]
	if !ok {
		return "", ErrNotFoundByKey
	}
	return value, nil
}

func (c *ToughtAsyncCache) Set(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *ToughtAsyncCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

type AsyncCache struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewAsyncCache() Cache {
	return &AsyncCache{
		data: make(map[string]string),
		mu:   sync.RWMutex{},
	}
}

func (c *AsyncCache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	if !ok {
		return "", ErrNotFoundByKey
	}
	return value, nil
}

func (c *AsyncCache) Set(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *AsyncCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

type SyncCache struct {
	data map[string]string
}

func NewSyncCache() Cache {
	return &SyncCache{
		data: make(map[string]string),
	}
}

func (sc *SyncCache) Set(key, value string) error {
	sc.data[key] = value
	return nil
}

func (sc *SyncCache) Get(key string) (value string, err error) {
	value, ok := sc.data[key]
	if !ok {
		return "", ErrNotFoundByKey
	}
	return value, nil
}

func (sc *SyncCache) Delete(key string) error {
	_, ok := sc.data[key]
	if !ok {
		return ErrNotFoundByKey
	}
	delete(sc.data, key)
	return nil
}
