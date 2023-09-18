package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	_db     *gorm.DB
	_dbOnce sync.Once
)

func Gorm() (*gorm.DB, error) {
	var err error
	_dbOnce.Do(func() {
		_db, err = Init()
	})
	if err != nil {
		return nil, err
	}
	return _db, nil
}

func Init() (*gorm.DB, error) {
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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,                 // 慢 SQL 阈值
			LogLevel:      logger.LogLevel(level.code), // Log level
			Colorful:      dbc.colorful,                // 彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(dbc.maxIdle)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(dbc.maxOpen)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(dbc.maxLifetime) * time.Second)

	_dbOnce.Do(func() {
		_db = db
	})

	return _db, nil
}
