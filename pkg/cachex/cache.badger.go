package cachex

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type SConfigBadger struct {
	Path string
}

func NewCacheBadger(cfg *SConfigBadger, opts ...FOption) ICache {
	defaultOptions := &SOptions{Delimiter: defaultDelimiter}
	for _, o := range opts {
		o(defaultOptions)
	}

	badgerOpts := badger.DefaultOptions(cfg.Path)
	badgerOpts = badgerOpts.WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(badgerOpts)
	if err != nil {
		panic(err)
	}
	return &SCacheBadger{
		opts: defaultOptions,
		db:   db,
	}
}

type SCacheBadger struct {
	opts *SOptions
	db   *badger.DB
}

func (c *SCacheBadger) key(ns, key string) string {
	return fmt.Sprintf("%s%s%s", ns, c.opts.Delimiter, key)
}

func (c *SCacheBadger) Set(ctx context.Context, ns, key, value string, expiration time.Duration) error {
	return c.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(c.key(ns, key)), []byte(value))
		if expiration > 0 {
			entry = entry.WithTTL(expiration)
		}
		return txn.SetEntry(entry)
	})
}

func (c *SCacheBadger) Get(ctx context.Context, ns, key string) (string, bool, error) {
	var value string
	var ok bool

	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(c.key(ns, key)))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		value = string(val)
		ok = true
		return nil
	})
	return value, ok, err
}

func (c *SCacheBadger) Delete(ctx context.Context, ns, key string) error {
	b, err := c.Exists(ctx, ns, key)
	if err != nil {
		return err
	}
	if !b {
		return nil
	}
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(c.key(ns, key)))
	})
}

func (c *SCacheBadger) GetAndDelete(ctx context.Context, ns, key string) (string, bool, error) {
	value, ok, err := c.Get(ctx, ns, key)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	if err = c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(c.key(ns, key)))
	}); err != nil {
		return "", false, err
	}
	return value, true, nil
}

func (c *SCacheBadger) Exists(ctx context.Context, ns, key string) (bool, error) {
	var exists bool
	return exists, c.db.View(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(c.key(ns, key))); err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		exists = true
		return nil
	})
}

func (c *SCacheBadger) Close(ctx context.Context) error {
	return c.db.Close()
}

func (c *SCacheBadger) Iterator(ctx context.Context, ns string, fn func(ctx context.Context, key string, value string) bool) error {
	return c.db.View(func(txn *badger.Txn) error {
		iterOpts := badger.DefaultIteratorOptions
		iterOpts.Prefix = []byte(c.key(ns, ""))
		it := txn.NewIterator(iterOpts)
		defer it.Close()

		it.Rewind()
		for it.Valid() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			key := string(item.Key())
			if !fn(ctx, strings.TrimPrefix(key, c.key(ns, "")), string(val)) {
				break
			}
			it.Next()
		}
		return nil
	})
}
