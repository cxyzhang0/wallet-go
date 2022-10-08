package test

import (
	tran "github.com/cxyzhang0/wallet-go/ethtran_pkcs11"
	kmssdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"testing"
)

// address1: 0x389ac41522E3019886ACB003843E62d84FfA70bB
// address3: 0xf9aD19e7a38FaDB98C8A5cC7a14aBcfE80AC657b
func TestGetAddressPubKey(t *testing.T) {
	pf := "Slot Token 0"
	keyLabel := kmssdk.KeyLabel{
		Prefix:    pf,
		KeyRing:   "WIM-test",
		Key:       "secp256k1-hsm-1",
		Version:   3,
		Algorithm: kmssdk.Secp256k1,
	}

	addrPubKey, addr, err := tran.GetAddressPubKey(keyLabel, _sdk)
	if err != nil {
		t.Errorf("failed to get pub key addr for %s", keyLabel.String())
	}

	t.Logf("got pub key addr for %s: %+v; %s", keyLabel.String(), addrPubKey, addr)
}
