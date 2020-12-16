package controllers

import (
	"beeapi/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type TokenController struct {
	TokenBaseController
}

// @Title 发行代币
// @router /releaseTokens [post]
func (this *TokenController) ReleaseTokens() {
	var coinBase models.ReleaseTokens
	json.Unmarshal(this.Ctx.Input.RequestBody, &coinBase)

	payPwd := coinBase.PayPwd //支付密码
	if payPwd == "" {
		this.ReturnMsg(false, 0, "请输入支付密码。", nil)
	}
	if account.PayPwd != DataEncryption(payPwd) {
		this.ReturnMsg(false, 0, "支付密码错误！，请重新输入。", nil)
	}

	this.issueTokenPower()

	_name := coinBase.Name              //代币名称
	_symbol := coinBase.Symbol          //代币符号
	totalSupply := coinBase.TotalSupply //代币总量

	if _name == "" {
		this.ReturnMsg(false, 0, "请输入代币名称。", nil)
	}

	if _symbol == "" {
		this.ReturnMsg(false, 0, "请输入代币符号。", nil)
	}

	if totalSupply == "" {
		this.ReturnMsg(false, 0, "请输入代币总量。", nil)
	}

	_totalSupply, _ := strconv.ParseFloat(totalSupply, 64) //代币总量
	_description := coinBase.Description                   //代币介绍

	this.initialSupply(_name, _symbol, _totalSupply, _description, &account)
}

//发行代币逻辑
func (this *TokenController) initialSupply(_name string, _symbol string, _supply float64, _description string, _account *models.Account) {
	var id string
	sqlQuery = "SELECT id FROM token_token_manage WHERE symbol = $1"
	err := db.QueryRow(sqlQuery, _symbol).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that symbol.")
	case err != nil:
		log.Fatal(err)
	default:
		this.ReturnMsg(false, 10010, "代币已经存在！", nil)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	create_time := time.Now().UTC()

	//写入代币管理表
	sqlInsert = "INSERT INTO token_token_manage (user_id,name,symbol,total_supply,description,lock,create_time) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	if _, err = tx.Exec(sqlInsert, _account.UserId, _name, _symbol, _supply, _description, false, create_time); err != nil {
		log.Fatal(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
		}
	}

	//写入用户余额表
	sqlInsert = "INSERT INTO token_balance (user_id,symbol,token_balance,token_balance_lock,create_time) VALUES ($1,$2,$3,$4,$5)"
	if _, err = tx.Exec(sqlInsert, _account.UserId, _symbol, _supply, _account.BalanceLock, create_time); err != nil {
		log.Fatal(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
		}
	}
	//获取用户余额的id
	var balance_id string
	sqlQuery = "SELECT id FROM token_balance WHERE user_id = $1 AND symbol = $2"
	err = tx.QueryRow(sqlQuery, _account.UserId, _symbol).Scan(&balance_id)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that symbol,user_id")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("balance_id is %s\n", balance_id)
	}

	sqlInsert = "INSERT INTO token_balance_action (user_id,balance_id,amount,behavior,to_user_id,create_time) VALUES ($1,$2,$3,$4,$5,$6)"
	if _, err = tx.Exec(sqlInsert, _account.UserId, balance_id, _supply, "in", UserId, create_time); err != nil {
		log.Fatal(err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
		}
	}
	if err = tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			this.ReturnMsg(false, 0, "回滚事务失败！", nil)
		}
	}
	this.ReturnMsg(true, 10011, "代币发布成功！", nil)
}

// @Title 查询所有类型余额
// @router /balanceTokenAll [get]
func (this *TokenController) BalanceTokenAll() {
	sqlQuery = "SELECT symbol,token_balance FROM token_balance WHERE user_id = $1"
	rows, err := db.Query(sqlQuery, UserId)
	if err != nil {
		log.Fatal(err)
	}
	var symbol string
	var token_balance float64
	var result = make(map[string]float64)
	for rows.Next() {
		if err := rows.Scan(&symbol, &token_balance); err != nil {
			log.Fatal(err)
		}
		result[symbol] = token_balance
	}
	this.ReturnMsg(true, 10007, "查询成功", result)
}

