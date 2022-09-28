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
func TestGetAddressPubKey(t *testing.T) {
	keyLabel := gcpsdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west1",
		KeyRing:  "WIM-test",
		Key:      "secp256k1-hsm-1",
		Version:  1,
	}

	addrPubKey, addr, err := tran.GetAddressPubKey(keyLabel, _sdk, networkParams)
	if err != nil {
		t.Errorf("failed to get pub key addr for %s", keyLabel.String())
	}

	t.Logf("got pub key addr for %s: %+v; %s", keyLabel.String(), addrPubKey, addr)
}
