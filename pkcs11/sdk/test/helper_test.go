package test

import (
	"encoding/hex"
	"fmt"
	"github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"testing"
)

func TestKeyLabel(t *testing.T) {
	pf := "projects/quantum-pilot-360000/locations/us-west1"
	tmp := sdk.KeyLabel{
		pf,
		"WIM-test",
		"secp256k1-hsm-1",
		1,
		sdk.Secp256k1,
	}
	got := tmp.Label()
	want := fmt.Sprintf("%s/keyRings/WIM-test/cryptoKeys/secp256k1-hsm-1/cryptoKeyVersions/1", pf)
	if got != want {
		t.Errorf("got: %s; wanted %s", got, want)
	}

	next := tmp.Next()
	got = next.Label()
	if got == want {
		t.Errorf("got: %s; wanted: %s", got, want)
	}

	gotKeyLabel := tmp
	wantKeyLabel, err := sdk.StringToKeyLabel(want)
	if err != nil {
		t.Errorf("want: %s is wrong format err: %+v", want, err)
	}
	if gotKeyLabel != *wantKeyLabel {
		t.Errorf("got key label: %+v; wanted key label: %+v", gotKeyLabel, *wantKeyLabel)
	}

	gotKeyLabel = next
	wantKeyLabel, err = sdk.StringToKeyLabel(want)
	if err != nil {
		t.Errorf("want: %s is wrong format", want)
	}
	if gotKeyLabel == *wantKeyLabel {
		t.Errorf("got key label: %+v; wanted key label: %+v", gotKeyLabel, *wantKeyLabel)
	}
}

func TestShortKeyLabel(t *testing.T) {
	pf := "Slot Token 0"
	tmp := sdk.KeyLabel{
		pf,
		"WIM-test",
		"secp256k1-hsm-1",
		1,
		sdk.Secp256k1,
	}
	got := tmp.ShortLabel()
	want := fmt.Sprintf("%s/WIM-test/secp256k1-hsm-1/1", pf)
	if got != want {
		t.Errorf("got: %s; wanted %s", got, want)
	}

	next := tmp.Next()
	got = next.ShortLabel()
	if got == want {
		t.Errorf("got: %s; wanted: %s", got, want)
	}

	gotKeyLabel := tmp
	wantKeyLabel, err := sdk.ShortStringToKeyLabel(want)
	if err != nil {
		t.Errorf("want: %s is wrong format err: %+v", want, err)
	}
	if gotKeyLabel != *wantKeyLabel {
		t.Errorf("got key label: %+v; wanted key label: %+v", gotKeyLabel, *wantKeyLabel)
	}

	gotKeyLabel = next
	wantKeyLabel, err = sdk.ShortStringToKeyLabel(want)
	if err != nil {
		t.Errorf("want: %s is wrong format", want)
	}
	if gotKeyLabel == *wantKeyLabel {
		t.Errorf("got key label: %+v; wanted key label: %+v", gotKeyLabel, *wantKeyLabel)
	}
}

func TestGetAllSlots(t *testing.T) {
	p, err := initPKCS11Context(module)
	if err != nil {
		t.Errorf("failed to initiate context: %+v", err)
	}
	sdk.ListAllSlots(p)
}

func TestGetSlot(t *testing.T) {
	p, err := initPKCS11Context(module)
	if err != nil {
		t.Errorf("failed to initiate context: %+v", err)
	}
	slotId, err := sdk.GetSlot(p, tokenLabel)
	if err != nil {
		t.Errorf("failed to get slot for label %s: %+v", tokenLabel, err)
	}

	t.Logf("got slot id: %d for slot label: %s", slotId, tokenLabel)
}

func TestDigestSHA256(t *testing.T) {
	message := "test message"
	want := "3f0a377ba0a4a460ecb616f6507ce0d8cfa3e704025d4fda3ed0c5ca05468728"
	//want, err := hex.DecodeString(want)
	//if err != nil {
	//	t.Errorf("failed to decode hex string %s: %+v", want, err)
	//}

	hash, err := sdk.DigestSHA256(_sdk.P, _sdk.Session, message)
	if err != nil {
		t.Errorf("failed to hash message %s: %+v", message, err)
	}

	got := hex.EncodeToString(hash)
	if want != got {
		t.Errorf("got: %s; wanted: %s", got, want)
	}
}