// @Title 查询余额
// @router /balanceToken [post]
func (this *TokenController) BalanceToken() {
	var token models.ReleaseTokens
	json.Unmarshal(this.Ctx.Input.RequestBody, &token)

	_symbol := token.Symbol
	if _symbol == "" {
		this.ReturnMsg(false, 0, "请输入代币符号", nil)
	}

	var symbol string
	var token_balance float64
	var result = make(map[string]float64)

	sqlQuery = "SELECT symbol,token_balance FROM token_balance WHERE user_id=$1 AND symbol=$2"
	err := db.QueryRow(sqlQuery, UserId, _symbol).Scan(&symbol, &token_balance)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that symbol.")
	case err != nil:
		log.Fatal(err)
	default:
		result[symbol] = token_balance
	}
	this.ReturnMsg(true, 10007, "查询成功", result)
}

// @Title 冻结账号
// @router /frozenAccount [post]
func (this *TokenController) FrozenAccount() {
	var frozen models.Account
	json.Unmarshal(this.Ctx.Input.RequestBody, &frozen)
	_status := frozen.Frozen

	var status bool
	var msg string
	var code int
	if _status == true {
		status = true
		msg = "账号冻结！"
		code = 10002
	} else {
		status = false
		msg = "账号解冻！"
		code = 10008
	}
	account.Frozen = status

	sqlUpdate = "UPDATE token_user SET frozen = $1 WHERE id = $2"
	if _, err := db.Exec(sqlUpdate, status, UserId); err != nil {
		log.Fatal(err)
	}
	this.ReturnMsg(true, code, msg, nil)
}

// @Title 转账
// @router /transferToken [post]
func (this *TokenController) TransferToken() {
	var transferToken models.ReleaseTokens
	json.Unmarshal(this.Ctx.Input.RequestBody, &transferToken)

	payPwd := transferToken.PayPwd
	_to := transferToken.Account
	_symbol := transferToken.Symbol
	_amount, _ := strconv.ParseFloat(transferToken.TotalSupply, 64)

	if _amount <= 0 {
		this.ReturnMsg(false, 0, "totalSupply必须大于0", nil)
	}

	this.checkSymbol(_symbol)

	if account.PayPwd != DataEncryption(payPwd) {
		this.ReturnMsg(false, 0, "支付密码不正确！", nil)
	}

	var (
		user_id, email, balance_id, to_balance_id string
		mobile                                    int
		token_balance                             float64
		frozen                                    bool
	)

	if VerifyMobile(_to) == true {
		sqlQuery = "SELECT id,email,mobile,frozen FROM token_user WHERE mobile = $1"
	} else if VerifyEmail(_to) == true {
		sqlQuery = "SELECT id,email,mobile,frozen FROM token_user WHERE email = $1"
	} else {
		this.ReturnMsg(false, 10003, "账号格式错误！请重新输入。", nil)
	}
	var toAccount models.Account
	err := db.QueryRow(sqlQuery, _to).Scan(&user_id, &email, &mobile, &frozen)
	switch {
	case err == sql.ErrNoRows:
		this.ReturnMsg(false, 10004, "账户不存在！请先注册账户。", nil)
	case err != nil:
		log.Fatal(err)
	default:
		toAccount = models.Account{UserId: user_id, Email: email, Mobile: mobile, Frozen: frozen}
	}

	//查询转账账户的余额
	sqlQuery = "SELECT id,token_balance FROM token_balance WHERE user_id = $1 AND symbol = $2"
	err = db.QueryRow(sqlQuery, UserId, _symbol).Scan(&balance_id, &token_balance)
	switch {
	case err == sql.ErrNoRows:
		this.ReturnMsg(false, 10005, "没有此种代币！", nil)
	case err != nil:
		log.Fatal(err)
	default:
		account.BalanceOf = token_balance
	}

	res := this.transfer(&account, &toAccount, _symbol, _amount)

	if res == nil {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		sqlQuery = "SELECT id,token_balance FROM token_balance WHERE user_id = $1 AND symbol = $2"
		err = tx.QueryRow(sqlQuery, toAccount.UserId, _symbol).Scan(&to_balance_id, &token_balance)
		switch {
		case err == sql.ErrNoRows:
			sqlQuery = "INSERT INTO token_balance (user_id,symbol,token_balance,token_balance_lock,create_time) VALUES ($1,$2,$3,$4,$5)"
			_, err = tx.Exec(sqlQuery, toAccount.UserId, _symbol, toAccount.BalanceOf, toAccount.BalanceLock, time.Now().UTC())
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
				}
				log.Fatal(err)
			}
			sqlQuery = "SELECT id FROM token_balance WHERE user_id = $1 AND symbol = $2"
			err = tx.QueryRow(sqlQuery, toAccount.UserId, _symbol).Scan(&to_balance_id)
			switch {
			case err == sql.ErrNoRows:
				log.Printf("没有此种代币")
			case err != nil:
				log.Fatal(err)
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
				}
			default:
				fmt.Println("to_balance_id:", balance_id)
			}
		case err != nil:
			log.Fatal(err)
		default:
			sqlUpdate = "UPDATE token_balance SET token_balance = token_balance + $1 WHERE user_id = $2 AND symbol = $3"
			_, err = tx.Exec(sqlUpdate, toAccount.BalanceOf, toAccount.UserId, _symbol)
			if err != nil {
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
				}
				log.Fatal(err)
			}
		}

		sqlUpdate = "UPDATE token_balance SET token_balance = $1 WHERE user_id = $2 AND symbol = $3"
		_, err = tx.Exec(sqlUpdate, account.BalanceOf, UserId, _symbol)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
			log.Fatal(err)
		}

		sqlInsert = "INSERT INTO token_balance_action (user_id, balance_id,amount,behavior,to_user_id,create_time) VALUES ($1,$2,$3,$4,$5,$6)"
		if _, err = tx.Exec(sqlInsert, toAccount.UserId, to_balance_id, toAccount.BalanceOf, "in", UserId, time.Now().UTC()); err != nil {
			log.Fatal("err: ", err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}

		sqlInsert = "INSERT INTO token_balance_action (user_id, balance_id,amount,behavior,to_user_id,create_time) VALUES ($1,$2,$3,$4,$5,$6)"
		if _, err := tx.Exec(sqlInsert, UserId, balance_id, account.BalanceOf, "out", toAccount.UserId, time.Now().UTC()); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
			log.Fatal(err)
		}
		if err = tx.Commit(); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
			log.Fatal(err)
		}
	}

	this.ReturnMsg(true, 10022, "转账成功", nil)
}

