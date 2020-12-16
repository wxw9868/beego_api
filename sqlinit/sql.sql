-- connStr := "postgresql://tuser:123456@104.168.163.18:26257/bank?sslmode=require"
-- 应用表
CREATE TABLE token_tenant (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    app_id VARCHAR NOT NULL, -- 应用ID
    app_secret VARCHAR NOT NULL, -- 应用秘钥
    app_name VARCHAR NOT NULL,  -- 应用名称
    app_desc STRING NULL, -- 应用介绍
    user_id VARCHAR NOT NULL, -- 应用标识
    create_time TIMESTAMP not null, -- 创建时间
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    FAMILY "primary" (id, app_id, app_secret, app_name, app_desc, app_token, create_time)
);

-- 应用代币绑定表
CREATE TABLE token_app_banding (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    app_id VARCHAR NOT NULL, -- 应用ID
    manage_id VARCHAR NOT NULL, -- 应用秘钥
    user_id VARCHAR NOT NULL, -- 应用标识
    create_time TIMESTAMP not null, -- 创建时间
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    FAMILY "primary" (id, app_id, manage_id, user_id, create_time)
);


-- DROP TABLE public.token_tenant;
DROP TABLE bank.token_tenant;

-- 用户表
CREATE TABLE token_user (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    tenant_id VARCHAR NOT NULL, -- 应用id
    nickname VARCHAR NULL,  -- 用户名
    mobile INT8 NOT NULL,  -- 手机号
    email VARCHAR NOT NULL,  -- 邮箱
    login_pwd VARCHAR NOT NULL, -- 登录密码
    pay_pwd VARCHAR NOT NULL, -- 支付密码
    user_type INT8 NOT NULL, -- 用户类型：2:企业用户;1:个人用户
    balance DECIMAL NOT NULL DEFAULT 0:::DECIMAL, -- 余额
    balance_lock DECIMAL NOT NULL DEFAULT 0:::DECIMAL, -- 锁定余额
    frozen BOOL NOT NULL DEFAULT false,  -- true为账号冻结
    create_time TIMESTAMP NOT NULL,  -- 创建时间
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    FAMILY "primary" (id, nickname, mobile, email, login_pwd, pay_pwd, user_type, balance, balance_lock, frozen, create_time)
);
-- 添加字段
ALTER TABLE token_user ADD COLUMN tenant_id STRING NOT NULL DEFAULT '0';
ALTER TABLE token_balance ADD COLUMN tenant_id STRING NOT NULL DEFAULT '0';
ALTER TABLE token_balance_action ADD COLUMN tenant_id STRING NOT NULL DEFAULT '0';
ALTER TABLE token_token_manage ADD COLUMN tenant_id STRING NOT NULL DEFAULT '0';
ALTER TABLE token_token_add_log ADD COLUMN tenant_id STRING NOT NULL DEFAULT '0';
SHOW COLUMNS FROM token_user;

-- 用户余额表
CREATE TABLE token_balance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR NOT NULL, -- 应用id
    user_id STRING NOT NULL, -- 用户id
    symbol STRING NOT NULL, -- 代币符号
    token_balance DECIMAL NOT NULL DEFAULT 0, -- 可用余额
    token_balance_lock DECIMAL NOT NULL DEFAULT 0, -- 锁定余额
    create_time TIMESTAMP NOT NULL -- 创建时间
);

-- 余额消费记录表
CREATE TABLE token_balance_action (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR NOT NULL, -- 应用表id
    user_id STRING NOT NULL, -- 用户id
    balance_id STRING NOT NULL , -- 余额id
    amount DECIMAL NOT NULL DEFAULT 0, -- 消费金额
    behavior STRING NOT NULL, -- in:收入,out:支出
    to_user_id STRING NOT NULL, -- 收入、支出的用户id
    create_time TIMESTAMP NOT NULL -- 创建时间
);

-- 代币管理表
CREATE TABLE token_token_manage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR NOT NULL, -- 应用表id
    user_id STRING NOT NULL, -- 用户id
    name STRING NOT NULL,  -- 代币名称
    symbol STRING NOT NULL, -- 代币符号
    total_supply DECIMAL NOT NULL DEFAULT 0, -- 代币发行量
    description STRING,  -- 代币简介
    lock bool NOT NULL, -- true为代币冻结
    create_time TIMESTAMP NOT NULL -- 创建时间
);

-- 代币增发记录
CREATE TABLE token_token_add_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR NOT NULL, -- 应用表id
    manage_id STRING NOT NULL, -- 代币管理表id
    user_id STRING NOT NULL, -- 用户id
    add_supply DECIMAL NOT NULL DEFAULT 0, -- 代币发行量
    create_time TIMESTAMP NOT NULL -- 创建时间
);







