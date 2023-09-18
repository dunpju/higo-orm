package him

import (
	"gorm.io/gorm"
)

func DbConfig(connect string) *DBConfig {
	return &DBConfig{connect: connect}
}

type DBConfig struct {
	connect       string
	username      string
	password      string
	host          string
	port          string
	database      string
	charset       string
	driver        string
	prefix        string
	maxIdle       int
	maxOpen       int
	maxLifetime   int
	logMode       string
	slowThreshold int
	colorful      bool
}

func (d *DBConfig) Connect() string {
	return d.connect
}

func (d *DBConfig) Username() string {
	return d.username
}

func (d *DBConfig) Password() string {
	return d.password
}

func (d *DBConfig) Host() string {
	return d.host
}

func (d *DBConfig) Port() string {
	return d.port
}

func (d *DBConfig) Database() string {
	return d.database
}

func (d *DBConfig) Charset() string {
	return d.charset
}

func (d *DBConfig) Driver() string {
	return d.driver
}

func (d *DBConfig) Prefix() string {
	return d.prefix
}

func (d *DBConfig) MaxIdle() int {
	return d.maxIdle
}

func (d *DBConfig) MaxOpen() int {
	return d.maxOpen
}

func (d *DBConfig) MaxLifetime() int {
	return d.maxLifetime
}

func (d *DBConfig) LogMode() string {
	return d.logMode
}

func (d *DBConfig) SlowThreshold() int {
	return d.slowThreshold
}

func (d *DBConfig) Colorful() bool {
	return d.colorful
}

func (d *DBConfig) SetUsername(username string) *DBConfig {
	d.username = username
	return d
}

func (d *DBConfig) SetPassword(password string) *DBConfig {
	d.password = password
	return d
}

func (d *DBConfig) SetHost(host string) *DBConfig {
	d.host = host
	return d
}

func (d *DBConfig) SetPort(port string) *DBConfig {
	d.port = port
	return d
}

func (d *DBConfig) SetDatabase(database string) *DBConfig {
	d.database = database
	return d
}

func (d *DBConfig) SetCharset(charset string) *DBConfig {
	d.charset = charset
	return d
}

func (d *DBConfig) SetDriver(driver string) *DBConfig {
	d.driver = driver
	return d
}

func (d *DBConfig) SetPrefix(prefix string) *DBConfig {
	d.prefix = prefix
	return d
}

func (d *DBConfig) SetMaxIdle(maxIdle int) *DBConfig {
	d.maxIdle = maxIdle
	return d
}

func (d *DBConfig) SetMaxOpen(maxOpen int) *DBConfig {
	d.maxOpen = maxOpen
	return d
}

func (d *DBConfig) SetMaxLifetime(maxLifetime int) *DBConfig {
	d.maxLifetime = maxLifetime
	return d
}

func (d *DBConfig) SetLogMode(logMode string) *DBConfig {
	d.logMode = logMode
	return d
}

func (d *DBConfig) SetSlowThreshold(slowThreshold int) *DBConfig {
	d.slowThreshold = slowThreshold
	return d
}

func (d *DBConfig) SetColorful(colorful bool) *DBConfig {
	d.colorful = colorful
	return d
}

func (d *DBConfig) Init() (*gorm.DB, error) {
	return Init(d)
}