//转账交易逻辑
func (this *TokenController) transfer(_from *models.Account, _to *models.Account, _symbol string, _value float64) []byte {
	if token.Lock {
		this.ReturnMsg(false, 10030, "锁仓状态，停止一切转账活动", nil)
	}
	if _from.Frozen {
		this.ReturnMsg(false, 10002, "账号冻结！", nil)
	}
	if _to.Frozen {
		this.ReturnMsg(false, 10002, "账号冻结！", nil)
	}

	if _from.BalanceOf >= _value {
		_from.BalanceOf -= _value
		_to.BalanceOf += _value
		return nil
	} else {
		this.ReturnMsg(false, 10020, "余额不足！", nil)
	}
	return []byte("1")
}

// @Title 代币仓冻结
// @router /setLock [post]
func (this *TokenController) SetLock() {
	var setLock models.Token
	json.Unmarshal(this.Ctx.Input.RequestBody, &setLock)

	_symbol := setLock.Symbol
	_lock := setLock.Lock

	this.checkSymbol(_symbol)
	if token.UserId != UserId {
		this.ReturnMsg(false, 10014, "没有冻结代币仓的权限！", nil)
	}

	var msg *models.Message

	if _lock == true {
		token.Lock = true
		msg = &models.Message{Status: true, Code: 10013, Message: "代币仓冻结，停止一切交易活动"}
	} else {
		token.Lock = false
		msg = &models.Message{Status: true, Code: 10017, Message: "代币解冻，交易活动恢复"}
	}

	sqlUpdate = "UPDATE token_token_manage SET lock = $1 WHERE symbol = $2"
	if _, err := db.Exec(sqlUpdate, _lock, _symbol); err == nil {
		this.Data["json"] = msg
	} else {
		msg = &models.Message{Status: false, Code: 10015, Message: "代币仓冻结失败"}
		this.Data["json"] = msg
	}
	this.ServeJSON()
}

//检查代币是否存在
func (this *TokenController) checkSymbol(_symbol string) {
	sqlQuery = "SELECT id,user_id,name,total_supply,description,lock FROM token_token_manage WHERE symbol = $1"
	err := db.QueryRow(sqlQuery, _symbol).Scan(&manage_id, &token_user_id, &token_name, &total_supply, &description, &lock)
	switch {
	case err == sql.ErrNoRows:
		this.ReturnMsg(false, 0, "代币不存在！", nil)
	case err != nil:
		log.Fatal(err)
	default:
		token = models.Token{Lock: lock, UserId: token_user_id, Name: token_name, Symbol: _symbol, TotalSupply: total_supply, Description: description, Id: manage_id}
	}
}

