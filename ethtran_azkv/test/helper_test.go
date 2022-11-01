package test

import (
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	tran "github.com/cxyzhang0/wallet-go/ethtran_azkv"
	"testing"
)

// address1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
// address2: 0x4A2EBB506da083caC4d61f9305dF8967E595D16b
func TestGetAddressPubKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}

	addrPubKey, addr, err := tran.GetAddressPubKey(keyLabel, _sdk)
	if err != nil {
		t.Errorf("failed to get pub key addr for %+v", keyLabel)
	}

	t.Logf("got pub key addr for %+v: %+v; %s", keyLabel, addrPubKey, addr)
}

func TestDeployContrct(t *testing.T) {

}


