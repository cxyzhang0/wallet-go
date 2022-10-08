package test

import (
	tran "github.com/cxyzhang0/wallet-go/btctran_pkcs11"
	pkcs11sdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"testing"
)

/**
addr1: miDQWKs9G6T3P6VgTwBm1bc9LpccKiiDpU
addr3: mt3NKnJargGTTJvLL8zFS9xUQBXNK56Ftg
*/
func TestGetLegacyAddressPubKey(t *testing.T) {
	pf := "Slot Token 0"
	keyLabel := pkcs11sdk.KeyLabel{
		Prefix:    pf,
		KeyRing:   "WIM-test",
		Key:       "secp256k1-hsm-1",
		Version:   3,
		Algorithm: pkcs11sdk.Secp256k1,
	}

	addrPubKey, addr, err := tran.GetLegacyAddressPubKey(keyLabel, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get pub key addr for %s: %+v", keyLabel.String(), err)
	}

	t.Logf("got pub key addr for %s: %+v; %s", keyLabel.String(), addrPubKey, addr)
}

/**
addr1: tb1q4aja7yqcjv52mypw20pdqsf0msjewf9ac52uxn
addr3: tb1q4sap5cfy9m00dda3zxem9ylycr6gujrsjlklhc
*/
func TestGetBech32AddressPublicKey(t *testing.T) {
	pf := "Slot Token 0"
	keyLabel := pkcs11sdk.KeyLabel{
		Prefix:    pf,
		KeyRing:   "WIM-test",
		Key:       "secp256k1-hsm-1",
		Version:   3,
		Algorithm: pkcs11sdk.Secp256k1,
	}

	addrPubKey, addr, err := tran.GetBech32AddressPublicKey(keyLabel, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get pub key addr for %s: %+v", keyLabel.String(), err)
	}

	t.Logf("got pub key addr for %s: %+v; %s", keyLabel.String(), addrPubKey, addr)
}