// @Title 增发代币
// @router /addTokens [post]
func (this *TokenController) AddTokens() {
	var addToken models.ReleaseTokens
	json.Unmarshal(this.Ctx.Input.RequestBody, &addToken)

	payPwd := addToken.PayPwd
	_symbol := addToken.Symbol
	amount := addToken.TotalSupply
	if payPwd == "" {
		this.ReturnMsg(false, 0, "请填写支付密码！", nil)
	}
	if _symbol == "" {
		this.ReturnMsg(false, 0, "请填写代币符号！", nil)
	}
	if amount == "" {
		this.ReturnMsg(false, 0, "请填写代币数量！", nil)
	}

	_amount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		this.ReturnMsg(false, 0, "增发的代币不准确！", nil)
	}
	if _amount <= 0 {
		this.ReturnMsg(false, 0, "增发的代币量必须大于0！", nil)
	}

	if account.PayPwd != DataEncryption(payPwd) {
		this.ReturnMsg(false, 0, "支付密码不正确！", nil)
	}

	this.issueTokenPower()
	this.checkSymbol(_symbol)

	if token.UserId != UserId {
		this.ReturnMsg(false, 10014, "没有冻结代币仓的权限！", nil)
	}
	this.doAddToken(_symbol, _amount, &account)
}

//增发代币的逻辑
func (this *TokenController) doAddToken(_symbol string, _value float64, account *models.Account) {
	if token.Lock {
		this.ReturnMsg(false, 10013, "代币仓冻结，停止一切交易活动！", nil)
	} else if account.Frozen {
		this.ReturnMsg(false, 10002, "账号冻结", nil)
	} else {
		var balance_id string
		sqlQuery = "SELECT id FROM token_balance WHERE symbol = $1 AND user_id = $2"
		err := db.QueryRow(sqlQuery, _symbol, token.UserId).Scan(&balance_id)
		switch {
		case err == sql.ErrNoRows:
			log.Printf("No user with that ID.")
		case err != nil:
			log.Fatal(err)
		default:
			fmt.Printf("balance_id is %s\n", balance_id)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(nil)
		}
		sqlInsert = "INSERT INTO token_token_add_log (manage_id,user_id,add_supply,create_time) VALUES ($1,$2,$3,$4)"
		_, err = tx.Exec(sqlInsert, token.Id, token.UserId, _value, time.Now().UTC())
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
			log.Fatal(err)
		}
		sqlUpdate = "UPDATE token_balance SET token_balance = token_balance + $1 WHERE symbol = $2 AND user_id = $3"
		_, err = tx.Exec(sqlUpdate, _value, _symbol, token.UserId)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}

		sqlInsert = "INSERT INTO token_balance_action (user_id,balance_id,amount,behavior,to_user_id,create_time) VALUES ($1,$2,$3,$4,$5,$6)"
		if _, err = tx.Exec(sqlInsert, account.UserId, balance_id, _value, "in", UserId, time.Now().UTC()); err != nil {
			log.Fatal(err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}
		if err = tx.Commit(); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
		}
		this.ReturnMsg(true, 10016, "代币增发成功", nil)
	}
}

func (this *TokenController) issueTokenPower() {
	if account.UserType != 2 {
		this.ReturnMsg(false, 10006, "用户没有代币发行权限！", nil)
	}
}

