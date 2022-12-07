package test

import (
	tran "github.com/cxyzhang0/wallet-go/avaxtran_azkv"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"testing"
)

/**
hsm
address1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
address2: 0x4A2EBB506da083caC4d61f9305dF8967E595D16b
address3: 0x4357fB73aF4359D2ec2dc449B90D73495F7794DD
soft
address1: 0xBfDdb90EFD9323db8E9cB136cE27E5b38fE29Da8
address2: 0xCf8A60d998F1Cc6c9346A3feb8464c5a8aCED274
address3: 0x6c0a79d56c9C861Dc8907469728e6Bcc4aEe5fbB
*/
func TestGetAddressPubKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     softKeyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version: "0ff7adfdbe0a4b69881c4dac6b0f81f4",
		//Version:   "485248105ef047aaa1f33cf0baaf9a03",
		Algorithm: kmssdk.Secp256k1,
	}
	/*
		keyLabel := kmssdk.KeyLabel{
			Key:     "secp256k1-hsm-1",
			Version: "cb848fb15e3a40b49bc41cbe957ea438",
			//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
			//Version:   "b6aec266b6a147f7a1c40fe842504650",
			Algorithm: kmssdk.Secp256k1,
		}
	*/
	addrPubKey, addr, err := tran.GetAddressPubKey(keyLabel, _sdk)
	if err != nil {
		t.Errorf("failed to get pub key addr for %+v", keyLabel)
	}

	t.Logf("got pub key addr for %+v: %+v; %s", keyLabel, addrPubKey, addr)
}
