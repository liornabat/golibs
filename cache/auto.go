package cache

import (
	"sync"
	"time"
)

type AutoCache struct {
	ttl time.Duration
	m   sync.Map
}

func NewAutoCache(ttl time.Duration) *AutoCache {
	ac := &AutoCache{
		ttl: ttl,
	}
	go ac.runCleanCache()
	return ac
}

func (ac *AutoCache) Put(key string, value interface{}) interface{} {
	ac.m.Store(key, &cacheEntry{
		key:        key,
		value:      value,
		expiration: time.Now().Add(ac.ttl),
	})
	return value
}

func (ac *AutoCache) Get(key string) interface{} {
	v, ok := ac.m.Load(key)
	if ok {
		return v.(*cacheEntry).value
	}
	return nil
}

func (ac *AutoCache) Exist(key string) bool {
	_, ok := ac.m.Load(key)
	return ok
}

func (ac *AutoCache) Delete(key string) {
	ac.m.Delete(key)
}

func (ac *AutoCache) Size() int {
	c := 0
	ac.m.Range(func(key, val interface{}) bool {
		c++
		return true
	})
	return c
}

func (ac *AutoCache) CompareAndSwap(key string, old, new interface{}) (interface{}, bool) {
	ac.Put(key, new)
	return new, true
}
func (ac *AutoCache) clean() {
	keys := []string{}
	ac.m.Range(func(key, val interface{}) bool {
		k, v := key.(string), val.(*cacheEntry)
		if time.Now().After(v.expiration) {
			keys = append(keys, k)
		}
		return true
	})
	for _, key := range keys {
		ac.m.Delete(key)
	}
}

func (ac *AutoCache) runCleanCache() {
	for {
		time.Sleep(time.Second)
		ac.clean()
	}
}
