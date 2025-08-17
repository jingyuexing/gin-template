package global

import (
	"sync"
	"template/global/config"
	"template/global/database"
	"template/global/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Once sync.Once
var Config *config.Configure
var DB *gorm.DB
var Logger *zap.Logger

func init() {
	Once.Do(func() {
		DB = database.CreateConnect(Config.Database)
		Logger = logger.NewLogger(Config.System)
		Config = config.LoadingConfigure()
	})
}
