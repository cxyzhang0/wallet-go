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
func TestGetLegacyAddressPubKey(t *testing.T) {
	keyLabel := gcpsdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west1",
		KeyRing:  "WIM-test",
		Key:      "secp256k1-hsm-1",
		Version:  1,
	}

	addrPubKey, addr, err := tran.GetLegacyAddressPubKey(keyLabel, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get pub key addr for %s", keyLabel.String())
	}

	t.Logf("got pub key addr for %s: %+v; %s", keyLabel.String(), addrPubKey, addr)
}

/**
addr1: tb1qzxvvx3054x34d2ga79sjkhry3zmuqemnetlvza
addr2: tb1qzqz3e8qm96yyz9j08g6zlmm20qctfqeuh4jv8y
*/
func TestGetBech32AddressPublicKey(t *testing.T) {
	keyLabel := gcpsdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west1",
		KeyRing:  "WIM-test",
		Key:      "secp256k1-hsm-1",
		Version:  1,
	}

	addrPubKey, addr, err := tran.GetBech32AddressPublicKey(keyLabel, _sdk, networkParams)
	if err != nil {
		t.Fatalf("failed to get pub key addr for %s", keyLabel.String())
	}

	t.Logf("got pub key addr for %s: %+v; %s", keyLabel.String(), addrPubKey, addr)
}
