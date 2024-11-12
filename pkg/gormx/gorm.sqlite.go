package gormx

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func newGORMSqlite(cfg *SConfig) (*gorm.DB, error) {
	_ = os.MkdirAll(filepath.Dir(cfg.DSN), os.ModePerm)
	dialectal := sqlite.Open(cfg.DSN)

	c := &gorm.Config{
		Logger: logger.Discard,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: true,
		},
	}
	if cfg.Debug {
		c.Logger = logger.Default
	}

	db, err := gorm.Open(dialectal, c)
	if err != nil {
		return nil, err
	}

	if err = setResolver(cfg, db); err != nil {
		return nil, err
	}

	if cfg.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)

	return db, nil
}

func setResolverSqlite(resolver *dbresolver.DBResolver, r *SResolver) {
	resolverCfg := dbresolver.Config{}
	for _, replica := range r.Replicas {
		resolverCfg.Replicas = append(resolverCfg.Replicas, sqlite.Open(replica))
	}
	for _, source := range r.Sources {
		resolverCfg.Sources = append(resolverCfg.Sources, sqlite.Open(source))
	}
	tables := stringSliceToAnySlice(r.Tables)
	resolver.Register(resolverCfg, tables...)
	fmt.Printf(
		"Use resolver, #tables: %v, #replicas: %v, #sources: %v\n",
		tables,
		r.Replicas,
		r.Sources,
	)
}
