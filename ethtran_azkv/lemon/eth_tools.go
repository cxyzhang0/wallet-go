package lemon

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type addrInfo struct {
	privkHex, pubkHex, address string
}

func (ad *addrInfo) toECDSAKey() *ecdsa.PrivateKey {
	k, err := crypto.HexToECDSA(ad.privkHex)
	if err != nil {
		panic(err)
	}
	return k
}

// 生成一个新的地址（随机密钥）
func genNewAddress() (*addrInfo, error) {
	privk, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	// crypto.Fro

	privkBytes := crypto.FromECDSA(privk)
	privkHex := hexutil.Encode(privkBytes)[2:]

	pub := privk.Public()
	pubkECDSA, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		panic("not ok to cast")
	}
	pubkBytes := crypto.FromECDSAPub(pubkECDSA)
	pubHex := hexutil.Encode(pubkBytes)[4:]
	address := crypto.PubkeyToAddress(*pubkECDSA).Hex()
	return &addrInfo{
		privkHex: privkHex, pubkHex: pubHex, address: address,
	}, nil
}