//支付接口
//payType(RMB:代表使用人民币支付；TOKEN:代表使用代币支付；MP(Mixed payment):代表使用混合支付方式)
//goodsType(EPG:代表体验商品)
//func (this *TokenController) pay() {
//	order_sn := this.GetString("order_sn")             //订单号
//	order_name := this.GetString("order_name")         //订单名称
//	order_amount, _ := this.GetFloat("order_amount")  	//订单总金额
//	order_token, _ := this.GetFloat("order_token")		//代币总额
//	order_note := this.GetString("order_note")         //订单备注
//	payType := this.GetString("payType")               //支付方式
//	goodsType := this.GetString("payType")             //商品类型
//	appId := this.GetString("payType")                 //商户（用户）id
//	payPwd := this.GetString("payType")                //支付密码
//	_symbol := this.GetString("_symbol")               //代币符号
//	_time := this.GetString("time")                    //代币锁定时间
//
//	timestamp := time.Now().Unix()
//
//	if number == "" {
//		return shim.Error("订单号不能为空！")
//	}
//	if name == "" {
//		return shim.Error("订单名称不能为空！")
//	}
//	if args[2] == "" {
//		return shim.Error("订单总金额不能为空！")
//	}
//	if payType == "" {
//		return shim.Error("支付方式不能为空！")
//	}
//	if payType == "MP" && args[3] == "" {
//		return shim.Error("代币总额不能为空！")
//	}
//	if appId == "" {
//		return shim.Error("商户ID不能为空！")
//	}
//	if timestamp == "" {
//		return shim.Error("支付时间不能为空！")
//	}
//	if _time == "" {
//		return shim.Error("代币锁定时间不能为空！")
//	}
//	ti, _ := strconv.Atoi(_time)
//	if ti <= 0 {
//		return shim.Error("代币锁定时间不能为0！")
//	}
//
//	//校验支付用户信息
//	userBytes, err := stub.GetState(_account)
//	if string(userBytes) == "" {
//		return shim.Error(StringBuilder(_account, "用户不存在!"))
//	}
//	account := Account{}
//	json.Unmarshal(userBytes, &account)
//	if account.LoginPwd != EncryptedPasswords(loginPwd) {
//		return shim.Error(StringBuilder(_account, "的登录密码错误！"))
//	}
//	if account.PayPwd != EncryptedPasswords(payPwd) {
//		return shim.Error(StringBuilder(_account, "的支付密码错误！"))
//	}
//
//	//支付操作
//	shopBytes, _ := stub.GetState(appId)
//	if string(shopBytes) == "" {
//		return shim.Error("商户不存在!")
//	}
//	shop := Account{}
//	json.Unmarshal(shopBytes, &shop)
//
//	tokenBytes, _ := stub.GetState(_symbol)
//	if string(tokenBytes) == "" {
//		return shim.Error(StringBuilder(_symbol, "代币不存在!"))
//	}
//	token := Token{}
//	json.Unmarshal(tokenBytes, &token)
//
//	var issuingAccountBytes []byte
//	var issuingAccount Account
//	if shop.Email != token.User {
//		issuingAccountBytes, err = stub.GetState(token.User)
//		if err != nil {
//			return shim.Error(err.Error())
//		}
//		issuingAccount = Account{}
//		json.Unmarshal(issuingAccountBytes, &issuingAccount)
//	} else {
//		issuingAccount = shop
//	}
//
//	result := token.transferA(&account, &shop, &issuingAccount, _symbol, _amount, _token, payType, goodsType)
//
//	message := Message{}
//	json.Unmarshal(result, &message)
//	if message.Code != 10000 {
//		return shim.Success(result)
//	}
//
//	accountBytes, err := json.Marshal(account)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	stub.PutState(_account, accountBytes)
//	stub.PutState(account.UserId, accountBytes)
//
//	shopAsBytes, err := json.Marshal(shop)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	fmt.Println(shop)
//	stub.PutState(shop.Email, shopAsBytes)
//	stub.PutState(appId, shopAsBytes)
//
//	if shop.Email != token.User {
//		issuingAccountBytes, err = json.Marshal(issuingAccount)
//		if err != nil {
//			return shim.Error(err.Error())
//		}
//		stub.PutState(token.User, issuingAccountBytes)
//		stub.PutState(issuingAccount.UserId, issuingAccountBytes)
//	}
//
//	lockTokenList := LockTokenList{}
//	lockTokenList.setTokenLockTime(_symbol, _amount, _account, account.UserId, _time, number, payType, stub)
//
//	//生成支付订单
//	orderLists := OrderLists{}
//	newId, key := createNewId(appId, number)
//	orderListsBytes, err := stub.GetState(newId)
//	if string(orderListsBytes) == "" {
//		orderLists.OrderLists = map[string]OrderMsg{key: {OrderNumber: number, OrderName: name, OrderNote: note, TotalAmount: _amount, UserId: account.UserId, Timestamp: timestamp, PayStatus: "SUCCESS", RefundStatus: "", RefundTime: "", PayType: payType}}
//		orderListsBytes, err = json.Marshal(orderLists)
//		if err != nil {
//			return shim.Error(err.Error())
//		}
//		stub.PutState(newId, orderListsBytes)
//	} else {
//		json.Unmarshal(orderListsBytes, &orderLists)
//		if orderLists.OrderLists[key].PayStatus == "SUCCESS" {
//			return shim.Error("订单已经支付，不可重复支付！")
//		}
//		orderLists.OrderLists = map[string]OrderMsg{key: {OrderNumber: number, OrderName: name, OrderNote: note, TotalAmount: _amount, UserId: account.UserId, Timestamp: timestamp, PayStatus: "SUCCESS", RefundStatus: "", RefundTime: "", PayType: payType}}
//		orderListsBytes, err = json.Marshal(orderLists)
//		if err != nil {
//			return shim.Error(err.Error())
//		}
//		err = stub.PutState(newId, orderListsBytes)
//		if err != nil {
//			shim.Error(err.Error())
//		}
//	}
//
//	return shim.Success(orderListsBytes)
//}

