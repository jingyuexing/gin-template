package database

import (
	"fmt"
	"template/global/config"
	"time"
	_ "time/tzdata"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateConnect(config config.Database) *gorm.DB {
	var result *gorm.DB
	dsn := DSNFactory(
		config.Type,
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Config,
		config.DBName,
	)
	switch config.Type {
	case "mysql":
		db, err := gorm.Open(mysql.Open(dsn))
		if err != nil {
			fmt.Printf("mysql has error:%s\n", err.Error())
		}
		result = db
	case "pgsql":
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Printf("postgres has error:%s\n", err.Error())
		}
		result = db
	case "clickhouse":
		db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Printf("clickhouse has error:%s\n", err.Error())
		}
		result = db
	}
	db,err := result.DB()
	if err != nil {
		// log.Fatalf("failed to get *sql.DB object: %v", err)
	}

	// 设置数据库连接池参数
	db.SetMaxOpenConns(100)           // 设置最大打开连接数
	db.SetMaxIdleConns(10)            // 设置最大空闲连接数
	db.SetConnMaxLifetime(1 * time.Hour)  // 设置连接的最大生命周期
	db.SetConnMaxIdleTime(30 * time.Minute)  // 设置空闲连接的最大存活时间
	return result
}
