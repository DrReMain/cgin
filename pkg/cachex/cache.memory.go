package cachex

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type SConfigMemory struct {
	CleanupInterval time.Duration
}

func NewCacheMemory(cfg *SConfigMemory, opts ...FOption) ICache {
	defaultOptions := &SOptions{Delimiter: defaultDelimiter}
	for _, o := range opts {
		o(defaultOptions)
	}
	return &SCacheMemory{
		opts:  defaultOptions,
		cache: cache.New(0, cfg.CleanupInterval),
	}
}

type SCacheMemory struct {
	opts  *SOptions
	cache *cache.Cache
}

func (c *SCacheMemory) key(ns, key string) string {
	return fmt.Sprintf("%s%s%s", ns, c.opts.Delimiter, key)
}

func (c *SCacheMemory) Set(ctx context.Context, ns, key, value string, expiration time.Duration) error {
	c.cache.Set(c.key(ns, key), value, expiration)
	return nil
}

func (c *SCacheMemory) Get(ctx context.Context, ns, key string) (string, bool, error) {
	if val, ok := c.cache.Get(c.key(ns, key)); ok {
		return val.(string), ok, nil
	}
	return "", false, nil
}

func (c *SCacheMemory) Delete(ctx context.Context, ns, key string) error {
	c.cache.Delete(c.key(ns, key))
	return nil
}

func (c *SCacheMemory) GetAndDelete(ctx context.Context, ns, key string) (string, bool, error) {
	value, ok, err := c.Get(ctx, ns, key)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}

	c.cache.Delete(c.key(ns, key))
	return value, true, nil
}

func (c *SCacheMemory) Exists(ctx context.Context, ns, key string) (bool, error) {
	_, ok := c.cache.Get(c.key(ns, key))
	return ok, nil
}

func (c *SCacheMemory) Close(ctx context.Context) error {
	c.cache.Flush()
	return nil
}

func (c *SCacheMemory) Iterator(ctx context.Context, ns string, fn func(ctx context.Context, key string, value string) bool) error {
	for k, v := range c.cache.Items() {
		if strings.HasPrefix(k, c.key(ns, "")) {
			if !fn(ctx, strings.TrimPrefix(k, c.key(ns, "")), v.Object.(string)) {
				break
			}
		}
	}
	return nil
}
