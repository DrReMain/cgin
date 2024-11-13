package jwtx

import "github.com/DrReMain/cgin/pkg/encoding/json"

type IToken interface {
	GetAccessToken() string
	GetTokenType() string
	GetExpiresAt() int64
	EncodeToJSON() ([]byte, error)
}

type SToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (t *SToken) GetAccessToken() string {
	return t.AccessToken
}

func (t *SToken) GetTokenType() string {
	return t.TokenType
}

func (t *SToken) GetExpiresAt() int64 {
	return t.ExpiresAt
}

func (t *SToken) EncodeToJSON() ([]byte, error) {
	return json.Marshal(t)
}
