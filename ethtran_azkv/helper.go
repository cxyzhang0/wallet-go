package ethtran_azkv

import (
	"fmt"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func GetAddressPubKey(keyLabel kmssdk.KeyLabel, sdk *kmssdk.SDK) (*common.Address, string, error) { // address pub key, address, error
	pubkey, err := sdk.GetECDSAPublicKey(keyLabel)
	if err != nil {
		return nil, "", err
	}

	addrPubKey := crypto.PubkeyToAddress(*pubkey)

	addr := addrPubKey.Hex()

	return &addrPubKey, addr, nil
}

func StringToBigInt(str string) (*big.Int, error) {
	bigFloat, _, err := big.ParseFloat(str, 10, 0, big.ToZero)
	if err != nil {
		return nil, fmt.Errorf("failed to pars big int: %+v", err)
	}

	bigInt, _ := bigFloat.Int(nil)

	return bigInt, nil
}
