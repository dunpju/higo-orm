package orm

import (
	"sync"
)

var (
	dbc     *dbConfig
	dbcOnce sync.Once
)

func DbConfig() *dbConfig {
	dbcOnce.Do(func() {
		dbc = &dbConfig{}
	})
	return dbc
}

type dbConfig struct {
	username    string
	password    string
	host        string
	port        string
	database    string
	charset     string
	driver      string
	prefix      string
	maxIdle     int
	maxOpen     int
	maxLifetime int
	logMode     string
	colorful    bool
}

func (d *dbConfig) SetUsername(username string) *dbConfig {
	d.username = username
	return d
}

func (d *dbConfig) SetPassword(password string) *dbConfig {
	d.password = password
	return d
}

func (d *dbConfig) SetHost(host string) *dbConfig {
	d.host = host
	return d
}

func (d *dbConfig) SetPort(port string) *dbConfig {
	d.port = port
	return d
}

func (d *dbConfig) SetDatabase(database string) *dbConfig {
	d.database = database
	return d
}

func (d *dbConfig) SetCharset(charset string) *dbConfig {
	d.charset = charset
	return d
}

func (d *dbConfig) SetDriver(driver string) *dbConfig {
	d.driver = driver
	return d
}

func (d *dbConfig) SetPrefix(prefix string) *dbConfig {
	d.prefix = prefix
	return d
}

func (d *dbConfig) SetMaxIdle(maxIdle int) *dbConfig {
	d.maxIdle = maxIdle
	return d
}

func (d *dbConfig) SetMaxOpen(maxOpen int) *dbConfig {
	d.maxOpen = maxOpen
	return d
}

func (d *dbConfig) SetMaxLifetime(maxLifetime int) *dbConfig {
	d.maxLifetime = maxLifetime
	return d
}

func (d *dbConfig) SetLogMode(logMode string) *dbConfig {
	d.logMode = logMode
	return d
}

func (d *dbConfig) SetColorful(colorful bool) *dbConfig {
	d.colorful = colorful
	return d
}
