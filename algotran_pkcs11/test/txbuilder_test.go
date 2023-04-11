package test

import (
	"context"
	"encoding/json"
	"github.com/algorand/go-algorand-sdk/future"
	tran "github.com/cxyzhang0/wallet-go/algotran_pkcs11"
	kmssdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestBuildPaymentTx(t *testing.T) {
	pf := "Slot Token 0"
	from := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		1,
		kmssdk.Ed25519,
	}
	fromPubKey, fromAddress, fromAddr, err := tran.GetAddressPubKey(from, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")
	to := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		3,
		kmssdk.Ed25519,
	}
	_, _, toAddr, err := tran.GetAddressPubKey(to, _sdk)
	// overwrite toAddr
	//toAddr = "3DZYSDRDFHORTFPGWPSIKOYIOZYAWI67IQZ6T264OBNJKOB74WRZ2LDPQA"
	//toAddr = "6CCALINGHNFIQXESWJJMO6EAELG7IJF5WY3W37FM23YRL3QN7VYGJHMIVQ"

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	FailOnErr(t, err, "FonSuggestedParams")
	// uncomment the next two lines if overriding suggested fees
	//txParams.FlatFee = true
	//txParams.Fee = 1000

	req := tran.TxReq{
		From:        from,
		FromPubKey:  fromPubKey,
		FromAddress: fromAddress,
		FromAddr:    fromAddr,
		To:          to,
		ToAddr:      toAddr,
		Amount:      1000000000,
		TxParams:    &txParams,
		Note:        []byte("Hello Payment"),
	}

	txId, signedTx, err := tran.BuildPaymentTx(req, _sdk)
	FailOnErr(t, err, "FonBuildPaymentTx")

	t.Logf("\ntx hash: %s\nsigned tx: %+v", txId, signedTx)

	sendResp, err := algodClient.SendRawTransaction(signedTx).Do(context.Background())
	FailOnErr(t, err, "FonSendRawTransaction")

	t.Logf("submitted tx: %s\n", sendResp)

	confirmedTx, err := future.WaitForConfirmation(algodClient, txId, 4, context.Background())
	FailOnErr(t, err, "FonWaitForConfirmation")

	t.Logf("confirmed: %+v\n", confirmedTx)

	txJSON, err := json.MarshalIndent(confirmedTx.Transaction.Txn, "", "\t")
	FailOnErr(t, err, "FonMarshalIndent")

	t.Logf("tx json: %s\n", txJSON)
}

