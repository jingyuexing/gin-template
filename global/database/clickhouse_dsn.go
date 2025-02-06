package database

import "github.com/jingyuexing/go-utils"

// clickhouse://gorm:gorm@localhost:9942/gorm?dial_timeout=10s&read_timeout=20s

type ClickHouse struct {
	template string
	host     string
	config   string
	port     int
	dbname   string
	username string
	password string
}

func NewClickHouse(username, password, host, database, config string, port int) *ClickHouse {
	return &ClickHouse{
		template: "clickhouse://{username}:{password}@{host}/{database}{config}",
		host:     host,
		username: username,
		password: password,
		dbname:   database,
		config:   config,
		port:     port,
	}
}

func (click *ClickHouse) DSN() string {
	host := click.host
	if click.port > 0 {
		host = host + ":" + utils.ToString(click.port)
	}
	config := ""
	if click.config != "" {
		config = "?" + click.config
	}
	return utils.Template(click.template, map[string]any{
		"username": click.username,
		"password": click.password,
		"host":     host,
		"database": click.dbname,
		"config":   config,
	}, "{}")
}
