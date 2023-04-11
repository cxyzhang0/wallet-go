package test

import (
	"github.com/btcsuite/btcd/txscript"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	tran "github.com/cxyzhang0/wallet-go/btctran_azkv"
	"testing"
)

/**
addr1: mhE4NMciTWxnDyc2rsBSLsALjbRvdQfNaM
addr2: mywkQeemfaVX9eNrcELxSNYSxkFPccRipn
*/
func TestGetLegacyAddressPubKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}

	addrPubKey, addr, err := tran.GetLegacyAddressPubKey(keyLabel, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get pub key addr for %+v: %+v", keyLabel, err)
	}

	t.Logf("got pub key addr for %+v: %+v; %s", keyLabel, addrPubKey, addr)
}

/**
soft
addr1: tb1q2222hcm4rv40t2gwee373gm5699dxcelxvda36
hsm
addr1: tb1q3twchl2qkwjdd8nprf4x75qcqhas2c8kgq54nu
addr2: tb1qefz780p9922jf7fyj8yet3pa8wukx8xrmn45af
addr3: tb1qcf3p6cdsjflzmcsc286mp4dyrktslv0n4crdyy
*/
func TestGetBech32AddressPublicKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     softKeyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version: "0ff7adfdbe0a4b69881c4dac6b0f81f4",
		//Version:   "485248105ef047aaa1f33cf0baaf9a03",
		Algorithm: kmssdk.Secp256k1,
	}

	/*
		keyLabel := kmssdk.KeyLabel{
			Key:     hsmKeyName,
			Version: "cb848fb15e3a40b49bc41cbe957ea438",
			//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
			//Version:   "b6aec266b6a147f7a1c40fe842504650",
			Algorithm: kmssdk.Secp256k1,
		}
	*/
	addrPubKey, addr, err := tran.GetBech32AddressPublicKey(keyLabel, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get pub key addr for %+v: %+v", keyLabel, err)
	}

	t.Logf("got pub key addr for %+v: %+v; %s", keyLabel, addrPubKey, addr)
}

/**
2 of 3 multisig P2SH addr: 2MuGNME38DnqsHVjhms2bkVZUDcurFwxjmk
script:
2 032088a41f1a35d4c29906e7c2bbbada0a96c04fb5b4c67add0d6534a9569cf3d3 02e245f74cf1d763c23eef45ff068db34ae93a6dd914be5df5d87d37eb84eed6e3 0330c2b8caf05320a939bd698531d6b5ccd6bd6ad48cd3d0f122bbd24c22e2f8e9 3 OP_CHECKMULTISIG
*/
func TestGetP2SHRedeemScript(t *testing.T) {
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

	redeemScript, scriptPubKey, addr, err := tran.GetP2SHRedeemScript(keyLabel1, keyLabel2, keyLabel3, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get P2SH redeem script: %+v", err)
	}

	redeemStr, err := txscript.DisasmString(redeemScript)
	if err != nil {
		t.Fatalf("failed to disassemble redeem script %+v: %+v", redeemScript, err)
	}

	scriptStr, err := txscript.DisasmString(scriptPubKey)
	if err != nil {
		t.Fatalf("failed to disassemble script pub key %+v: %+v", scriptPubKey, err)
	}

	t.Logf("got P2SH addr:\n%s\nredeem script:\n%s\nscript pub key:\n%s ", addr, redeemStr, scriptStr)
}

/**
// same result as above
2 of 3 multisig P2SH addr:
2MuGNME38DnqsHVjhms2bkVZUDcurFwxjmk
script:
2 032088a41f1a35d4c29906e7c2bbbada0a96c04fb5b4c67add0d6534a9569cf3d3 02e245f74cf1d763c23eef45ff068db34ae93a6dd914be5df5d87d37eb84eed6e3 0330c2b8caf05320a939bd698531d6b5ccd6bd6ad48cd3d0f122bbd24c22e2f8e9 3 OP_CHECKMULTISIG
*/
func TestGetMultisigRedeemScript(t *testing.T) {
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

	keyLabels := []kmssdk.KeyLabel{keyLabel1, keyLabel2, keyLabel3}
	redeemScript, scriptPubKey, addr, err := tran.GetMultisigRedeemScript(keyLabels, 2, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get P2SH redeem script: %+v", err)
	}

	redeemStr, err := txscript.DisasmString(redeemScript)
	if err != nil {
		t.Fatalf("failed to disassemble redeem script %+v: %+v", redeemScript, err)
	}

	scriptStr, err := txscript.DisasmString(scriptPubKey)
	if err != nil {
		t.Fatalf("failed to disassemble script pub key%+v: %+v", scriptPubKey, err)
	}
	t.Logf("got P2SH addr:\n%s\nredeem script:\n%s\nscript pub key:\n%s", addr, redeemStr, scriptStr)
}

/*
// 2 of 3 multisig P2WSH addr:
tb1qg6qjqxy9nv90y2rd5vyp4cfrwevntffpkjs99g68a5gkzlhf23vsdruduz
script:
2 032088a41f1a35d4c29906e7c2bbbada0a96c04fb5b4c67add0d6534a9569cf3d3 02e245f74cf1d763c23eef45ff068db34ae93a6dd914be5df5d87d37eb84eed6e3 0330c2b8caf05320a939bd698531d6b5ccd6bd6ad48cd3d0f122bbd24c22e2f8e9 3 OP_CHECKMULTISIG
*/
func TestGetSegWitMultisigRedeemScript(t *testing.T) {
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

	keyLabels := []kmssdk.KeyLabel{keyLabel1, keyLabel2, keyLabel3}
	redeemScript, scriptPubKey, addr, err := tran.GetSegWitMultisigRedeemScript(keyLabels, 2, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get P2SH redeem script: %+v", err)
	}

	redeemStr, err := txscript.DisasmString(redeemScript)
	if err != nil {
		t.Fatalf("failed to disassemble redeem script %+v: %+v", redeemScript, err)
	}

	scriptStr, err := txscript.DisasmString(scriptPubKey)
	if err != nil {
		t.Fatalf("failed to disassemble script pub key%+v: %+v", scriptPubKey, err)
	}

	t.Logf("got Multisig Bech32 addr:\n%s\nredeem script:\n%s\nscript pub key:\n%s", addr, redeemStr, scriptStr)
}
