package cachex

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type SConfigRedis struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func NewCacheRedis(cfg *SConfigRedis, opts ...FOption) ICache {
	defaultOptions := &SOptions{Delimiter: defaultDelimiter}
	for _, o := range opts {
		o(defaultOptions)
	}
	return &SCacheRedis{
		opts: defaultOptions,
		cli: redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Username: cfg.Username,
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
	}
}

type SCacheRedis struct {
	opts *SOptions
	cli  *redis.Client
}

func (c *SCacheRedis) key(ns, key string) string {
	return fmt.Sprintf("%s%s%s", ns, c.opts.Delimiter, key)
}

func (c *SCacheRedis) Set(ctx context.Context, ns, key, value string, expiration time.Duration) error {
	return c.cli.Set(ctx, c.key(ns, key), value, expiration).Err()
}

func (c *SCacheRedis) Get(ctx context.Context, ns, key string) (string, bool, error) {
	cmd := c.cli.Get(ctx, c.key(ns, key))
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}
		return "", false, err
	}
	return cmd.Val(), true, nil
}

func (c *SCacheRedis) Delete(ctx context.Context, ns, key string) error {
	if b, err := c.Exists(ctx, ns, key); err != nil {
		return err
	} else if !b {
		return nil
	}

	cmd := c.cli.Del(ctx, c.key(ns, key))
	if err := cmd.Err(); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}

func (c *SCacheRedis) GetAndDelete(ctx context.Context, ns, key string) (string, bool, error) {
	value, ok, err := c.Get(ctx, ns, key)
	if err != nil {
		return "", false, err
	} else if !ok {
		return "", false, nil
	}

	cmd := c.cli.Del(ctx, c.key(ns, key))
	if err := cmd.Err(); err != nil && !errors.Is(err, redis.Nil) {
		return "", false, err
	}
	return value, true, nil
}

func (c *SCacheRedis) Exists(ctx context.Context, ns, key string) (bool, error) {
	cmd := c.cli.Exists(ctx, c.key(ns, key))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (c *SCacheRedis) Close(ctx context.Context) error {
	return c.cli.Close()
}

func (c *SCacheRedis) Iterator(ctx context.Context, ns string, fn func(ctx context.Context, key string, value string) bool) error {
	var cursor uint64
Loop:
	for {
		cmd := c.cli.Scan(ctx, cursor, c.key(ns, "*"), 100)
		if err := cmd.Err(); err != nil {
			return err
		}

		keys, cur, err := cmd.Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			cmd := c.cli.Get(ctx, key)
			if err := cmd.Err(); err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				return err
			}
			if next := fn(ctx, strings.TrimPrefix(key, c.key(ns, "")), cmd.Val()); !next {
				break Loop
			}
		}

		if cur == 0 {
			break
		}
		cursor = cur
	}
	return nil
}
