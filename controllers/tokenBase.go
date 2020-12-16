package controllers

import (
	"beeapi/models"
	"beeapi/sqlinit"
	"github.com/astaxie/beego"
)

type TokenBaseController struct {
	beego.Controller
}

var db = sqlinit.Db
var UserId = "0"
var account models.Account
var token models.Token
var typeMsg = map[int]string{1: "企业账户", 2: "个人账户"}
var sqlQuery, sqlUpdate, sqlInsert string

//token_token_manage表的全局变量
var (
	manage_id, token_user_id, description, token_name string
	total_supply                                      float64
	lock                                              bool
	create_time                                       int
)

func (this *TokenBaseController) Prepare() {
	s := this.GetSession("account")
	if s == nil {
		msg := &models.Message{Status: false, Code: 10000, Message: "需要重新登陆！"}
		this.Data["json"] = &msg
		this.ServeJSON()
	}
	//类型断言
	v, ok := s.(models.Account)
	if ok {
		if v.UserId != "0" {
			UserId = v.UserId
			account = v
		}
	}
	this.checkLogin()
}

// @Title 检测登录状态
// @router /checkLogin [post]
func (this *TokenBaseController) checkLogin() {
	if UserId == "0" {
		this.ReturnMsg(false, 10000, "需要重新登陆！", nil)
	}
}

//信息提示
func (this *TokenBaseController) ReturnMsg(_status bool, _code int, _value string, _data interface{}) {
	msg := map[string]interface{}{"Status": _status, "Code": _code, "Message": _value}
	if _data != nil {
		msg["data"] = _data
	}
	this.Data["json"] = &msg
	this.ServeJSON()
	this.StopRun()
}
