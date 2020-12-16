package models

type AccountUser struct {
	Account  string
	Password string
}

type Account struct {
	Name        string  `json:"Name"`        //Name
	Frozen      bool    `json:"Frozen"`      //Frozen
	BalanceOf   float64 `json:"BalanceOf"`   //BalanceOf
	BalanceLock float64 `json:"BalanceLock"` //BalanceLock
	Mobile      int     `json:"Mobile"`      //手机号
	Email       string  `json:"Email"`       //电子邮箱
	LoginPwd    string  `json:"LoginPwd"`    //登录密码
	PayPwd      string  `json:"PayPwd"`      //支付密码
	UserType    int     `json:"UserType"`    //用户类型：2:企业用户;1:个人用户
	UserId      string  `json:"UserId"`      //用户ID
	RegType     string  `json:"regType"`     //注册类型
}

type ReleaseTokens struct {
	PayPwd      string
	Name        string
	Symbol      string
	TotalSupply string
	Description string
	Account     string
}

/*消息代码表
  10000   需要重新登陆
  10001   登录成功
  10002   账号冻结
  10003   账号格式错误！请重新输入。
  10004   账户不存在！请先注册账户。
  10005   代币余额为空！
  10006   用户没有代币发行权限
  10007   查询成功
  10008   账号解冻
  10010   注册成功

  10010   代币已经存在
  10011   代币发布成功
  10012   代币不存在
  10013   代币仓冻结，停止一切交易活动
  10014   没有冻结代币仓的权限
  10015   代币仓冻结失败
  10016   代币增发成功
  10017   代币解冻，交易活动恢复

  10020   余额不足
  10021   支付成功
  10022   转账成功
  10023   支付成功，不能重复支付订单
  10024   退款成功
  10025   提现成功
  10026   充值成功
*/

type Token struct {
	Lock        bool    `json:"Lock"`        /*true为代币冻结*/
	UserId      string  `json:"UserId"`      /*用户id*/
	Name        string  `json:"Name"`        /*代币名称*/
	Symbol      string  `json:"Symbol"`      /*代币符号*/
	TotalSupply float64 `json:"TotalSupply"` /*代币发行量*/
	Description string  `json:"Description"` /*代币简介*/
	Create_time int     `json:"Create_time"` /*创建时间*/
	Id          string  `json:"Id"`
}

type Message struct {
	Status  bool   `json:"Status"`
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

type Balance struct {
	Currency map[string]Currency `json:"Currency"`
}

type Currency struct {
	Name        string  `json:"Name"`
	Symbol      string  `json:"Symbol"`
	TotalSupply float64 `json:"TotalSupply"`
}

type OrderMsg struct {
	OrderNumber  string  `json:"OrderNumber"`  //商品订单号
	OrderName    string  `json:"OrderName"`    //商品名称
	TotalAmount  float64 `json:"TotalAmount"`  //商品总金额
	OrderNote    string  `json:"OrderNote"`    //商品描述
	UserId       string  `json:"UserId"`       //用户ID
	Timestamp    string  `json:"Timestamp"`    //发送请求的时间
	PayStatus    string  `json:"PayStatus"`    //支付状态
	RefundStatus string  `json:"RefundStatus"` //退款状态
	RefundTime   string  `json:"RefundTime"`   //退款时间
	PayType      string  `json:"PayType"`      //支付类型
}

type OrderLists struct {
	OrderLists map[string]OrderMsg
}

type LockToken struct {
	UserId      string  `json:"UserId"`
	Amount      float64 `json:"Amount"`
	Symbol      string  `json:"Symbol"`
	StartTime   int64   `json:"StartTime"`
	EndTime     int64   `json:"EndTime"`
	Status      int     `json:"Status"`      //0代表锁定，1代表解锁
	OrderNumber string  `json:"OrderNumber"` //商品订单号
	PayType     string  `json:"PayType"`     //支付类型
}

type LockTokenList struct {
	LockList map[string]LockToken
}
