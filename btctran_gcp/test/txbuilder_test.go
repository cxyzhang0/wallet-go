package test

import (
	tran "github.com/cxyzhang0/wallet-go/btctran_gcp"
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
	"testing"
)

/**
legacy
addr1: mqFSPcH3XifSDwFaPaMBkHa94A5jCssdD4
addr2: n3E9UG4eP4MaT4LTkmESuECU8xaCyZf7E4
bech32
addr1: tb1qzxvvx3054x34d2ga79sjkhry3zmuqemnetlvza
addr2: tb1qzqz3e8qm96yyz9j08g6zlmm20qctfqeuh4jv8y
*/
func TestBuildTx(t *testing.T) {
	from := gcpsdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west1",
		KeyRing:  "WIM-test",
		Key:      "secp256k1-hsm-1",
		Version:  1,
	}
	to := gcpsdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west1",
		KeyRing:  "WIM-test",
		Key:      "secp256k1-hsm-1",
		Version:  2,
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
	//tran.BcyAPI.PushTX(signedTxReq)
}