//支付交易
//func (token *Token) transferA(_from *Account, _to *Account, issuingAccount *Account, _symbol string, _value float64, _token float64, payType string, goodsType string) []byte {
//	if token.Lock {
//		msg := &Message{Status: true, Code: 10001, Message: "锁仓状态，停止一切转账活动"}
//		res, _ := json.Marshal(msg)
//		return res
//	}
//	if _from.Frozen {
//		msg := &Message{Status: true, Code: 10002, Message: "账号冻结！"}
//		res, _ := json.Marshal(msg)
//		return res
//	}
//	if _to.Frozen {
//		msg := &Message{Status: true, Code: 10002, Message: "账号冻结！"}
//		res, _ := json.Marshal(msg)
//		return res
//	}
//	msg := &Message{Status: true, Code: 10000, Message: "支付成功！"}
//	res, _ := json.Marshal(msg)
//	if payType == "RMB" {
//		if goodsType == "EPG" {
//			fmt.Println(_value)
//			issuingAccount.BalanceOf[_symbol] -= _value
//			_from.BalanceLock[_symbol] += _value
//		}
//		return res
//	}
//	if payType == "MP" {
//		_from.BalanceOf[_symbol] -= _token
//		_to.BalanceOf[_symbol] += _token
//		if goodsType == "EPG" {
//			r := _value - _token
//			issuingAccount.BalanceOf[_symbol] -= r
//			_from.BalanceLock[_symbol] += _value
//		}
//		return res
//	}
//	if payType == "TOKEN" && _from.BalanceOf[_symbol] >= _value {
//		_from.BalanceOf[_symbol] -= _value
//		_to.BalanceOf[_symbol] += _value
//		if goodsType == "EPG" {
//			_from.BalanceLock[_symbol] += _value
//		}
//		return res
//	} else {
//		msg := &Message{Status: false, Code: 10003, Message: "余额不足！"}
//		res, _ := json.Marshal(msg)
//		return res
//	}
//}

//锁定部分代币一段时间
//func (l *LockTokenList) setTokenLockTime(_symbol string, _amount float64, _account string, userId string, _time string, number string, _type string, stub shim.ChaincodeStubInterface) {
//	startTime := time.Now().Unix()
//	timestamp, _ := strconv.Atoi(_time)
//	endTime := startTime + int64(timestamp)*24*3600
//	key, _ := stub.CreateCompositeKey("LockToken", []string{_account, userId})
//	resultBytes, _ := stub.GetState(key)
//	if string(resultBytes) != "" {
//		json.Unmarshal(resultBytes, l)
//	}
//	l.LockList = map[string]LockToken{number: {UserId: userId, Amount: _amount, Symbol: _symbol, StartTime: startTime, EndTime: endTime, Status: 0, OrderNumber: number, PayType: _type}}
//	resultBytes, _ = json.Marshal(l)
//	stub.PutState(key, resultBytes)
//}

//支付查询接口
//func (t *TokenController) payQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	if len(args) != 2 {
//		return shim.Error("参数个数不正确，为2个")
//	}
//	appId := args[0]
//	number := args[1]
//
//	newId, key := createNewId(appId, number)
//	orderLists := OrderLists{}
//	orderListsBytes, _ := stub.GetState(newId)
//	if string(orderListsBytes) == "" {
//		return shim.Error("请检查AppId和商品订单号是否正确！")
//	}
//	json.Unmarshal(orderListsBytes, &orderLists)
//
//	var resultBytes []byte
//
//	if orderLists.OrderLists[key].PayStatus == "SUCCESS" {
//		result := &Message{Status: true, Code: 10006, Message: "支付成功。"}
//		resultBytes, _ = json.Marshal(result)
//	} else {
//		return shim.Error("订单号码错误！")
//	}
//	return shim.Success(resultBytes)
//}

