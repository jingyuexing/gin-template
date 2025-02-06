package database

type DSNOption[T any] func(dsn *T)

type DSN interface {
	DSN() string
}

func DSNFactory(DBType string, user string, password string, host string, port int, config string, dbname string) string {
	result := ""
	switch DBType {
	case "mysql":
		result = newMySQLDSN(
			WithMySQLHost(host),
			WithMySQLPort(port),
			WithMySQLPassword(password),
			WithMySQLUsername(user),
			WithMySQLDBName(dbname),
			WithMySQLDBConfig(config),
		// host, port, dbname, user, password, config,
		).DSN()
	case "pgsql":
		result = newPgSQLDSN(
			WithPgSQLDSNDBName(dbname),
			WithPgSQLPassword(password),
			WithPgSQLUser(user),
			WithPgSQLHost(host),
			WithPgSQLPort(port),
			WithPgSQLConfig(config),
		).DSN()
	case "oracle":
		result = newOracleDSN(user, password, host, port, dbname, config).DSN()
	case "mongodb":
		result = NewMongoDSN(user, password, host, dbname, config, port).DSN()
	case "clickhouse":
		result = NewClickHouse(user, password, host, dbname, config, port).DSN()
	case "redis":
		result = newRedisDSN(
			WithRedisHost(host),
			WithRedisPassword(password),
			WithRedisPort(port),
			WithRedisDB(dbname),
			WithRedisConfig(config),
		).DSN()
	}

	return result
}
