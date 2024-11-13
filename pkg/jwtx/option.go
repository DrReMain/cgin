package jwtx

import (
	"github.com/golang-jwt/jwt"
)

type SOptions struct {
	signingMethod jwt.SigningMethod
	signingKeyNew []byte
	SigningKeyOld []byte
	keyFns        []func(*jwt.Token) (any, error)
	expired       int
	tokenType     string
}

type FOption func(*SOptions)

func SetExpired(expired int) FOption {
	return func(o *SOptions) {
		o.expired = expired
	}
}

func SetSigningKey(key, old string) FOption {
	return func(o *SOptions) {
		o.signingKeyNew = []byte(key)
		if old != "" && key != old {
			o.SigningKeyOld = []byte(old)
		}
	}
}

func SetSigningMethod(method jwt.SigningMethod) FOption {
	return func(o *SOptions) {
		o.signingMethod = method
	}
}
