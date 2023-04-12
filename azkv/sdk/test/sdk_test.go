package test

import (
	"encoding/hex"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"testing"
)

/*
*
key: secp256k1-soft-1 versions
0eab9a0cc2e84018be05f90e5d914142
0ff7adfdbe0a4b69881c4dac6b0f81f4
485248105ef047aaa1f33cf0baaf9a03
key: secp256k1-hsm-1 versions
cb848fb15e3a40b49bc41cbe957ea438
0179a6204ed7491ea5b27a87b541d5cb
b6aec266b6a147f7a1c40fe842504650
*/
func TestGenerateKeyPair(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:       keyName,
		Algorithm: kmssdk.Secp256k1,
	}

	resp, err := _sdk.GenerateKeyPair(keyLabel)
	FailOnErr(t, err, "FonGenerateKeyPair")

	t.Logf("key pair created %+v: %+v", keyLabel, resp)
}

func TestGetKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version:   "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	res, err := _sdk.GetKey(keyLabel)
	FailOnErr(t, err, "FonGetKey")

	t.Logf("got key for %+v: %+v", keyLabel, res)
}

func TestGetECDSAPublicKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	resp, err := _sdk.GetECDSAPublicKey(keyLabel)
	FailOnErr(t, err, "FonGetECDSAPublicKey")

	t.Logf("got ecdsa public key for %+v: %+v", keyLabel, resp)
}

func TestGetSECP256K1PublicKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	resp, err := _sdk.GetSECP256K1PublicKey(keyLabel)
	FailOnErr(t, err, "FonGetSECP256K1PublicKey")

	t.Logf("got secp256k1 public key for %+v: %+v", keyLabel, resp)
}

func TestGetCosmosSECP256K1PubKey(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	resp, err := _sdk.GetCosmosSECP256K1PubKey(keyLabel)
	FailOnErr(t, err, "FonGetCosmosSECP256K1PubKey")

	t.Logf("got cosmos secp256k1 pubkey for %+v: %+v", keyLabel, resp)
}

func TestSignAndVerifySig(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version:   "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}
	message := "sign me"
	hash, err := kmssdk.SecureHash(message)
	FailOnErr(t, err, "FonSecureHash")

	sigResp, err := _sdk.Sign(keyLabel, hash)
	FailOnErr(t, err, "FonSign")

	verifyResp, err := _sdk.VerifySig(keyLabel, hash, sigResp.KeyOperationResult.Result)
	FailOnErr(t, err, "FonVerifySig")

	require.True(t, *verifyResp.Value)
	//t.Logf("verified: %t", *verifyResp.Value)
}

func TestSignAndVerifySig_negative(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version:   "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}
	message := "sign me"
	hash, err := kmssdk.SecureHash(message)
	FailOnErr(t, err, "FonSecureHash")

	message1 := "sign me 1"
	hash1, err := kmssdk.SecureHash(message1)
	FailOnErr(t, err, "FonSecureHash")

	sigResp, err := _sdk.Sign(keyLabel, hash)
	FailOnErr(t, err, "FonSign")

	verifyResp, err := _sdk.VerifySig(keyLabel, hash1, sigResp.KeyOperationResult.Result)
	FailOnErr(t, err, "FonVerifySig")

	require.False(t, *verifyResp.Value)
	//t.Logf("verified: %t", *verifyResp.Value)
}

// Signature directly from HSM - malleability
func TestCosmosSignAndVerifySig_HalfFail(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version:   "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	msg := []byte("sign me")

	// both SecureHashByteArray and crypto.Sha256 work
	//hash, err := kmssdk.SecureHashByteArray(msg)
	//FailOnErr(t, err, "FonSecureHash")
	hash := crypto.Sha256(msg)
	//t.Logf("hash: %v", hash)

	sigResp, err := _sdk.Sign(keyLabel, hash)
	//sigResp, err := _sdk.Sign(keyLabel, hash)
	//sigResp, err := _sdk.Sign(keyLabel, hash)
	FailOnErr(t, err, "FonSign")

	pubKey, err := _sdk.GetCosmosSECP256K1PubKey(keyLabel)
	FailOnErr(t, err, "FonGetCosmosSECP256K1PubKey")

	verified := pubKey.VerifySignature(msg, sigResp.Result)

	t.Logf("sig:     %s", hex.EncodeToString(sigResp.Result))
	t.Logf("pub key: %s", hex.EncodeToString(pubKey.GetKey()))
	//  pk: 02e9f0f33d79d0459b22ecb253cb3ec52c4ad12f4419fbea03ec825cc74fb3b693
	require.True(t, verified)
	//t.Logf("verified: %t", verified)
}

