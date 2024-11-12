package toml

import (
	"github.com/BurntSushi/toml"
)

var (
	Marshal   = toml.Marshal
	Unmarshal = toml.Unmarshal
)

type Primitive = toml.Primitive

func MarshalToString(v any) (string, error) {
	b, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
