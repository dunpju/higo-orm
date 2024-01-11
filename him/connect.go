package him

import (
	"fmt"
	"sync"
)

const (
	DefaultConnect = "Default"
)

var (
	_connect sync.Map
)

func DBConnect(connection string) (*DB, error) {
	conn, err := getConnect(connection)
	if err != nil {
		return nil, err
	}
	return conn.DB(), nil
}

func GetConnect(connection string) (*connect, error) {
	return getConnect(connection)
}

func getConnect(connection string) (*connect, error) {
	if conn, ok := _connect.Load(connection); ok {
		return newConnect(conn.(*connect).Dbc(), newDB(conn.(*connect).DB().GormDB(), conn.(*connect).dbc.connect, conn.(*connect).dbc, false)), nil
	}
	return nil, fmt.Errorf("connect nonexistent")
}

type connect struct {
	dbc *DBConfig
	db  *DB
}

func newConnect(dbc *DBConfig, db *DB) *connect {
	return &connect{dbc: dbc, db: db.Rest()}
}

func (c connect) Dbc() *DBConfig {
	return c.dbc
}

func (c connect) DB() *DB {
	return c.db
}
