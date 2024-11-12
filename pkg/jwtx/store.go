package jwtx

import (
	"context"
	"time"
)

type IStore interface {
	Check(ctx context.Context, token string) (bool, error)
	Set(ctx context.Context, token string, expiration time.Duration) error
	Delete(ctx context.Context, token string) error
	Close(ctx context.Context) error
}

type SOptionsStore struct {
	CacheNS string
}

type FOptionStore func(*SOptionsStore)

type SStoreImplements struct {
	opts *SOptionsStore
	c    ICache
}

func (s *SStoreImplements) Check(ctx context.Context, token string) (bool, error) {
	return s.c.Exists(ctx, s.opts.CacheNS, token)
}

func (s *SStoreImplements) Set(ctx context.Context, token string, expiration time.Duration) error {
	return s.c.Set(ctx, s.opts.CacheNS, token, "", expiration)
}

func (s *SStoreImplements) Delete(ctx context.Context, token string) error {
	return s.c.Delete(ctx, s.opts.CacheNS, token)
}

func (s *SStoreImplements) Close(ctx context.Context) error {
	return s.c.Close(ctx)
}

type ICache interface {
	Set(ctx context.Context, ns, key, value string, expiration time.Duration) error
	Get(ctx context.Context, ns, key string) (string, bool, error)
	Delete(ctx context.Context, ns, key string) error
	Exists(ctx context.Context, ns, key string) (bool, error)
	Close(ctx context.Context) error
}

func NewStoreWithCache(cache ICache, opts ...FOptionStore) IStore {
	s := &SStoreImplements{
		c: cache,
		opts: &SOptionsStore{
			CacheNS: "jwt",
		},
	}
	for _, opt := range opts {
		opt(s.opts)
	}
	return s
}
