package database

import "github.com/jingyuexing/go-utils"

type Oracle struct {
	template string
	path     string
	port     int
	config   string
	dbname   string
	username string
	password string
}

func newOracleDSN(username string, password string, path string, port int, dbname string, config string) *Oracle {
	return &Oracle{
		template: "oracle://{username}:{password}@{path}:{port}/{dbname}?{config}",
		username: username,
		password: password,
		path:     path,
		port:     port,
		dbname:   dbname,
		config:   config,
	}
}

func (self *Oracle) DSN() string {
	return utils.Template(self.template, map[string]any{
		"username": self.username,
		"password": self.password,
		"path":     self.path,
		"port":     self.port,
		"dbname":   self.dbname,
		"config":   self.config,
	}, "{}")
}
