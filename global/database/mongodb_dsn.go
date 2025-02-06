package database

import "github.com/jingyuexing/go-utils"

type MongoDB struct {
	template string
	host     string
	config   string
	port     int
	dbname   string
	username string
	password string
}

func NewMongoDSN(username, password, host, database, config string, port int) *MongoDB {
	return &MongoDB{
		template: "mongodb://{username}:{password}@{host}/{database}{config}",
		username: username,
		password: password,
		host:     host,
		dbname:   database,
		port:     port,
		config:   config,
	}
}

// mongodb://{username}:{password}@{host}/{database}?config
func (mg *MongoDB) DSN() string {
	host := mg.host
	if mg.port > 0 {
		host = host + ":" + utils.ToString(mg.port)
	}
	config := ""
	if mg.config != "" {
		config = "?" + mg.config
	}
	return utils.Template(mg.template, map[string]any{
		"username": mg.username,
		"password": mg.password,
		"host":     host,
		"database": mg.dbname,
		"config":   config,
	}, "{}")
}