// Signature directly from HSM - malleability
func TestCosmosSignAndVerifySig_HalfFailLoop(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version:   "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	msg := []byte("sign me")

	// both SecureHashByteArray and crypto.Sha256 work
	//hash, err := kmssdk.SecureHashByteArray(msg)
	//FailOnErr(t, err, "FonSecureHash")
	hash := crypto.Sha256(msg)
	//t.Logf("hash: %v", hash)
	var successes int
	var runs = 100
	for i := 0; i < runs; i++ {
		sigResp, err := _sdk.Sign(keyLabel, hash)
		//sigResp, err := _sdk.Sign(keyLabel, hash)
		FailOnErr(t, err, "FonSign")

		pubKey, err := _sdk.GetCosmosSECP256K1PubKey(keyLabel)
		FailOnErr(t, err, "FonGetCosmosSECP256K1PubKey")

		verified := pubKey.VerifySignature(msg, sigResp.Result)

		//t.Logf("sig:     %s", hex.EncodeToString(sigResp.Result))
		//t.Logf("pub key: %s", hex.EncodeToString(pubKey.GetKey()))
		// sig: c01a8eaf89e5722f15b6acf7da52e638ebd98ce6f81a806ee92e54ed7fec745434215d29438eefde9d04186a3d6fc8b36ff4ec02b3ae45ef74e20466271b3030
		//	pk:	04e9f0f33d79d0459b22ecb253cb3ec52c4ad12f4419fbea03ec825cc74fb3b693dbd35affe186096eacabaa5819961ca3fa42c255ac580fdacec23e784dbdd59a
		//  pk: 02e9f0f33d79d0459b22ecb253cb3ec52c4ad12f4419fbea03ec825cc74fb3b693
		//require.True(t, verified)
		if verified {
			successes++
		}
		//t.Logf("verified: %t", verified)
	}
	t.Logf("%d of %d", successes, runs)
	require.Equal(t, runs, successes)

}

func TestCosmosSignAndVerifySig(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version:   "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	msg := []byte("sign me")

	// both SecureHashByteArray and crypto.Sha256 work
	//hash, err := kmssdk.SecureHashByteArray(msg)
	//FailOnErr(t, err, "FonSecureHash")
	hash := crypto.Sha256(msg)
	//t.Logf("hash: %v", hash)

	sig, err := _sdk.GetCosmosChainSignature(keyLabel, hash)
	//sigResp, err := _sdk.Sign(keyLabel, hash)
	FailOnErr(t, err, "FonSign")

	pubKey, err := _sdk.GetCosmosSECP256K1PubKey(keyLabel)
	FailOnErr(t, err, "FonGetCosmosSECP256K1PubKey")

	verified := pubKey.VerifySignature(msg, sig)

	t.Logf("sig:     %s", hex.EncodeToString(sig))
	t.Logf("pub key: %s", hex.EncodeToString(pubKey.GetKey()))
	//  pk: 02e9f0f33d79d0459b22ecb253cb3ec52c4ad12f4419fbea03ec825cc74fb3b693
	require.True(t, verified)
	//t.Logf("verified: %t", verified)
}

func TestCosmosSignAndVerifySigLoop(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version:   "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}

	msg := []byte("sign me")

	// both SecureHashByteArray and crypto.Sha256 work
	//hash, err := kmssdk.SecureHashByteArray(msg)
	//FailOnErr(t, err, "FonSecureHash")
	hash := crypto.Sha256(msg)
	//t.Logf("hash: %v", hash)
	var successes int
	var runs = 100
	for i := 0; i < runs; i++ {
		sig, err := _sdk.GetCosmosChainSignature(keyLabel, hash)
		FailOnErr(t, err, "FonSign")

		pubKey, err := _sdk.GetCosmosSECP256K1PubKey(keyLabel)
		FailOnErr(t, err, "FonGetCosmosSECP256K1PubKey")

		verified := pubKey.VerifySignature(msg, sig)

		if verified {
			successes++
		}
	}
	t.Logf("%d of %d", successes, runs)
	require.Equal(t, runs, successes)
}

func TestGetChainSignature(t *testing.T) {
	keyLabel := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Algorithm: kmssdk.Secp256k1,
	}
	message := "sign me"
	hash, err := kmssdk.SecureHash(message)
	FailOnErr(t, err, "FonSecureHash")

	signature, err := _sdk.GetChainSignature(keyLabel, hash)
	FailOnErr(t, err, "FonGetChainSignature")

	t.Logf("got chain signature for key %+v and hash %+v: %+v", keyLabel, hash, signature)
}
