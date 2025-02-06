package boot

import (
	"fmt"
	"template/global"
	"template/router"
)

func init() {
	app := router.Routers()
	app.Run(":"+fmt.Sprintf("%d", global.Config.Env.Port))
}