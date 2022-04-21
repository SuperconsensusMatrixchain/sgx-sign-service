package main

import (
	"fmt"
	"testing"
)

// V6Avp9KqLfwGUaRFPV27a8VZhwxiYofoU
// test
func Test_SGXSign(t *testing.T) {
	// 创建账号
	addr, err := CreateXuperAccount()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(addr)

	// 签名服务
	sigxserve := NewXuperchainAccount(addr)

	msg := GetDate()

	// sign
	sign, err := sigxserve.Sign(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	// verify
	result, err := sigxserve.verify(sign, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

func GetDate() []byte {
	return []byte("123456")
}
