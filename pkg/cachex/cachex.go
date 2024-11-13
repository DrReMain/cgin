package cachex

import (
	"context"
	"time"
)

var defaultDelimiter = ":"

type ICache interface {
	Set(ctx context.Context, ns, key, value string, expiration time.Duration) error
	Get(ctx context.Context, ns, key string) (string, bool, error)
	Delete(ctx context.Context, ns, key string) error
	GetAndDelete(ctx context.Context, ns, key string) (string, bool, error)
	Exists(ctx context.Context, ns, key string) (bool, error)
	Close(ctx context.Context) error
	Iterator(ctx context.Context, ns string, fn func(ctx context.Context, key, value string) bool) error
}

type SOptions struct {
	Delimiter string
}

type FOption func(o *SOptions)

func WithDelimiter(delimiter string) FOption {
	return func(o *SOptions) {
		o.Delimiter = delimiter
	}
}
