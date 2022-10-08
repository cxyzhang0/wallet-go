package ethtran_gcp

import (
	"fmt"
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func GetAddressPubKey(keyLabel gcpsdk.KeyLabel, sdk *gcpsdk.SDK) (*common.Address, string, error) { // address pub key, address, error
	pubkey, err := sdk.GetECDSAPublicKeyForSecp256k1(keyLabel)
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
