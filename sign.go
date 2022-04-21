package main

import (
	"encoding/json"
	"log"
	crypto "sgx-sign-service/crypro"
)

var (
	//Parameters:
	//   - `language`：1中文，2英文。
	//   - `strength`：1弱（12个助记词），2中（18个助记词），3强（24个助记词）。
	language = 2
	strength = uint8(1)
)

// 签名服务接口-方便以后拓展
type SGXSignServe interface {
	// 签名
	Sign(msg []byte) ([]byte, error)
	// 验证
	verify(sign []byte, msg []byte) (bool, error)
}

// 超级链账号
type XuperChainAccount struct {
	// 合约账号
	contractAccount string
	// 地址
	Address string
	// 私钥
	PrivateKey string
	// 公钥
	PublicKey string
	// 助记词
	Mnemonic string
}

// 创建账号
func CreateXuperAccount() (string, error) {
	cli := crypto.GetCryptoClient()
	ecdsaAccount, err := cli.CreateNewAccountWithMnemonic(language, strength)
	if err != nil {
		log.Printf("CreateAccount CreateNewAccountWithMnemonic err: %v", err)
		return "CreateAccount CreateNewAccountWithMnemonic err", err
	}
	// 持久化
	GDB.Add(ecdsaAccount.Address, ecdsaAccount.Mnemonic)
	return ecdsaAccount.Address, nil
}

// 账号是否以及存在
func IsExist(address string) bool {
	return GDB.IsExist(address)
}

// 新建签名服务
func NewXuperchainAccount(address string) SGXSignServe {
	// 从数据库中读取account 的信息
	mnemonic, err := GDB.Query(address)
	if err != nil {
		return nil
	}
	// 恢复
	cryptoClient := crypto.GetCryptoClient()
	ecdsaAccount, err := cryptoClient.RetrieveAccountByMnemonic(mnemonic, language)
	if err != nil {
		return nil
	}

	return &XuperChainAccount{
		Address:    ecdsaAccount.Address,
		PublicKey:  ecdsaAccount.JsonPublicKey,
		PrivateKey: ecdsaAccount.JsonPrivateKey,
		Mnemonic:   ecdsaAccount.Mnemonic,
	}
}

// 签名
func (x *XuperChainAccount) Sign(msg []byte) ([]byte, error) {
	cryptoClient := crypto.GetCryptoClient()
	// 签名
	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(x.PrivateKey)
	if err != nil {
		return nil, err
	}

	sign, err := cryptoClient.SignECDSA(privateKey, msg)
	if err != nil {
		return nil, err
	}

	signInfo := struct {
		PublicKey string `json:"public_key"`
		Sign      []byte `json:"sign"`
	}{
		PublicKey: x.PublicKey,
		Sign:      sign,
	}

	data, err := json.Marshal(&signInfo)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 验证
func (x *XuperChainAccount) verify(sign []byte, msg []byte) (bool, error) {
	cryptoClient := crypto.GetCryptoClient()
	// 公钥
	publickey, err := cryptoClient.GetEcdsaPublicKeyFromJsonStr(x.PublicKey)
	if err != nil {
		return false, err
	}
	signInfo := struct {
		PublicKey string
		Sign      []byte
	}{}

	err = json.Unmarshal(sign, &signInfo)
	if err != nil {
		return false, err
	}
	// 验证
	return cryptoClient.VerifyECDSA(publickey, signInfo.Sign, msg)
}
