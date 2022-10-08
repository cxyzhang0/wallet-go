package test

import (
	tran "github.com/cxyzhang0/wallet-go/ethtran_gcp"
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
	"testing"
)

// address1: 0xaEC11A266C0e4AcaB346Bd7aE4033b3fFB81E401
// address2: 0x7720BBE9bc6201237AbCfD3Fb47317AC51981C71
func TestGetAddressPubKey(t *testing.T) {
	keyLabel := gcpsdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west1",
		KeyRing:  "WIM-test",
		Key:      "secp256k1-hsm-1",
		Version:  2,
	}

	addrPubKey, addr, err := tran.GetAddressPubKey(keyLabel, _sdk)
	if err != nil {
		t.Errorf("failed to get pub key addr for %s", keyLabel.String())
	}

	t.Logf("got pub key addr for %s: %+v; %s", keyLabel.String(), addrPubKey, addr)
}
