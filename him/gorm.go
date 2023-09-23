package him

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func Gorm(connect string) (*gorm.DB, error) {
	connection, err := getConnect(connect)
	return connection.db.GormDB(), err
}

func Default() (*gorm.DB, error) {
	return Gorm(DefaultConnect)
}

func Init(dbc *DBConfig) (*gorm.DB, error) {
	if dbc.connect == "" {
		return nil, fmt.Errorf("connect cannot be empty")
	}
	if conn, ok := _connect.Load(dbc.connect); ok {
		return conn.(*connect).DB().GormDB(), nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbc.username,
		dbc.password,
		dbc.host,
		dbc.port,
		dbc.database,
		dbc.charset,
	)

	level, err := LogLevel(dbc.logMode)
	if err != nil {
		return nil, err
	}

	slowThreshold := 3
	if dbc.slowThreshold > 0 {
		slowThreshold = dbc.slowThreshold
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second * time.Duration(slowThreshold), // 慢 SQL 阈值
			LogLevel:      logger.LogLevel(level.code),                // Log level
			Colorful:      dbc.colorful,                               // 彩色打印
		},
	)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(dbc.maxIdle)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(dbc.maxOpen)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(dbc.maxLifetime) * time.Second)

	_connect.Store(dbc.connect, newConnect(dbc, newDB(gormDB, dbc.connect)))

	return gormDB, nil
}
