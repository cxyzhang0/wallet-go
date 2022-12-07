package test

import (
	"github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	//pf := "us01hsm01test.virtucrypt.com:592"
	pf := "Slot Token 0"
	//pf := "projects/quantum-pilot-360000/locations/us-west1"
	keyLabel := sdk.KeyLabel{
		pf,
		"WIM-test",
		"ed25519-hsm-1",
		1,
		sdk.Ed25519,
	}
	/*
		keyLabel := sdk.KeyLabel{
			pf,
			"WIM-test",
			"secp256k1-hsm-1",
			5,
			sdk.Secp256k1,
		}
	*/
	persistent := false
	pbk, prvk, err := _sdk.GenerateKeyPair(keyLabel, persistent)
	FailOnErr(t, err, "FonGenerateKeyPair")
	//if err != nil {
	//	t.Fatalf("failed to generate key pair - curve: %s; keyLabel string: %s; persistent: %v: %+v", keyLabel.Algorithm.String(), keyLabel.ShortLabel(), persistent, err)
	//}

	t.Logf("generated key pair - pbk: %d; prvk: %d; curve: %s; keyLabel string: %s; persistent: %v", pbk, prvk, keyLabel.Algorithm.String(), keyLabel.ShortLabel(), persistent)
}

func TestSignAndVerify(t *testing.T) {
	//pf := "us01hsm01test.virtucrypt.com:592"
	pf := "Slot Token 0"
	//pf := "projects/quantum-pilot-360000/locations/us-west1"
	keyLabel := sdk.KeyLabel{
		pf,
		"WIM-test",
		"secp256k1-hsm-1",
		1,
		sdk.Secp256k1,
	}
	message := "sign me"
	hash, err := sdk.SecureHash(message)
	//hash, err := sdk.DigestSHA256(_sdk.P, _sdk.Session, message)
	if err != nil {
		t.Fatalf("failed to hash message %s: %+v", message, err)
	}

	sig, err := _sdk.Sign(keyLabel, hash)
	if err != nil {
		t.Fatalf("failed to sign hash %+v: %+v", hash, err)
	}

	if err := _sdk.VerifySig(keyLabel, hash, sig); err != nil {
		t.Fatalf("failed to verify sig %+v: %+v", sig, err)
	}

	t.Logf("signed and verifyed: %+v", sig)
}

func TestFindAllECKeys(t *testing.T) {
	objs, err := _sdk.GetAllECKeys()
	if err != nil {
		t.Errorf("failed to get all EC keys - error: %+v", err)
	}

	t.Logf("got all EC keys objs %+v", objs)
}

func TestFindAllPrivateKeys(t *testing.T) {
	objs, err := _sdk.GetAllPrivateKeys()
	if err != nil {
		t.Errorf("failed to get all private keys - error: %+v", err)
	}

	t.Logf("got all private keys objs %+v", objs)
}

func TestGetPrivateKeyHandle(t *testing.T) {
	//pf := "us01hsm01test.virtucrypt.com:592"
	pf := "Slot Token 0"
	//pf := "projects/quantum-pilot-360000/locations/us-west1"
	keyLabel := sdk.KeyLabel{
		pf,
		"WIM-test",
		"secp256k1-hsm-1",
		5,
		sdk.Secp256k1,
	}
	obj, err := _sdk.GetPrivateKeyHandle(keyLabel)
	if err != nil {
		t.Fatalf("failed to get private key for %s; error: %+v", keyLabel.ShortLabel(), err)
	}

	t.Logf("got private key obj %d for %s", obj, keyLabel.ShortLabel())
}

func TestGetPublicKeyHandle(t *testing.T) {
	//pf := "us01hsm01test.virtucrypt.com:592"
	pf := "Slot Token 0"
	//pf := "projects/quantum-pilot-360000/locations/us-west1"
	keyLabel := sdk.KeyLabel{
		pf,
		"WIM-test",
		"secp256k1-hsm-1",
		5,
		sdk.Secp256k1,
	}
	obj, err := _sdk.GetPublicKeyHandle(keyLabel)
	if err != nil {
		t.Fatalf("failed to get public key for %s; error: %+v", keyLabel.ShortLabel(), err)
	}

	t.Logf("got public key obj %d for %s", obj, keyLabel.ShortLabel())
}

func TestGetPublicKeyECPoint(t *testing.T) {
	//pf := "us01hsm01test.virtucrypt.com:592"
	pf := "Slot Token 0"
	//pf := "projects/quantum-pilot-360000/locations/us-west1"
	keyLabel := sdk.KeyLabel{
		pf,
		"WIM-test",
		"secp256k1-hsm-1",
		5,
		sdk.Secp256k1,
	}

	attr, err := _sdk.GetPublicKeyECPoint(keyLabel)
	if err != nil {
		t.Fatalf("failed to get public key attribute for %s: %+v", keyLabel.ShortLabel(), err)
	}

	t.Logf("got public key attr %+v with length %d for %s", *attr, len(attr.Value), keyLabel.ShortLabel())
}

func TestGetPublicKey(t *testing.T) {
	//pf := "us01hsm01test.virtucrypt.com:592"
	pf := "Slot Token 0"
	//pf := "projects/quantum-pilot-360000/locations/us-west1"
	keyLabel := sdk.KeyLabel{
		pf,
		"WIM-test",
		"secp256k1-hsm-1",
		5,
		sdk.Secp256k1,
	}

	pubKey, err := _sdk.GetECDSAPublicKey(keyLabel)
	if err != nil {
		t.Fatalf("failed to get ecdsa public key for %s: %+v", keyLabel.ShortLabel(), err)
	}

	t.Logf("got ecdsa public key for %s; \n%+v", keyLabel.ShortLabel(), pubKey)
}

func TestGetPublicKeyAttrFromPrivateKey(t *testing.T) {
	//pf := "us01hsm01test.virtucrypt.com:592"
	pf := "Slot Token 0"
	//pf := "projects/quantum-pilot-360000/locations/us-west1"
	keyLabel := sdk.KeyLabel{
		pf,
		"WIM-test",
		"secp256k1-hsm-1",
		5,
		sdk.Secp256k1,
	}

	attr, err := _sdk.GetPublicKeyAttrFromPrivateKey(keyLabel)
	if err != nil {
		t.Fatalf("failed to get public key attribute for %s: %+v", keyLabel.ShortLabel(), err)
	}

	t.Logf("got public key attr %+v with length %d for %s", *attr, len(attr.Value), keyLabel.ShortLabel())
}

// It is not working.
func TestGetPubKeyBySig(t *testing.T) {
	pf := "Slot Token 0"
	//pf := "projects/quantum-pilot-360000/locations/us-west1"
	keyLabel := sdk.KeyLabel{
		pf,
		"WIM-test",
		"secp256k1-hsm-1",
		1,
		sdk.Secp256k1,
	}

	pubKey, err := _sdk.GetPubKeyBySig(keyLabel)
	if err != nil {
		t.Errorf("failed to get public key for %s; error: %+v", keyLabel.ShortLabel(), err)
	}

	t.Logf("got public key obj %+v for %s", pubKey, keyLabel.ShortLabel())
}
