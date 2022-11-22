package test

import (
	"github.com/blockcypher/gobcy"
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

	// To override recipient in req.To
	//req.ToAddr = "2MtorgucUbH9jPXdet913qXwXJJUj4GEPt2"

	signedTxReq, txHash, fromAddr, err := tran.BuildTx(&req, _sdk, networkParams)
	if err != nil {
		t.Logf("fromAddr: %s", fromAddr) // this is dummy just to avoid not used error on fromAddr
		t.Fatalf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", signedTxReq, txHash)

	// broadcast?
	resp, err := tran.BcyAPI.PushTX(signedTxReq)
	if err != nil {
		t.Fatalf("failed to push tx %+v\n %+v", signedTxReq, err)
	}

	t.Logf("broadcasted tx: %+v \n +%v", signedTxReq, resp)

	// webhooks tx_confirmation
	// use pipedream server for receiving blockcypher events.
	// https://pipedream.com/@pingpong/requestbin-p_JZCd9Nz/inspect/2G721AuHAY7lhxXWP2Z6YGYmUpA
	hook, err := tran.BcyAPI.CreateHook(gobcy.Hook{
		Event:         "tx-confirmation",
		Address:       fromAddr,
		Hash:          txHash,
		Confirmations: 3,
		URL:           tran.Conf.Blockcypher.WebhookURL,
		//URL:           "https://eo87b9j94cnw82l.m.pipedream.net/",
	})
	if err != nil {
		t.Fatalf("failed to create webhook for tx-confirmation: +%v", err)
	}

	t.Logf("webhook created for tx-confirmation: %+v", hook)
	// webhooks unconfirmed-tx - for some reason, unconfirmed-tx does not show up
	//hook, err = tran.BcyAPI.CreateHook(gobcy.Hook{
	//	Event:   "unconfirmed-tx",
	//	Address: fromAddr,
	//	Hash:    txHash,
	//	//Confirmations: 6,
	//	URL: tran.Conf.Blockcypher.WebhookURL,
	//	//URL:           "https://eo87b9j94cnw82l.m.pipedream.net/",
	//})
	//if err != nil {
	//	t.Fatalf("failed to create webhook for unconfirmed-tx: +%v", err)
	//}
	//
	//t.Logf("webhook created for unconfirmed_tx: %+v", hook)

	// websocket test failed to connect to "wss://socket.blockcypher.com/v1/btc/test3?token=e905d13ae51748e2b618da1ba4ce0458"
	// it says blockcypher.com certificate failed.
	// so this functionality has not been fully tested
	//startWebsocketClient(t, fromAddr, txHash)
}

/**
P2SH address for sender: 		2MuGNME38DnqsHVjhms2bkVZUDcurFwxjmk
Bech32 address for recipient: 	tb1qcf3p6cdsjflzmcsc286mp4dyrktslv0n4crdyy
*/
func TestBuildMultisigTx(t *testing.T) {
	keyLabel1 := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	keyLabel2 := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Version: "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	keyLabel3 := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	from := []kmssdk.KeyLabel{keyLabel1, keyLabel2, keyLabel3}

	to := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	req := tran.MultisigTxReq{
		From:   from,
		M:      2,
		To:     to,
		Amount: 5000,
	}

	// To override recipient in req.To
	req.ToAddr = "2MtorgucUbH9jPXdet913qXwXJJUj4GEPt2"

	signedTxReq, txHash, fromAddr, err := tran.BuildMultisigTx(&req, _sdk, networkParams)
	if err != nil {
		t.Logf("fromAddr: %s", fromAddr) // this is dummy just to avoid not used error on fromAddr
		t.Fatalf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", signedTxReq, txHash)

	// broadcast
	resp, err := tran.BcyAPI.PushTX(signedTxReq)
	if err != nil {
		t.Fatalf("failed to push tx %+v\n %+v", signedTxReq, err)
	}

	t.Logf("broadcasted tx: %+v \n +%v", signedTxReq, resp)

	// webhooks tx_confirmation
	// use pipedream server for receiving blockcypher events.
	// https://pipedream.com/@pingpong/requestbin-p_JZCd9Nz/inspect/2G721AuHAY7lhxXWP2Z6YGYmUpA
	hook, err := tran.BcyAPI.CreateHook(gobcy.Hook{
		Event:         "tx-confirmation",
		Address:       fromAddr,
		Hash:          txHash,
		Confirmations: 3,
		URL:           tran.Conf.Blockcypher.WebhookURL,
		//URL:           "https://eo87b9j94cnw82l.m.pipedream.net/",
	})
	if err != nil {
		t.Fatalf("failed to create webhook for tx-confirmation: +%v", err)
	}

	t.Logf("webhook created for tx-confirmation: %+v", hook)

}

/**
P2WSH address for sender: 		tb1qg6qjqxy9nv90y2rd5vyp4cfrwevntffpkjs99g68a5gkzlhf23vsdruduz
Bech32 address for recipient: 	tb1qcf3p6cdsjflzmcsc286mp4dyrktslv0n4crdyy
Override with NYDIG P2SH recipient: 2MtorgucUbH9jPXdet913qXwXJJUj4GEPt2
*/
func TestBuildSegWitMultisigTx(t *testing.T) {
	keyLabel1 := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	keyLabel2 := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Version: "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	keyLabel3 := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	from := []kmssdk.KeyLabel{keyLabel1, keyLabel2, keyLabel3}

	to := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	req := tran.MultisigTxReq{
		From:   from,
		M:      2,
		To:     to,
		Amount: 5000,
	}

	// To override recipient in req.To
	req.ToAddr = "2MtorgucUbH9jPXdet913qXwXJJUj4GEPt2"

	signedTxReq, txHash, fromAddr, err := tran.BuildSegWitMultisigTx(&req, _sdk, networkParams)
	if err != nil {
		t.Logf("fromAddr: %s", fromAddr) // this is dummy just to avoid not used error on fromAddr
		t.Fatalf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", signedTxReq, txHash)

	// broadcast
	resp, err := tran.BcyAPI.PushTX(signedTxReq)
	if err != nil {
		t.Fatalf("failed to push tx %+v\n %+v", signedTxReq, err)
	}

	t.Logf("broadcasted tx: %+v \n +%v", signedTxReq, resp)

	// webhooks tx_confirmation
	// use pipedream server for receiving blockcypher events.
	// https://pipedream.com/@pingpong/requestbin-p_JZCd9Nz/inspect/2G721AuHAY7lhxXWP2Z6YGYmUpA
	hook, err := tran.BcyAPI.CreateHook(gobcy.Hook{
		Event:         "tx-confirmation",
		Address:       fromAddr,
		Hash:          txHash,
		Confirmations: 3,
		URL:           tran.Conf.Blockcypher.WebhookURL,
		//URL:           "https://eo87b9j94cnw82l.m.pipedream.net/",
	})
	if err != nil {
		t.Fatalf("failed to create webhook for tx-confirmation: +%v", err)
	}

	t.Logf("webhook created for tx-confirmation: %+v", hook)

}
