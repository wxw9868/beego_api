package controllers

import (
	"beeapi/models"
	"database/sql"
	"encoding/json"
	"github.com/astaxie/beego"
	"log"
	"strconv"
	"time"
)

type TokenAccountController struct {
	beego.Controller
}

// @Title CreateAccount 用户注册
// @router /createAccount [post]
func (this *TokenAccountController) CreateAccount() {
	json.Unmarshal(this.Ctx.Input.RequestBody, &account)

	nickname := account.Name     //用户昵称
	mobile := account.Mobile     //手机号
	email := account.Email       //邮箱
	loginPwd := account.LoginPwd //登录密码有【
	payPwd := account.PayPwd     //支付密码
	userType := account.UserType //用户类型:1:个人用户; 2:企业用户*/
	regType := account.RegType   //mobileReg:手机注册;emailReg:邮箱注册

	var id string

	switch {
	case userType == 1:
		if regType == "" {
			this.ReturnMsg(false, 0, "注册类型不能为空！", nil)
		}

		if regType == "mobileReg" {
			if mobile == 0 {
				this.ReturnMsg(false, 0, "手机号不能为空！", nil)
			}
			if VerifyMobile(strconv.Itoa(mobile)) == false {
				this.ReturnMsg(false, 0, "手机号格式不对！", nil)
			}
			sqlQuery = "SELECT id FROM token_user WHERE mobile = $1"
			err := db.QueryRow(sqlQuery, mobile).Scan(&id)
			switch {
			case err == sql.ErrNoRows:
				log.Printf("No user with that ID.")
			case err != nil:
				log.Fatal(err)
			default:
				this.ReturnMsg(false, 0, "手机号已注册", nil)
			}
		} else {
			if email == "" {
				this.ReturnMsg(false, 0, "邮箱不能为空！", nil)
			}
			if VerifyEmail(email) == false {
				this.ReturnMsg(false, 0, "邮箱格式不对！", nil)
			}
			sqlQuery = "SELECT id FROM token_user WHERE email = $1"
			err := db.QueryRow(sqlQuery, email).Scan(&id)
			switch {
			case err == sql.ErrNoRows:
				log.Printf("No user with that ID.")
			case err != nil:
				log.Fatal(err)
			default:
				this.ReturnMsg(false, 0, "邮箱已注册", nil)
			}
		}
	case userType == 2:
		if email == "" {
			this.ReturnMsg(false, 0, "邮箱不能为空！", nil)
		}
		if VerifyEmail(email) == false {
			this.ReturnMsg(false, 0, "邮箱格式不对！", nil)
		}
		sqlQuery = "SELECT id FROM token_user WHERE email = $1"
		err := db.QueryRow(sqlQuery, email).Scan(&id)
		switch {
		case err == sql.ErrNoRows:
			log.Printf("No user with that ID.")
		case err != nil:
			log.Fatal(err)
		default:
			this.ReturnMsg(false, 0, "邮箱已注册", nil)
		}
	default:
		this.ReturnMsg(false, 0, "用户类型不存在！", nil)
	}

	if loginPwd == "" {
		this.ReturnMsg(false, 0, "登录密码不能为空！", nil)
	}
	if payPwd == "" {
		this.ReturnMsg(false, 0, "支付密码不能为空！", nil)
	}
	loginPwd = DataEncryption(loginPwd)
	payPwd = DataEncryption(payPwd)

	if userType == 2 && mobile != 0 {
		if VerifyMobile(strconv.Itoa(mobile)) == false {
			this.ReturnMsg(false, 0, "手机号格式不对！", nil)
		}
		sqlQuery = "SELECT id FROM token_user WHERE mobile = $1"
		err := db.QueryRow(sqlQuery, mobile).Scan(&id)
		switch {
		case err == sql.ErrNoRows:
			log.Printf("No user with that ID.")
		case err != nil:
			log.Fatal(err)
		default:
			this.ReturnMsg(false, 0, "手机号已被使用", nil)
		}
	}

	sqlInsert = "INSERT INTO token_user (nickname,mobile,email,login_pwd,pay_pwd,user_type,balance,balance_lock,create_time,frozen) VALUES ($1,$2,$3,$4,$5,$6,DEFAULT,DEFAULT,$7,false)"
	if _, err := db.Exec(sqlInsert, nickname, mobile, email, loginPwd, payPwd, userType, time.Now().UTC()); err != nil {
		log.Fatal(err)
	}
	this.ReturnMsg(true, 10010, "注册成功！", nil)
}

// @Title 用户登录
// @router /accountLogin [post]
func (this *TokenAccountController) AccountLogin() {
	var accountUser models.AccountUser
	json.Unmarshal(this.Ctx.Input.RequestBody, &accountUser)

	_account := accountUser.Account
	_password := accountUser.Password

	if _password == "" {
		this.ReturnMsg(false, 0, "登录密码不能为空！", nil)
	}
	if VerifyMobile(_account) == true {
		sqlQuery = "SELECT id,nickname,email,login_pwd,pay_pwd,user_type,mobile,balance,balance_lock,frozen FROM token_user WHERE mobile = $1"
	} else if VerifyEmail(_account) == true {
		sqlQuery = "SELECT id,nickname,email,login_pwd,pay_pwd,user_type,mobile,balance,balance_lock,frozen FROM token_user WHERE email = $1"
	} else {
		this.ReturnMsg(false, 0, "账号格式错误！请重新输入。", nil)
	}

	var (
		user_id, nickname, email, login_pwd, pay_pwd string
		mobile, user_type                            int
		balance, balance_lock                        float64
		frozen                                       bool
	)
	err := db.QueryRow(sqlQuery, _account).Scan(&user_id, &nickname, &email, &login_pwd, &pay_pwd, &user_type, &mobile, &balance, &balance_lock, &frozen)
	switch {
	case err == sql.ErrNoRows:
		this.ReturnMsg(false, 10004, "账户不存在！请先注册账户。", nil)
	case err != nil:
		log.Fatal(err)
	default:
		if login_pwd != DataEncryption(_password) {
			this.ReturnMsg(false, 0, "密码错误！", nil)
		}
		if frozen == true {
			this.ReturnMsg(false, 10002, "账号已被冻结", nil)
		}

		account.UserId = user_id
		account.Name = nickname
		account.Mobile = mobile
		account.Email = email
		account.LoginPwd = login_pwd
		account.PayPwd = pay_pwd
		account.UserType = user_type
		account.Frozen = frozen
		account.BalanceOf = balance
		account.BalanceLock = balance_lock

		coinBase := make(map[string]interface{})
		coinBase["UserId"] = user_id
		coinBase["Name"] = nickname
		coinBase["Mobile"] = mobile
		coinBase["Email"] = email
		coinBase["UserType"] = user_type
		coinBase["Frozen"] = frozen
		coinBase["BalanceOf"] = balance
		coinBase["BalanceLock"] = balance_lock

		this.SetSession("account", account)
		this.ReturnMsg(true, 10001, "登录成功", coinBase)
	}
}

// @Title 用户退出
// @router /loginOut [get]
func (this *TokenAccountController) LoginOut() {
	this.DelSession("account")
	s := this.GetSession("account")
	if s == nil {
		this.ReturnMsg(true, 1, "退出成功", nil)
	}
}

//提示信息
func (this *TokenAccountController) ReturnMsg(_status bool, _code int, _value string, _data interface{}) {
	msg := map[string]interface{}{"Status": _status, "Code": _code, "Message": _value}
	if _data != nil {
		msg["data"] = _data
	}
	this.Data["json"] = &msg
	this.ServeJSON()
	this.StopRun()
}