//余额充值
//func (t *TokenController) recharge(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	if len(args) != 5 {
//		return shim.Error("Incorrect number of arguments. Expecting 5")
//	}
//	_account := args[0]                           //账号：邮箱/手机号
//	loginPwd := args[1]                           //
//	payPwd := args[2]                             //
//	_amount, _ := strconv.ParseFloat(args[3], 64) //
//	_symbol := args[4]                            //
//
//	symbolBytes, _ := stub.GetState(_symbol)
//	if string(symbolBytes) == "" {
//		return shim.Error(StringBuilder(_symbol, "代币不存在!"))
//	}
//
//	account := Account{}
//	accountBytes, _ := stub.GetState(_account)
//	if string(accountBytes) == "" {
//		return shim.Error(StringBuilder(_account, "用户不存在!"))
//	}
//	json.Unmarshal(accountBytes, &account)
//	if EncryptedPasswords(loginPwd) != account.LoginPwd {
//		return shim.Error("登录密码不正确！")
//	}
//	if EncryptedPasswords(payPwd) != account.PayPwd {
//		return shim.Error("支付密码不正确！")
//	}
//	if _amount <= 0 {
//		return shim.Error("充值金额不能小于等于0！")
//	}
//
//	token := Token{}
//	tokenBytes, err := stub.GetState(_symbol)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	json.Unmarshal(tokenBytes, &token)
//
//	issuingAccount := Account{}
//	issuingAccountBytes, err := stub.GetState(token.User)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	json.Unmarshal(issuingAccountBytes, &issuingAccount)
//
//	result := token.transfer(&issuingAccount, &account, _symbol, _amount)
//	message := Message{}
//	json.Unmarshal(result, &message)
//	if message.Code != 10010 {
//		return shim.Success(result)
//	}
//
//	accountBytes, err = json.Marshal(account)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	stub.PutState(_account, accountBytes)
//	stub.PutState(account.UserId, accountBytes)
//
//	issuingAccountBytes, err = json.Marshal(issuingAccount)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	stub.PutState(token.User, issuingAccountBytes)
//	stub.PutState(issuingAccount.UserId, issuingAccountBytes)
//
//	res := &Message{Status: true, Code: 10009, Message: "充值成功。"}
//	resBytes, _ := json.Marshal(res)
//
//	return shim.Success(resBytes)
//}

//提现接口
//func (t *TokenController) withdraw(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	if len(args) != 5 {
//		return shim.Error("Incorrect number of arguments. Expecting 5")
//	}
//	_account := args[0]
//	loginPwd := args[1]
//	payPwd := args[2]
//	_amount, _ := strconv.ParseFloat(args[3], 64)
//	_symbol := args[4]
//
//	account := Account{}
//	accountBytes, _ := stub.GetState(_account)
//	if string(accountBytes) == "" {
//		return shim.Error(StringBuilder(_account, "用户不存在！"))
//	}
//	json.Unmarshal(accountBytes, &account)
//	if EncryptedPasswords(loginPwd) != account.LoginPwd {
//		return shim.Error("登录密码不正确！")
//	}
//	if EncryptedPasswords(payPwd) != account.PayPwd {
//		return shim.Error("支付密码不正确！")
//	}
//	if _amount <= 0 {
//		return shim.Error("提现金额不能小于等于0！")
//	}
//
//	token := Token{}
//	tokenBytes, err := stub.GetState(_symbol)
//	if string(tokenBytes) == "" {
//		return shim.Error(StringBuilder(_symbol, "不存在！"))
//	}
//	json.Unmarshal(tokenBytes, &token)
//
//	issuingAccount := Account{}
//	issuingAccountBytes, err := stub.GetState(token.User)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	json.Unmarshal(issuingAccountBytes, &issuingAccount)
//
//	result := token.transfer(&account, &issuingAccount, _symbol, _amount)
//	message := Message{}
//	json.Unmarshal(result, &message)
//	if message.Code != 10010 {
//		return shim.Success(result)
//	}
//
//	accountBytes, err = json.Marshal(account)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	stub.PutState(_account, accountBytes)
//	stub.PutState(account.UserId, accountBytes)
//
//	issuingAccountBytes, err = json.Marshal(issuingAccount)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	stub.PutState(token.User, issuingAccountBytes)
//	stub.PutState(issuingAccount.UserId, issuingAccountBytes)
//
//	res := &Message{Status: true, Code: 10008, Message: "提现成功。"}
//	resBytes, _ := json.Marshal(res)
//
//	return shim.Success(resBytes)
//}

