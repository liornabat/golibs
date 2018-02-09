package tracing

import (
	"time"

	c "github.com/liornabat/golibs/cache"
)

type spanCache struct {
	cache *c.LRU
}

func newCache() *spanCache {
	s := &spanCache{
		cache: c.NewLRUWithOptions(100000,
			&c.Options{
				TimeNow: func() time.Time { return time.Now() },
				TTL:     2 * time.Minute,
			}),
	}
	return s
}

func (s *spanCache) getSpan(key string) (*Span, bool) {
	v := s.cache.Get(key)
	if v != nil {
		span, ok := v.(*Span)
		if ok {
			s.cache.Delete(key)
			return span, true
		}
	}
	return nil, false
}

func (s *spanCache) putSpan(key string, value *Span) {
	s.cache.Put(key, value)

}
