package database

import "github.com/jingyuexing/go-utils"

type MySQL struct {
	template string
	host     string
	config   string
	port     int
	dbname   string
	username string
	password string
}

func WithMySQLHost(host string) DSNOption[MySQL] {
	return func(dsn *MySQL) {
		if host == "" {
			return
		}
		dsn.host = host
	}
}
func WithMySQLPort(port int) DSNOption[MySQL] {
	return func(dsn *MySQL) {
		if port == 0 {
			return
		}
		dsn.port = port
	}
}
func WithMySQLUsername(user string) DSNOption[MySQL] {
	return func(dsn *MySQL) {
		if user == "" {
			return
		}
		dsn.username = user
	}
}
func WithMySQLPassword(password string) DSNOption[MySQL] {
	return func(dsn *MySQL) {
		if password == "" {
			return
		}
		dsn.password = password
	}
}
func WithMySQLDBName(dbname string) DSNOption[MySQL] {
	return func(dsn *MySQL) {
		if dbname == "" {
			return
		}
		dsn.dbname = dbname
	}
}

func WithMySQLDBConfig(config string) DSNOption[MySQL] {
	return func(dsn *MySQL) {
		if config == "" {
			return
		}
		dsn.config = config
	}
}


func newMySQLDSN(opts ...DSNOption[MySQL]) *MySQL {
	dsn := &MySQL{
		template: "{username}:{password}@tcp({path}:{port})/{dbname}?{config}",
		host:     "127.0.0.1",
		port:     3306,
		dbname:   "admin",
		username: "admin",
		password: "admin",
		config:   "charset=utf8mb4&parseTime=True&loc=Local",
	}
	for _, p := range opts {
		p(dsn)
	}
	return dsn
}

func (self *MySQL) DSN() string {
	return utils.Template(self.template, map[string]any{
		"username": self.username,
		"password": self.password,
		"path":     self.host,
		"port":     self.port,
		"dbname":   self.dbname,
		"config":   self.config,
	}, "{}")
}
