package test

import (
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
		t.Fatalf("failed to get pub key addr for %+v", keyLabel)
	}

	t.Logf("got pub key addr for %+v: %+v; %s", keyLabel, addrPubKey, addr)
}

/**
addr1: tb1q3twchl2qkwjdd8nprf4x75qcqhas2c8kgq54nu
addr2: tb1qefz780p9922jf7fyj8yet3pa8wukx8xrmn45af
*/
func TestGetBech32AddressPublicKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}

	addrPubKey, addr, err := tran.GetBech32AddressPublicKey(keyLabel, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get pub key addr for %+v", keyLabel)
	}

	t.Logf("got pub key addr for %+v: %+v; %s", keyLabel, addrPubKey, addr)
}
