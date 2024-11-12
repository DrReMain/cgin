package gormx

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type SConfig struct {
	TablePrefix  string
	Debug        bool
	DBType       string
	DSN          string
	MaxLifeTime  int
	MaxIdleTime  int
	MaxOpenConns int
	MaxIdleConns int

	Resolver []SResolver
}

type SResolver struct {
	Replicas []string
	Sources  []string
	Tables   []string
}

func NewGORM(cfg *SConfig) (*gorm.DB, error) {
	switch strings.ToLower(cfg.DBType) {
	case "mysql":
		return newGORMMysql(cfg)
	case "pgsql":
		return newGORMPgsql(cfg)
	case "sqlite":
		return newGORMSqlite(cfg)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}
}

func setResolver(cfg *SConfig, db *gorm.DB) error {
	if len(cfg.Resolver) > 0 {
		resolver := &dbresolver.DBResolver{}
		for _, r := range cfg.Resolver {
			switch strings.ToLower(cfg.DBType) {
			case "mysql":
				setResolverMysql(resolver, &r)
			case "pgsql":
				setResolverPgsql(resolver, &r)
			case "sqlite":
				setResolverSqlite(resolver, &r)
			default:
				continue
			}
		}

		resolver.
			SetMaxOpenConns(cfg.MaxOpenConns).
			SetMaxIdleConns(cfg.MaxIdleConns).
			SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second).
			SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)
		return db.Use(resolver)
	}
	return nil
}

func stringSliceToAnySlice(s []string) []any {
	r := make([]any, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}