//退款接口
//func (t *TokenController) refund(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	if len(args) != 4 {
//		return shim.Error("Incorrect number of arguments. Expecting 4")
//	}
//	appId := args[0]                              //商户（用户）ID
//	_symbol := args[1]                            //代币符号
//	number := args[2]                             //订单号
//	_amount, _ := strconv.ParseFloat(args[3], 64) //退款金额
//
//	symbolBytes, _ := stub.GetState(_symbol)
//	if string(symbolBytes) == "" {
//		return shim.Error(StringBuilder(_symbol, "代币不存在！"))
//	}
//
//	newId, key := createNewId(appId, number)
//
//	orderListsBytes, _ := stub.GetState(newId)
//	if string(orderListsBytes) == "" {
//		return shim.Error("请检查AppId和商品订单号是否正确！")
//	}
//	orderLists := OrderLists{}
//	json.Unmarshal(orderListsBytes, &orderLists)
//	fmt.Println(orderLists)
//	data := orderLists.OrderLists[key]
//	if data.RefundStatus == "SUCCESS" {
//		return shim.Error("已经退款成功，不能重复退款！")
//	}
//	if data.TotalAmount != _amount {
//		return shim.Error("退款金额不正确！")
//	}
//
//	shop := Account{}
//	shopBytes, err := stub.GetState(appId)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	json.Unmarshal(shopBytes, &shop)
//
//	account := Account{}
//	accountBytes, err := stub.GetState(data.UserId)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	json.Unmarshal(accountBytes, &account)
//
//	token := Token{}
//	result := token.transfer(&shop, &account, _symbol, _amount)
//	message := Message{}
//	json.Unmarshal(result, &message)
//	if message.Code != 10010 {
//		return shim.Success(result)
//	}
//
//	shopBytes, err = json.Marshal(shop)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	stub.PutState(appId, shopBytes)
//	stub.PutState(shop.Email, shopBytes)
//
//	accountAsBytes, err := json.Marshal(account)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	stub.PutState(data.UserId, accountAsBytes)
//	_account := strconv.Itoa(account.Mobile)
//	stub.PutState(_account, accountAsBytes)
//
//	data.RefundStatus = "SUCCESS"
//	timeStr := time.Now().Format("2006-01-02 15:04:05")
//	data.RefundTime = timeStr
//	orderLists.OrderLists[key] = data
//	orderListsBytes, err = json.Marshal(orderLists)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	stub.PutState(newId, orderListsBytes)
//
//	return shim.Success(orderListsBytes)
//}

//退款查询接口
//func (t *TokenController) refundQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
//	if len(args) != 2 {
//		return shim.Error("Incorrect number of arguments. Expecting 2")
//	}
//	AppId := args[0]
//	number := args[1]
//
//	newId, key := createNewId(AppId, number)
//
//	orderListsBytes, _ := stub.GetState(newId)
//	if string(orderListsBytes) == "" {
//		return shim.Error("请检查AppId和商品订单号是否正确！")
//	}
//	orderLists := OrderLists{}
//	json.Unmarshal(orderListsBytes, &orderLists)
//	data := orderLists.OrderLists[key]
//
//	var resultBytes []byte
//	if data.RefundStatus == "SUCCESS" {
//		result := &Message{Status: true, Code: 10007, Message: "退款成功。"}
//		resultBytes, _ = json.Marshal(result)
//	} else {
//		return shim.Error("订单号码错误！")
//	}
//	return shim.Success(resultBytes)
//}

//func (t *TokenController) userInfo(stub shim.ChaincodeStubInterface, args []string) {
//	if len(args) != 2 {
//		return shim.Error("Incorrect number of arguments. Expecting 2")
//	}
//	name := args[0]
//	loginPwd := args[1]
//
//	resultBytes, _ := stub.GetState(name)
//	if string(resultBytes) == "" {
//		return shim.Error("此账户不存在！")
//	}
//	account := Account{}
//	json.Unmarshal(resultBytes, &account)
//	if account.LoginPwd != DataEncryption(loginPwd) {
//		return shim.Error("密码错误，请重新输入")
//	}
//}
