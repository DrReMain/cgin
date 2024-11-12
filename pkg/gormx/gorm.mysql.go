package gormx

import (
	"database/sql"
	"fmt"
	"time"

	stdmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func newGORMMysql(cfg *SConfig) (*gorm.DB, error) {
	if err := createDatabase(cfg.DSN); err != nil {
		return nil, err
	}
	dialectal := mysql.Open(cfg.DSN)

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

func createDatabase(dsn string) error {
	cfg, err := stdmysql.ParseDSN(dsn)
	if err != nil {
		return err
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", cfg.User, cfg.Passwd, cfg.Addr))
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET = `utf8mb4`;", cfg.DBName)
	_, err = db.Exec(query)
	return err
}

func setResolverMysql(resolver *dbresolver.DBResolver, r *SResolver) {
	resolverCfg := dbresolver.Config{}
	for _, replica := range r.Replicas {
		resolverCfg.Replicas = append(resolverCfg.Replicas, mysql.Open(replica))
	}
	for _, source := range r.Sources {
		resolverCfg.Sources = append(resolverCfg.Sources, mysql.Open(source))
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
