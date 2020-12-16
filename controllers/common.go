package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"log"
	"net"
	"regexp"
	"strings"
)

//数据加密
func DataEncryption(password string) string {
	// DO NOT use this salt value; generate your own random salt. 8 bytes is
	// a good length.
	salt := []byte{0xc7, 0x29, 0xd2, 0x98, 0xb7, 0x7a, 0xcd, 0x7b}

	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(dk)
}

//验证邮箱
func VerifyEmail(str string) bool {
	matched, _ := regexp.MatchString("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", str)
	return matched
}

//验证手机号
func VerifyMobile(str string) bool {
	matched, _ := regexp.MatchString("^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\\d{8}$", str)
	return matched
}

//获取客户端ip
func GetIp() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

//拼接字符串
func StringBuilder(s1, s2 string) string {
	// strings.Builder的0值可以直接使用
	var builder strings.Builder

	// 向builder中写入字符/字符串
	builder.Write([]byte(s1))
	builder.WriteByte(' ')
	builder.WriteString(s2)

	// String() 方法获得拼接的字符串
	return builder.String()
}

//Md5密码加密
func Md5Encrypt(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