/**
address 1: IXTKWQLXMTOJSRRSYTXRSSRW7CS3YDKZB734FLJMFXKGE6NCNZ3QXY2WLI
address 2: YHFN62XQDRT5HDO5ZFCCZWLJIBBLT5HFVWV4JENM5GN7O6ADUHPH624Y44
address 3: EC7KDBFTC6TFOF4KZWZZFIOJKYM2IGI6O5V7LJVLZK5TMM56UVCTHHATWY

2 of 3 multisig address: 3DZYSDRDFHORTFPGWPSIKOYIOZYAWI67IQZ6T264OBNJKOB74WRZ2LDPQA
3 of 3 multisig address: 6CCALINGHNFIQXESWJJMO6EAELG7IJF5WY3W37FM23YRL3QN7VYGJHMIVQ
*/
func TestBuildMultisigPaymentTx(t *testing.T) {
	var m uint8 = 2
	pf := "Slot Token 0"
	keyLabels := []kmssdk.KeyLabel{
		{
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			1,
			kmssdk.Ed25519,
		},
		{
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			2,
			kmssdk.Ed25519,
		}, {
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			3,
			kmssdk.Ed25519,
		},
	}
	ma, fromAddress, pubKeys, fromAddr, err := tran.GetMultisigAddress(keyLabels, m, _sdk)
	FailOnErr(t, err, "FonGetMultisigAddress")

	to := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		3,
		kmssdk.Ed25519,
	}
	_, _, toAddr, err := tran.GetAddressPubKey(to, _sdk)
	// overwrite toAddr
	//toAddr = "6CCALINGHNFIQXESWJJMO6EAELG7IJF5WY3W37FM23YRL3QN7VYGJHMIVQ"

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	FailOnErr(t, err, "FonSuggestedParams")
	// uncomment the next two lines if overriding suggested fees
	//txParams.FlatFee = true
	//txParams.Fee = 1000

	req := tran.MultisigTxReq{
		From:        keyLabels,
		PubKeys:     pubKeys,
		FromAddress: fromAddress,
		FromAddr:    fromAddr,
		M:           int(m),
		To:          to,
		ToAddr:      toAddr,
		Amount:      10000,
		TxParams:    &txParams,
		Note:        []byte("Hello Multisig Payment"),
	}
	txId, signedTx, err := tran.BuildMultisigPaymentTx(req, *ma, _sdk)
	FailOnErr(t, err, "FonBuildMultisigPaymentTx")

	t.Logf("\ntx hash: %s\nsigned tx: %+v", txId, signedTx)

	sendResp, err := algodClient.SendRawTransaction(signedTx).Do(context.Background())
	FailOnErr(t, err, "FonSendRawTransaction")

	t.Logf("submitted tx: %s\n", sendResp)

	confirmedTx, err := future.WaitForConfirmation(algodClient, txId, 4, context.Background())
	FailOnErr(t, err, "FonWaitForConfirmation")

	t.Logf("confirmed: %+v\n", confirmedTx)

	txJSON, err := json.MarshalIndent(confirmedTx.Transaction.Txn, "", "\t")
	FailOnErr(t, err, "FonMarshalIndent")

	t.Logf("tx json: %s\n", txJSON)
}

