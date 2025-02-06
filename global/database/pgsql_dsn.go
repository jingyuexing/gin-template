package database

import (
	"net/url"
	"strings"

	"github.com/jingyuexing/go-utils"
)

type PgSQLDSN struct {
	template string
	host     string
	port     int
	dbname   string
	config   string
	path     string
	user     string
	password string
}

func WithPgSQLHost(host string) DSNOption[PgSQLDSN] {
	return func(dsn *PgSQLDSN) {
		if host == "" {
			return
		}
		dsn.host = host
	}
}
func WithPgSQLPort(port int) DSNOption[PgSQLDSN] {
	return func(dsn *PgSQLDSN) {
		if port == 0 {
			return
		}
		dsn.port = port
	}
}
func WithPgSQLUser(user string) DSNOption[PgSQLDSN] {
	return func(dsn *PgSQLDSN) {
		if user == "" {
			return
		}
		dsn.user = user
	}
}
func WithPgSQLPassword(password string) DSNOption[PgSQLDSN] {
	return func(dsn *PgSQLDSN) {
		if password == "" {
			return
		}
		dsn.password = password
	}
}
func WithPgSQLDSNDBName(dbname string) DSNOption[PgSQLDSN] {
	return func(dsn *PgSQLDSN) {
		if dbname == "" {
			return
		}
		dsn.dbname = dbname
	}
}

func WithPgSQLConfig(config string) DSNOption[PgSQLDSN]{
	return func(dsn *PgSQLDSN) {
		if config == "" {
			return
		}
		dsn.config = config
	}
}

// "host=" + p.Path + " user=" + p.Username + " password=" + p.Password + " dbname=" + p.Dbname + " port=" + p.Port + " " + p.Config
func newPgSQLDSN(opts ...DSNOption[PgSQLDSN]) *PgSQLDSN {
	dsn := &PgSQLDSN{
		template: "host={host} user={user} password={password} dbname={dbname} port={port} {config}",
		host:     "127.0.0.1",
		user:     "postgre",
		password: "postgre",
		dbname:   "postgre",
		port:     5432,
		config:   "sslmode=disable TimeZone=Asia/Shanghai client_encoding=UTF8",
	}
	for _, p := range opts {
		p(dsn)
	}
	return dsn
}

func (self *PgSQLDSN) DSN() string {
	if strings.Contains(self.config, "&") {
		c := make([]string, 0)
		value, err := url.ParseQuery(self.config)
		if err != nil {
			self.config = ""
		}

		for k, v := range value {
			klower := strings.ToLower(k)
			if klower == "charset" || klower == "client_encoding" || klower == "encoding" {
				c = append(c, "client_encoding="+v[0])
			}
			if klower == "loc" || klower == "timezone" {
				c = append(c, "TimeZone="+v[0])
			}
		}
		self.config = strings.Join(c, " ")
	}
	return utils.Template(self.template, map[string]any{
		"host":     self.host,
		"user":     self.user,
		"password": self.password,
		"dbname":   self.dbname,
		"port":     self.port,
		"config":   self.config,
	}, "{}")
}
