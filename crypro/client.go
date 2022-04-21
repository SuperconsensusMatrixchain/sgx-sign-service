package crypto

import (
	"github.com/xuperchain/crypto/client/service/base"
	//"github.com/xuperchain/crypto/client/service/gm"
	"github.com/xuperchain/crypto/client/service/xchain"
)

func getInstance() interface{} {
	return &xchain.XchainCryptoClient{}
}

// GetCryptoClient get crypto client
func GetCryptoClient() base.CryptoClient {
	cryptoClient := getInstance().(base.CryptoClient)
	return cryptoClient
}

// GetXchainCryptoClient get xchain crypto client
//func GetXchainCryptoClient() *xchain.XchainCryptoClient {
//	return &xchain.XchainCryptoClient{}
//}
