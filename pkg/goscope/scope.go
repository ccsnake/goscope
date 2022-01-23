package goscope

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed sqlite.sql
var SQLITE_CreateTable string

//go:embed mysql.sql
var MYSQL_CreateTable string

type Scope struct {
	DB     *sql.DB
	Config *Environment
}

func (s *Scope) GetConfig() *Environment {
	return s.Config
}

func (s *Scope) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}

func (s *Scope) setupDB() error {
	db, err := databaseSetup(databaseInformation{
		databaseType:          s.Config.GoScopeDatabaseType,
		connection:            s.Config.GoScopeDatabaseConnection,
		maxOpenConnections:    s.Config.GoScopeDatabaseMaxOpenConnections,
		maxIdleConnections:    s.Config.GoScopeDatabaseMaxIdleConnections,
		maxConnectionLifetime: s.Config.GoScopeDatabaseMaxConnLifetime,
	})
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	s.DB = db

	conn, err := db.Conn(context.TODO())
	if err != nil {
		_ = s.Close()
		return fmt.Errorf("scope:could not connect to database: %v", err)
	}
	defer conn.Close()

	if s.Config.GoScopeDatabaseType == "mysql" {
		_, err = conn.ExecContext(context.TODO(), MYSQL_CreateTable)

	} else if strings.HasPrefix(s.Config.GoScopeDatabaseType, "sqlite") {
		_, err = conn.ExecContext(context.TODO(), SQLITE_CreateTable)
	}
	if err != nil {
		return fmt.Errorf("scope:could not create table: %v", err)
	}

	return nil
}
