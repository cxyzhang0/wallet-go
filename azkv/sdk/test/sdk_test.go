package test

import (
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version:   "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	resp, err := _sdk.GenerateKeyPair(keyLabel)
	if err != nil {
		t.Fatalf("failed to create key %+v: %+v", keyLabel, err)
	}

	t.Logf("key pair created %+v: %+v", keyLabel, resp)
}

func TestGetKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:       "secp256k1-hsm-1",
		Version:   "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	res, err := _sdk.GetKey(keyLabel)
	if err != nil {
		t.Fatalf("failed to get key %+v: %+v", keyLabel, err)
	}

	t.Logf("got key for %+v: %+v", keyLabel, res)
}

func TestGetECDSAPublicKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}

	resp, err := _sdk.GetECDSAPublicKey(keyLabel)
	if err != nil {
		t.Fatalf("failed to get the ecdsa public key for %+v: %+v", keyLabel, err)
	}

	t.Logf("got ecdsa public key for %+v: %+v", keyLabel, resp)
}

func TestSignAndVerifySig(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}
	message := "sign me"
	hash, err := kmssdk.SecureHash(message)
	if err != nil {
		t.Fatalf("failed to hash %s", message)
	}

	sigResp, err := _sdk.Sign(keyLabel, hash)
	if err != nil {
		t.Fatalf("failed to sig hash %+v: %+v", hash, err)
	}

	verifyResp, err := _sdk.VerifySig(keyLabel, hash, sigResp.KeyOperationResult.Result)
	if err != nil {
		t.Fatalf("failed to verify hash %+v and sig %+v: %+v", hash, sigResp.KeyOperationResult.Result, err)
	}

	t.Logf("verified: %+v", verifyResp)
}

func TestGetChainSignature(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}
	message := "sign me"
	hash, err := kmssdk.SecureHash(message)
	if err != nil {
		t.Fatalf("failed to hash %s", message)
	}

	signature, err := _sdk.GetChainSignature(keyLabel, hash)
	if err != nil {
		t.Fatalf("failed to get chain signature for key %+v and hash %+v: %+v", keyLabel, hash, err)
	}

	t.Logf("got chain signature for key %+v and hash %+v: %+v", keyLabel, hash, signature)
}