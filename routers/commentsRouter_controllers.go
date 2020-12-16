package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["beeapi/controllers:ObjectController"] = append(beego.GlobalControllerRouter["beeapi/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:ObjectController"] = append(beego.GlobalControllerRouter["beeapi/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:ObjectController"] = append(beego.GlobalControllerRouter["beeapi/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:objectId`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:ObjectController"] = append(beego.GlobalControllerRouter["beeapi/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:objectId`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:ObjectController"] = append(beego.GlobalControllerRouter["beeapi/controllers:ObjectController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:objectId`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenAccountController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenAccountController"],
		beego.ControllerComments{
			Method:           "AccountLogin",
			Router:           `/accountLogin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenAccountController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenAccountController"],
		beego.ControllerComments{
			Method:           "CreateAccount",
			Router:           `/createAccount`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenAccountController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenAccountController"],
		beego.ControllerComments{
			Method:           "LoginOut",
			Router:           `/loginOut`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenBaseController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenBaseController"],
		beego.ControllerComments{
			Method:           "checkLogin",
			Router:           `/checkLogin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenController"],
		beego.ControllerComments{
			Method:           "AddTokens",
			Router:           `/addTokens`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenController"],
		beego.ControllerComments{
			Method:           "BalanceToken",
			Router:           `/balanceToken`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenController"],
		beego.ControllerComments{
			Method:           "BalanceTokenAll",
			Router:           `/balanceTokenAll`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenController"],
		beego.ControllerComments{
			Method:           "FrozenAccount",
			Router:           `/frozenAccount`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenController"],
		beego.ControllerComments{
			Method:           "ReleaseTokens",
			Router:           `/releaseTokens`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenController"],
		beego.ControllerComments{
			Method:           "SetLock",
			Router:           `/setLock`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:TokenController"] = append(beego.GlobalControllerRouter["beeapi/controllers:TokenController"],
		beego.ControllerComments{
			Method:           "TransferToken",
			Router:           `/transferToken`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:UserController"] = append(beego.GlobalControllerRouter["beeapi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:UserController"] = append(beego.GlobalControllerRouter["beeapi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:UserController"] = append(beego.GlobalControllerRouter["beeapi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Get",
			Router:           `/:uid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:UserController"] = append(beego.GlobalControllerRouter["beeapi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           `/:uid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:UserController"] = append(beego.GlobalControllerRouter["beeapi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:uid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:UserController"] = append(beego.GlobalControllerRouter["beeapi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beeapi/controllers:UserController"] = append(beego.GlobalControllerRouter["beeapi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Logout",
			Router:           `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
