package sqlinit

import (
	"beeapi/models"
	"encoding/gob"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/gomodule/redigo/redis"
)

func init() {
	//启用Session
	beego.BConfig.WebConfig.Session.SessionOn = true
	//Session的引擎
	session_provider := beego.AppConfig.String("session_provider")
	session_provider_config := beego.AppConfig.String("session_provider_config")
	beego.BConfig.WebConfig.Session.SessionProvider = session_provider
	beego.BConfig.WebConfig.Session.SessionProviderConfig = session_provider_config
	//初始化配置
	initDatabase()
	//注册接口
	registerInterface()
}

func registerInterface() {
	gob.Register(models.Account{})
}
