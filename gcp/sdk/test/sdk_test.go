package test

import (
	"github.com/cxyzhang0/wallet-go/gcp/sdk"
	"testing"
)

func TestCreateKeyRing(t *testing.T) {
	keyLabel := sdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west2",
		KeyRing:  "WIM-test-2",
	}
	result, err := _sdk.CreateKeyRing(keyLabel)
	if err != nil {
		t.Errorf("failed to create key ring for %s: %+v", keyLabel.String(), err)
	}

	t.Logf("key ring %s created at %v", result.Name, result.CreateTime)
}

func TestGenerateKeyPair(t *testing.T) {
	keyLabel := sdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west2",
		KeyRing:  "WIM-test-2",
		Key:      "secp256p-soft-1",
		//Key:       "secp256k1-hsm-1",
		Algorithm: sdk.Secp256p,
		//Algorithm: sdk.Secp256k1,
	}
	result, err := _sdk.GenerateKeyPair(keyLabel)
	if err != nil {
		t.Errorf("failed to generate key pair for %s: %+v", keyLabel.String(), err)
	}

	t.Logf("key pair %s created at %v", result.Name, result.CreateTime)
}

func TestCreateKeyVersion(t *testing.T) {
	keyLabel := sdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west2",
		KeyRing:  "WIM-test-2",
		Key:      "secp256p-soft-1",
		//Key:       "secp256k1-hsm-1",
		Algorithm: sdk.Secp256p,
		//Algorithm: sdk.Secp256k1,
	}
	result, err := _sdk.CreateKeyVersion(keyLabel)
	if err != nil {
		t.Errorf("failed to create key version %s: %+v", keyLabel.String(), err)
	}

	t.Logf("key vsersion %s created at %v", result.Name, result.CreateTime)
}

/**
Only applicable for non-signing keys
*/
func TestUpdateKeySetPrimary(t *testing.T) {
	keyLabel := sdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west2",
		KeyRing:  "WIM-test-2",
		Key:      "secp256p-soft-1",
		//Key:       "secp256k1-hsm-1",
		Version:   1,
		Algorithm: sdk.Secp256p,
		//Algorithm: sdk.Secp256k1,
	}
	result, err := _sdk.UpdateKeySetPrimary(keyLabel)
	if err != nil {
		t.Errorf("failed to update primary key version %s: %+v", keyLabel.String(), err)
	}

	t.Logf("primary key version for %s: %s; created at %v", result.Name, result.Primary, result.CreateTime)
}

func TestGetPublicKeyPem(t *testing.T) {
	keyLabel := sdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west2",
		KeyRing:  "WIM-test-2",
		//Key:      "secp256p-soft-1",
		Key:     "secp256k1-hsm-1",
		Version: 1,
		//Algorithm: sdk.Secp256p,
		Algorithm: sdk.Secp256k1,
	}
	result, err := _sdk.GetPublicKeyPem(keyLabel)
	if err != nil {
		t.Errorf("failed to get public key pem %s: %+v", keyLabel.String(), err)
	}

	t.Logf("got pem for %s: %s", keyLabel.String(), result)
}

func TestGetECDSAPublicKey(t *testing.T) {
	keyLabel := sdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west2",
		KeyRing:  "WIM-test-2",
		Key:      "secp256p-soft-1",
		//Key:     "secp256k1-hsm-1",
		Version:   1,
		Algorithm: sdk.Secp256p,
		//Algorithm: sdk.Secp256k1,
	}
	result, err := _sdk.GetECDSAPublicKey(keyLabel)
	if err != nil {
		t.Errorf("failed to get ecdsa public key %s: %+v", keyLabel.String(), err)
	}

	t.Logf("got ecdsa public key for %s: %+v", keyLabel.String(), result)
}

func TestGetECDSAPublicKeyForSecp256k1(t *testing.T) {
	keyLabel := sdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west2",
		KeyRing:  "WIM-test-2",
		//Key:      "secp256p-soft-1",
		Key:     "secp256k1-hsm-1",
		Version: 1,
		//Algorithm: sdk.Secp256p,
		Algorithm: sdk.Secp256k1,
	}
	result, err := _sdk.GetECDSAPublicKeyForSecp256k1(keyLabel)
	if err != nil {
		t.Errorf("failed to get ecdsa public key %s: %+v", keyLabel.String(), err)
	}

	t.Logf("got ecdsa public key for %s: %+v", keyLabel.String(), result)
}

func TestSignString(t *testing.T) {
	keyLabel := sdk.KeyLabel{
		Project:  "quantum-pilot-360000",
		Location: "us-west1",
		KeyRing:  "WIM-test",
		Key:      "secp256k1-hsm-1",
		Version:  1,
	}

	message := "sign me"
	sig, hash, err := _sdk.SignString(keyLabel, message)
	if err != nil {
		t.Errorf("failed to sign message: %s", message)
	}

	t.Logf("got sig and hash for %s: %+v; %+v", message, sig, hash)
}