/*
{
	"index": 97,
	"params": {
		"clawback": "3DZYSDRDFHORTFPGWPSIKOYIOZYAWI67IQZ6T264OBNJKOB74WRZ2LDPQA",
		"creator": "3DZYSDRDFHORTFPGWPSIKOYIOZYAWI67IQZ6T264OBNJKOB74WRZ2LDPQA",
		"decimals": 2,
		"freeze": "3DZYSDRDFHORTFPGWPSIKOYIOZYAWI67IQZ6T264OBNJKOB74WRZ2LDPQA",
		"manager": "3DZYSDRDFHORTFPGWPSIKOYIOZYAWI67IQZ6T264OBNJKOB74WRZ2LDPQA",
		"metadata-hash": "MHhhY2MwYWUwOGZlNDRmMDY1YWY3ZTlhOTYzMTJjYjU=",
		"name": "WFDC",
		"name-b64": "V0ZEQw==",
		"reserve": "3DZYSDRDFHORTFPGWPSIKOYIOZYAWI67IQZ6T264OBNJKOB74WRZ2LDPQA",
		"total": 1000000000,
		"unit-name": "WFUSD",
		"unit-name-b64": "V0ZVU0Q="
	}
}
*/
func TestBuildMultisigCreateAssetTx(t *testing.T) {
	var m uint8 = 2
	pf := "Slot Token 0"
	keyLabels := []kmssdk.KeyLabel{
		{
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			1,
			kmssdk.Ed25519,
		},
		{
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			2,
			kmssdk.Ed25519,
		}, {
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			3,
			kmssdk.Ed25519,
		},
	}
	ma, fromAddress, pubKeys, fromAddr, err := tran.GetMultisigAddress(keyLabels, m, _sdk)
	FailOnErr(t, err, "FonGetMultisigAddress")

	to := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		3,
		kmssdk.Ed25519,
	}
	_, _, toAddr, err := tran.GetAddressPubKey(to, _sdk)
	// overwrite toAddr
	//toAddr = "6CCALINGHNFIQXESWJJMO6EAELG7IJF5WY3W37FM23YRL3QN7VYGJHMIVQ"

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	FailOnErr(t, err, "FonSuggestedParams")
	// uncomment the next two lines if overriding suggested fees
	//txParams.FlatFee = true
	//txParams.Fee = 1000

	req := tran.MultisigTxReq{
		From:        keyLabels,
		PubKeys:     pubKeys,
		FromAddress: fromAddress,
		FromAddr:    fromAddr,
		M:           int(m),
		To:          to,
		ToAddr:      toAddr,
		Amount:      10000,
		TxParams:    &txParams,
		Note:        []byte("Hello Multisig Asset"),
	}

	param := tran.AssetParam{
		CreatorAddress:  fromAddr,
		AssetName:       "WFDC",
		UnitName:        "WFUSD",
		Total:           1e9,
		Decimals:        2,
		DefaultFrozen:   false,
		URL:             "",
		MetaDataHash:    crypto.Keccak256Hash([]byte("WFDC-WFUSD")).String()[:32],
		ManagerAddress:  fromAddr,
		ReserveAddress:  fromAddr,
		FreezeAddress:   fromAddr,
		ClawbackAddress: fromAddr,
	}
	txId, signedTx, err := tran.BuildMultisigCreateAssetTx(req, *ma, param, _sdk)
	FailOnErr(t, err, "FonBuildMultisigCreateAssetTx")

	t.Logf("\ntx hash: %s\nsigned tx: %+v", txId, signedTx)

	sendResp, err := algodClient.SendRawTransaction(signedTx).Do(context.Background())
	FailOnErr(t, err, "FonSendRawTransaction")

	t.Logf("submitted tx: %s\n", sendResp)

	confirmedTx, err := future.WaitForConfirmation(algodClient, txId, 4, context.Background())
	FailOnErr(t, err, "FonWaitForConfirmation")

	t.Logf("confirmed: %+v\n", confirmedTx)

	assetID := confirmedTx.AssetIndex
	t.Logf("Asset ID: %d\n", assetID)
	printCreatedAsset(assetID, fromAddr, algodClient)
	printAssetHolding(assetID, fromAddr, algodClient)

	//txJSON, err := json.MarshalIndent(confirmedTx.Transaction.Txn, "", "\t")
	//FailOnErr(t, err, "FonMarshalIndent")
	//
	//t.Logf("tx json: %s\n", txJSON)

}

/**
Asset ID: 97
*/
func TestBuildAssetAcceptanceTx(t *testing.T) {
	assetID := uint64(97)
	pf := "Slot Token 0"
	from := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		1,
		kmssdk.Ed25519,
	}
	fromPubKey, fromAddress, fromAddr, err := tran.GetAddressPubKey(from, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	FailOnErr(t, err, "FonSuggestedParams")
	// uncomment the next two lines if overriding suggested fees
	//txParams.FlatFee = true
	//txParams.Fee = 1000

	req := tran.TxReq{
		From:        from,
		FromPubKey:  fromPubKey,
		FromAddress: fromAddress,
		FromAddr:    fromAddr,
		TxParams:    &txParams,
		Note:        []byte("Hello Asset Acceptance"),
	}

	txId, signedTx, err := tran.BuildAssetAcceptanceTx(req, assetID, _sdk)
	FailOnErr(t, err, "FonBuildAssetAcceptanceTx")

	t.Logf("\ntx hash: %s\nsigned tx: %+v", txId, signedTx)

	sendResp, err := algodClient.SendRawTransaction(signedTx).Do(context.Background())
	FailOnErr(t, err, "FonSendRawTransaction")

	t.Logf("submitted tx: %s\n", sendResp)

	confirmedTx, err := future.WaitForConfirmation(algodClient, txId, 4, context.Background())
	FailOnErr(t, err, "FonWaitForConfirmation")

	t.Logf("confirmed: %+v\n", confirmedTx)

	t.Logf("account: %s\n", fromAddr)
	printAssetHolding(assetID, fromAddr, algodClient)

	//txJSON, err := json.MarshalIndent(confirmedTx.Transaction.Txn, "", "\t")
	//FailOnErr(t, err, "FonMarshalIndent")
	//
	//t.Logf("tx json: %s\n", txJSON)
}

