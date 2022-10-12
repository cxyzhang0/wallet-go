package test

import (
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	tran "github.com/cxyzhang0/wallet-go/btctran_azkv"
	"testing"
)

/**
legacy
addr1: mhE4NMciTWxnDyc2rsBSLsALjbRvdQfNaM
addr2: mywkQeemfaVX9eNrcELxSNYSxkFPccRipn
bech32
addr1: tb1q3twchl2qkwjdd8nprf4x75qcqhas2c8kgq54nu
addr2: tb1qefz780p9922jf7fyj8yet3pa8wukx8xrmn45af
*/
func TestBuildTx(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}
	to := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}
	req := tran.TxReq{
		From:         from,
		IsFromLegacy: false,
		To:           to,
		IsToLegacy:   false,
		Amount:       5000,
	}

	signedTxReq, txHash, err := tran.BuildTx(&req, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", signedTxReq, txHash)

	// broadcast?
	resp, err := tran.BcyAPI.PushTX(signedTxReq)
	if err != nil {
		t.Fatalf("failed to push tx %+v\n %+v", signedTxReq, err)
	}

	t.Logf("pushed tx %+v \n +%v", signedTxReq, resp)
}
