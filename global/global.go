package global

import (
	"template/global/config"
	"template/global/database"
	"template/global/logger"
)
var Config = config.LoadingConfigure()
var DB = database.CreateConnect(Config.Database)
var Logger = logger.NewLogger(Config.System)
