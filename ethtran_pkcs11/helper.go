package ethtran_pkcs11

import (
	"fmt"
	pkcs11sdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func GetAddressPubKey(keyLabel pkcs11sdk.KeyLabel, sdk *pkcs11sdk.SDK) (*common.Address, string, error) { // address pub key, address, error
	pubkey, err := sdk.GetPublicKey(keyLabel)
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
