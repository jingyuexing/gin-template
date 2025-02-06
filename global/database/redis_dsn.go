package database

import "github.com/jingyuexing/go-utils"

// Redis 配置结构体
type Redis struct {
	template string
	host     string
	port     int
	password string // 可选密码
	db       string // 默认是 0
	config   string // 可选配置
}

// 设置 Redis 密码
func WithRedisPassword(password string) DSNOption[Redis] {
	return func(r *Redis) {
		if password == "" {
			return
		}
		r.password = password
	}
}

// 设置 Redis 密码
func WithRedisPort(port int) DSNOption[Redis] {
	return func(r *Redis) {
		if port == 0 {
			return
		}
		r.port = port
	}
}

// 设置 Redis 数据库
func WithRedisDB(db string) DSNOption[Redis] {
	return func(r *Redis) {
		if db == "0" {
			return
		}
		r.db = db
	}
}

func WithRedisHost(host string) DSNOption[Redis] {
	return func(dsn *Redis) {
		if host == "" {
			return
		}
		dsn.host = host
	}
}

// 设置 Redis 配置
func WithRedisConfig(config string) DSNOption[Redis] {
	return func(r *Redis) {
		if config == "" {
			return
		}
		r.config = config
	}
}

// newRedisDSN 创建 Redis 配置的构造函数，使用可选参数设置配置项
func newRedisDSN(opts ...DSNOption[Redis]) *Redis {
	// 使用默认配置
	redis := &Redis{
		template: "redis://{host}:{port}/{db}?{config}",
		host:     "127.0.0.1", // 默认 host
		port:     6379,        // 默认 port
		password: "",          // 默认无密码
		db:       "0",         // 默认数据库为 0
		config:   "simple",    // 默认无配置
	}

	// 通过传入的选项修改默认配置
	for _, opt := range opts {
		opt(redis)
	}
	if redis.config == "simple" {
		redis.template = "{host}:{port}"
		return redis
	}
	if redis.password != "" {
		redis.template = "redis://{password}@{host}:{port}/{db}?{config}"
	}
	return redis
}

func (self *Redis) DSN() string {
	return utils.Template(self.template, map[string]any{
		"password": self.password,
		"host":     self.host,
		"port":     self.port,
		"dbname":   self.db,
		"config":   self.config,
	}, "{}")
}
