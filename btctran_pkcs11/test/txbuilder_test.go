package test

import (
	tran "github.com/cxyzhang0/wallet-go/btctran_pkcs11"
	pkcs11sdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"testing"
)

/**
legacy
addr1: miDQWKs9G6T3P6VgTwBm1bc9LpccKiiDpU
addr3: mt3NKnJargGTTJvLL8zFS9xUQBXNK56Ftg
bech32
addr1: tb1q4aja7yqcjv52mypw20pdqsf0msjewf9ac52uxn
addr3: tb1q4sap5cfy9m00dda3zxem9ylycr6gujrsjlklhc
*/
func TestBuildTx(t *testing.T) {
	pf := "Slot Token 0"
	from := pkcs11sdk.KeyLabel{
		Prefix:    pf,
		KeyRing:   "WIM-test",
		Key:       "secp256k1-hsm-1",
		Version:   1,
		Algorithm: pkcs11sdk.Secp256k1,
	}
	to := pkcs11sdk.KeyLabel{
		Prefix:    pf,
		KeyRing:   "WIM-test",
		Key:       "secp256k1-hsm-1",
		Version:   3,
		Algorithm: pkcs11sdk.Secp256k1,
	}
	req := tran.TxReq{
		From:         from,
		IsFromLegacy: true,
		To:           to,
		IsToLegacy:   true,
		Amount:       5000,
	}

	signedTxReq, txHash, err := tran.BuildTx(&req, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", signedTxReq, txHash)

	// broadcast?
	//tran.BcyAPI.PushTX(signedTxReq)
}
