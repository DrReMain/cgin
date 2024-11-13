package configx

import (
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/creasty/defaults"

	"github.com/DrReMain/cgin/pkg/encoding/json"
	"github.com/DrReMain/cgin/pkg/encoding/toml"
)

var once sync.Once
var supportExts = []string{".toml", ".json"}

func parse(v any, env string) error {
	ext := filepath.Ext(env)
	if ext == "" || !slices.Contains(supportExts, ext) {
		return nil
	}

	buf, err := os.ReadFile(env)
	if err != nil {
		return err
	}

	switch ext {
	case ".json":
		err = json.Unmarshal(buf, v)
	case ".toml":
		err = toml.Unmarshal(buf, v)
	}
	return err
}

func Load(v any, dir string, env string) error {
	if err := defaults.Set(v); err != nil {
		return err
	}

	fullName := filepath.Join(dir, env)
	info, err := os.Stat(fullName)
	if err != nil {
		return err
	}

	if info.IsDir() {
		err := filepath.WalkDir(fullName, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			return parse(v, path)
		})
		if err != nil {
			return err
		}
	}

	if err := parse(v, fullName); err != nil {
		return err
	}
	return nil
}

func MustLoad(v any, dir string, env string) {
	once.Do(func() {
		if err := Load(v, dir, env); err != nil {
			panic(err)
		}
	})
}