func TestGetAllAssetHoldings(t *testing.T) {
	pf := "Slot Token 0"
	from := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		1,
		kmssdk.Ed25519,
	}
	_, _, fromAddr, err := tran.GetAddressPubKey(from, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")

	act, err := algodClient.AccountInformation(fromAddr).Do(context.Background())
	FailOnErr(t, err, "FonAccountInformation")

	t.Logf("account: %s\n", fromAddr)
	for _, assetHolding := range act.Assets {
		prettyPrint(assetHolding)
	}
}

func TestGetAssetHolding(t *testing.T) {
	assetID := uint64(97)
	pf := "Slot Token 0"
	from := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		1,
		kmssdk.Ed25519,
	}
	_, _, fromAddr, err := tran.GetAddressPubKey(from, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")

	algodClient.AccountAssetInformation(fromAddr, assetID).Do(context.Background())
	FailOnErr(t, err, "FonAccountAssetInformation")

	act, err := algodClient.AccountInformation(fromAddr).Do(context.Background())
	FailOnErr(t, err, "FonAccountInformation")

	t.Logf("account: %s\n", fromAddr)
	for _, assetHolding := range act.Assets {
		if assetID == assetHolding.AssetId {
			prettyPrint(assetHolding)
		}
	}
}

func TestGetAllCreatedAssets(t *testing.T) {
	pf := "Slot Token 0"
	m := 2
	/*
		from := kmssdk.KeyLabel{
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			1,
			kmssdk.Ed25519,
		}
		_, _, fromAddr, err := tran.GetAddressPubKey(from, _sdk)
	*/
	keyLabels := []kmssdk.KeyLabel{
		{
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			1,
			kmssdk.Ed25519,
		},
		{
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			2,
			kmssdk.Ed25519,
		}, {
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			3,
			kmssdk.Ed25519,
		},
	}
	_, _, _, fromAddr, err := tran.GetMultisigAddress(keyLabels, uint8(m), _sdk)

	FailOnErr(t, err, "FonGetAddressPubKey")

	act, err := algodClient.AccountInformation(fromAddr).Do(context.Background())
	FailOnErr(t, err, "FonAccountInformation")

	t.Logf("account: %s\n", fromAddr)
	for _, assetHolding := range act.CreatedAssets {
		prettyPrint(assetHolding)
	}
}

