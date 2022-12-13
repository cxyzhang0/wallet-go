package test

import (
	tran "github.com/cxyzhang0/wallet-go/algotran_pkcs11"
	kmssdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"testing"
)

/**
address 1: IXTKWQLXMTOJSRRSYTXRSSRW7CS3YDKZB734FLJMFXKGE6NCNZ3QXY2WLI
address 2: YHFN62XQDRT5HDO5ZFCCZWLJIBBLT5HFVWV4JENM5GN7O6ADUHPH624Y44
address 3: EC7KDBFTC6TFOF4KZWZZFIOJKYM2IGI6O5V7LJVLZK5TMM56UVCTHHATWY

./sandbox goal account list
./sandbox goal clerk send -a 1000000000000000 -f RSNHVAMLPZPY2BQ64YEHZUVPR5P5KKTKJNLDJA5KWTYN7VOMUI7V3S7QIU -t IXTKWQLXMTOJSRRSYTXRSSRW7CS3YDKZB734FLJMFXKGE6NCNZ3QXY2WLI
./sandbox goal account balance -a IXTKWQLXMTOJSRRSYTXRSSRW7CS3YDKZB734FLJMFXKGE6NCNZ3QXY2WLI

*/
func TestGetAddress(t *testing.T) {
	pf := "Slot Token 0"
	keyLabel := kmssdk.KeyLabel{
		pf,
		"WIM-test",
		"id-ed25519-hsm-1",
		3,
		kmssdk.Ed25519,
	}
	pubKey, address, addStr, err := tran.GetAddressPubKey(keyLabel, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")

	t.Logf("\npubKey: %+v\naddress byptes: %+v\naddress: %s", pubKey, address, addStr)
}

/**
2 of 3 multisig address: 3DZYSDRDFHORTFPGWPSIKOYIOZYAWI67IQZ6T264OBNJKOB74WRZ2LDPQA
3 of 3 multisig address: 6CCALINGHNFIQXESWJJMO6EAELG7IJF5WY3W37FM23YRL3QN7VYGJHMIVQ
*/
func TestGetMultisigAddress(t *testing.T) {
	pf := "Slot Token 0"
	var m uint8 = 2
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
	ma, address, pubKeys, addr, err := tran.GetMultisigAddress(keyLabels, m, _sdk)
	FailOnErr(t, err, "FonGetMultisigAddress")

	t.Logf("\naddress bytes: %+v\naddress: %s\npub keys: %+v\nmultisig accont: %+v", address, addr, pubKeys, ma)
}
