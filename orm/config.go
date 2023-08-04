package orm

import (
	"fmt"
	"github.com/dunpju/higo-config/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	dbc         *dbConfig
	dbOnce      sync.Once
	confDefault *config.Configure
)

func DbConfig() *dbConfig {
	return dbc
}

type dbConfig struct {
	Username    string
	Password    string
	Host        string
	Port        string
	Database    string
	Charset     string
	Driver      string
	Prefix      string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime int
	LogMode     bool
}

func newGorm() (*gorm.DB, error) {
	dbOnce.Do(func() {
		confDefault = config.Db("DB.Default").(*config.Configure)
		dbc = &dbConfig{
			Username:    confDefault.Get("Username").(string),
			Password:    confDefault.Get("Password").(string),
			Host:        confDefault.Get("Host").(string),
			Port:        confDefault.Get("Port").(string),
			Database:    confDefault.Get("Database").(string),
			Charset:     confDefault.Get("Charset").(string),
			Driver:      confDefault.Get("Driver").(string),
			Prefix:      confDefault.Get("Prefix").(string),
			MaxIdle:     confDefault.Get("MaxIdle").(int),
			MaxOpen:     confDefault.Get("MaxOpen").(int),
			MaxLifetime: confDefault.Get("MaxLifetime").(int),
			LogMode:     confDefault.Get("LogMode").(bool),
		}
	})
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbc.Username,
		dbc.Password,
		dbc.Host,
		dbc.Port,
		dbc.Database,
		dbc.Charset,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Logger
	db.LogMode(logMode)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxOpen)
	db.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)
	if registerCallbackCounter == 1 {
		if db.Callback().Query().Get("gorm:Query") == nil {
			db.Callback().Query().Before("gorm:Query").Register("Query", sqlReplace)
		}
		if db.Callback().RowQuery().Get("gorm:Query") == nil {
			db.Callback().RowQuery().Before("gorm:RowQuery").Register("RowQuery", sqlReplace)
		}
		if db.Callback().Create().Get("gorm:Create") == nil {
			db.Callback().Create().Before("gorm:Create").Register("Create", sqlReplace)
		}
		if db.Callback().Update().Get("gorm:Update") == nil {
			db.Callback().Update().Before("gorm:Update").Register("Update", sqlReplace)
		}
		if db.Callback().Delete().Get("gorm:Delete") == nil {
			db.Callback().Delete().Before("gorm:Delete").Register("Delete", sqlReplace)
		}
		logger.Logrus.Infoln(fmt.Sprintf("DB %s:%s Connection success!", dbConfig.Host,
			dbConfig.Port))
	}
	return db
}