func TestBuildMultisigMintAssetTx(t *testing.T) {
	assetID := uint64(97)
	var m uint8 = 2
	pf := "Slot Token 0"
	keyLabels := []kmssdk.KeyLabel{
		{
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			1,
			kmssdk.Ed25519,
		},
		{
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			2,
			kmssdk.Ed25519,
		}, {
			pf,
			"WIM-test",
			"id-ed25519-hsm-1",
			3,
			kmssdk.Ed25519,
		},
	}
	ma, fromAddress, pubKeys, fromAddr, err := tran.GetMultisigAddress(keyLabels, m, _sdk)
	FailOnErr(t, err, "FonGetMultisigAddress")

	to := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		2,
		kmssdk.Ed25519,
	}
	_, _, toAddr, err := tran.GetAddressPubKey(to, _sdk)
	// overwrite toAddr
	//toAddr = "6CCALINGHNFIQXESWJJMO6EAELG7IJF5WY3W37FM23YRL3QN7VYGJHMIVQ"

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	FailOnErr(t, err, "FonSuggestedParams")
	// uncomment the next two lines if overriding suggested fees
	//txParams.FlatFee = true
	//txParams.Fee = 1000

	req := tran.MultisigTxReq{
		From:        keyLabels,
		PubKeys:     pubKeys,
		FromAddress: fromAddress,
		FromAddr:    fromAddr,
		M:           int(m),
		To:          to,
		ToAddr:      toAddr,
		Amount:      20000,
		TxParams:    &txParams,
		Note:        []byte("Hello Mint"),
		AssetID:     assetID,
	}

	txId, signedTx, err := tran.BuildMultisigMintAssetTx(req, *ma, _sdk)
	FailOnErr(t, err, "FonBuildMultisigMintAssetTx")

	t.Logf("\ntx hash: %s\nsigned tx: %+v", txId, signedTx)

	sendResp, err := algodClient.SendRawTransaction(signedTx).Do(context.Background())
	FailOnErr(t, err, "FonSendRawTransaction")

	t.Logf("submitted tx: %s\n", sendResp)

	confirmedTx, err := future.WaitForConfirmation(algodClient, txId, 4, context.Background())
	FailOnErr(t, err, "FonWaitForConfirmation")

	t.Logf("confirmed: %+v\n", confirmedTx)

	t.Logf("account: %s\n", toAddr)
	printAssetHolding(assetID, toAddr, algodClient)

	t.Logf("account: %s\n", fromAddr)
	printAssetHolding(assetID, fromAddr, algodClient)

	//txJSON, err := json.MarshalIndent(confirmedTx.Transaction.Txn, "", "\t")
	//FailOnErr(t, err, "FonMarshalIndent")
	//
	//t.Logf("tx json: %s\n", txJSON)

}

func TestBuildTransferAssetTx(t *testing.T) {
	assetID := uint64(97)
	pf := "Slot Token 0"
	from := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		1,
		kmssdk.Ed25519,
	}
	fromPubKey, fromAddress, fromAddr, err := tran.GetAddressPubKey(from, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")
	to := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		2,
		kmssdk.Ed25519,
	}
	_, _, toAddr, err := tran.GetAddressPubKey(to, _sdk)
	// overwrite toAddr
	//toAddr = "3DZYSDRDFHORTFPGWPSIKOYIOZYAWI67IQZ6T264OBNJKOB74WRZ2LDPQA"
	//toAddr = "6CCALINGHNFIQXESWJJMO6EAELG7IJF5WY3W37FM23YRL3QN7VYGJHMIVQ"

	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	FailOnErr(t, err, "FonSuggestedParams")
	// uncomment the next two lines if overriding suggested fees
	//txParams.FlatFee = true
	//txParams.Fee = 1000

	req := tran.TxReq{
		From:        from,
		FromPubKey:  fromPubKey,
		FromAddress: fromAddress,
		FromAddr:    fromAddr,
		To:          to,
		ToAddr:      toAddr,
		Amount:      100,
		TxParams:    &txParams,
		Note:        []byte("Hello Transfer Asset"),
		AssetID:     assetID,
	}

	txId, signedTx, err := tran.BuildTransferAssetTx(req, _sdk)
	FailOnErr(t, err, "FonBuildPaymentTx")

	t.Logf("\ntx hash: %s\nsigned tx: %+v", txId, signedTx)

	sendResp, err := algodClient.SendRawTransaction(signedTx).Do(context.Background())
	FailOnErr(t, err, "FonSendRawTransaction")

	t.Logf("submitted tx: %s\n", sendResp)

	confirmedTx, err := future.WaitForConfirmation(algodClient, txId, 4, context.Background())
	FailOnErr(t, err, "FonWaitForConfirmation")

	t.Logf("confirmed: %+v\n", confirmedTx)

	t.Logf("account: %s\n", toAddr)
	printAssetHolding(assetID, toAddr, algodClient)

	t.Logf("account: %s\n", fromAddr)
	printAssetHolding(assetID, fromAddr, algodClient)
	//txJSON, err := json.MarshalIndent(confirmedTx.Transaction.Txn, "", "\t")
	//FailOnErr(t, err, "FonMarshalIndent")
	//
	//t.Logf("tx json: %s\n", txJSON)
}
