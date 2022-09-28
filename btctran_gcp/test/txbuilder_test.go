package test

import (
	tran "github.com/cxyzhang0/wallet-go/btctran_gcp"
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
	"testing"
)

/**
addr1: mqFSPcH3XifSDwFaPaMBkHa94A5jCssdD4
addr2: n3E9UG4eP4MaT4LTkmESuECU8xaCyZf7E4
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
		From:   from,
		To:     to,
		Amount: 5000,
	}

	signedTxReq, txHash, err := tran.BuildTx(&req, _sdk, networkParams)
	if err != nil {
		t.Errorf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", signedTxReq, txHash)
}
