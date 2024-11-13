package jwtx

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const defaultKey = "eO5A5MMBceOJakUT6pKJ"

var ErrInvalidToken = errors.New("invalid token")

func GetSigningMethod(v string) (method jwt.SigningMethod) {
	switch v {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}
	return
}

type IAuth interface {
	// GenerateToken Generate a JWT with the provided subject
	GenerateToken(ctx context.Context, subject string) (IToken, error)
	// DestroyToken Invalidate a token by removing it from the token store
	DestroyToken(ctx context.Context, accessToken string) error
	// ParseSubject Parse the subject (or user identifier) from a given access token
	ParseSubject(ctx context.Context, accessToken string) (string, error)
	// Close any resources held by the SJWTAuth instance
	Close(ctx context.Context) error
}

func New(store IStore, opts ...FOption) IAuth {
	o := SOptions{
		signingMethod: jwt.SigningMethodHS512,
		signingKeyNew: []byte(defaultKey),
		SigningKeyOld: nil,
		keyFns:        nil,
		expired:       7200,
		tokenType:     "Bearer",
	}

	for _, opt := range opts {
		opt(&o)
	}

	o.keyFns = append(o.keyFns, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return o.signingKeyNew, nil
	})

	if o.SigningKeyOld != nil {
		o.keyFns = append(o.keyFns, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return o.SigningKeyOld, nil
		})
	}

	return &SJWTAuth{
		opts:  &o,
		store: store,
	}
}

type SJWTAuth struct {
	opts  *SOptions
	store IStore
}

func (j *SJWTAuth) callStore(fn func(IStore) error) error {
	if store := j.store; store != nil {
		return fn(store)
	}
	return nil
}

func (j *SJWTAuth) parseToken(t string) (*jwt.StandardClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	for _, keyFn := range j.opts.keyFns {
		if token, err := jwt.ParseWithClaims(t, &jwt.StandardClaims{}, keyFn); err != nil || token == nil || !token.Valid {
			continue
		}
		break
	}

	if err != nil || token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return token.Claims.(*jwt.StandardClaims), nil
}

func (j *SJWTAuth) GenerateToken(ctx context.Context, subject string) (IToken, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(j.opts.expired) * time.Second).Unix()
	tt := jwt.NewWithClaims(j.opts.signingMethod, jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   subject,
	})

	t, err := tt.SignedString([]byte(j.opts.signingKeyNew))
	if err != nil {
		return nil, err
	}

	ti := &SToken{
		AccessToken: t,
		TokenType:   j.opts.tokenType,
		ExpiresAt:   expiresAt,
	}
	return ti, nil
}

func (j *SJWTAuth) DestroyToken(ctx context.Context, t string) error {
	claims, err := j.parseToken(t)
	if err != nil {
		return err
	}

	return j.callStore(func(store IStore) error {
		expired := time.Until(time.Unix(claims.ExpiresAt, 0))
		return store.Set(ctx, t, expired)
	})
}

func (j *SJWTAuth) ParseSubject(ctx context.Context, t string) (string, error) {
	if t == "" {
		return "", ErrInvalidToken
	}

	claims, err := j.parseToken(t)
	if err != nil {
		return "", err
	}

	err = j.callStore(func(store IStore) error {
		if exists, err := store.Check(ctx, t); err != nil {
			return err
		} else if exists {
			return ErrInvalidToken
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}

func (j *SJWTAuth) Close(ctx context.Context) error {
	return j.callStore(func(store IStore) error {
		return store.Close(ctx)
	})
}
