package sdk

import (
	kms "cloud.google.com/go/kms/apiv1"
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"hash/crc32"
	"strconv"
	"time"

	//kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"runtime"
)

type SDK struct {
	ctx *context.Context
	kmc *kms.KeyManagementClient
}

func NewSDK() (*SDK, error) {
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	s := &SDK{
		&ctx,
		client,
	}

	runtime.SetFinalizer(s, func(s *SDK) {
		s.kmc.Close()
	})

	return s, nil
}

// CreateKeyRing
/**
keyLabel only needs
	Project
	Location
	KeyRing
*/
func (s *SDK) CreateKeyRing(keyLabel KeyLabel) (*kmspb.KeyRing, error) {
	parent := keyLabel.ParentLocation()
	req := &kmspb.CreateKeyRingRequest{
		Parent:    parent,
		KeyRingId: keyLabel.KeyRing,
	}
	if result, err := s.kmc.CreateKeyRing(*s.ctx, req); err != nil {
		return nil, fmt.Errorf("failed to create key ring for %s %s: %+v", parent, keyLabel.KeyRing, err)
	} else {
		return result, nil
	}
}

func (s *SDK) GenerateKeyPair(keyLabel KeyLabel) (*kmspb.CryptoKey, error) {
	parent := keyLabel.ParentKeyRing()

	var algo kmspb.CryptoKeyVersion_CryptoKeyVersionAlgorithm
	var template *kmspb.CryptoKeyVersionTemplate
	switch keyLabel.Algorithm {
	case RSA2048:
		algo = kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256
		template = &kmspb.CryptoKeyVersionTemplate{
			Algorithm:       algo,
			ProtectionLevel: kmspb.ProtectionLevel_SOFTWARE,
		}
	case Secp256p:
		algo = kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256
		template = &kmspb.CryptoKeyVersionTemplate{
			Algorithm:       algo,
			ProtectionLevel: kmspb.ProtectionLevel_SOFTWARE,
		}
	case Secp256k1:
		algo = kmspb.CryptoKeyVersion_EC_SIGN_SECP256K1_SHA256
		template = &kmspb.CryptoKeyVersionTemplate{
			Algorithm:       algo,
			ProtectionLevel: kmspb.ProtectionLevel_HSM,
		}
	}

	req := &kmspb.CreateCryptoKeyRequest{
		Parent:      parent,
		CryptoKeyId: keyLabel.Key,
		CryptoKey: &kmspb.CryptoKey{
			Purpose:         kmspb.CryptoKey_ASYMMETRIC_SIGN,
			VersionTemplate: template,
			// Optional: customize how long key versions should be kept before destroying.
			DestroyScheduledDuration: durationpb.New(24 * time.Hour),
		},
	}

	result, err := s.kmc.CreateCryptoKey(*s.ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateKeyVersion
/**
Create a new version of an existing key
*/
func (s *SDK) CreateKeyVersion(keyLabel KeyLabel) (*kmspb.CryptoKeyVersion, error) {
	parent := keyLabel.ParentKey()

	req := &kmspb.CreateCryptoKeyVersionRequest{
		Parent: parent,
	}

	result, err := s.kmc.CreateCryptoKeyVersion(*s.ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateKeySetPrimary
/**
Update the primary key version for an existing key
NOTE: Only applicable for non-signing keys
*/
func (s *SDK) UpdateKeySetPrimary(keyLabel KeyLabel) (*kmspb.CryptoKey, error) {
	key := keyLabel.ParentKey()
	version := keyLabel.Version

	req := &kmspb.UpdateCryptoKeyPrimaryVersionRequest{
		Name:               key,
		CryptoKeyVersionId: strconv.Itoa(int(version)),
	}
	result, err := s.kmc.UpdateCryptoKeyPrimaryVersion(*s.ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SDK) GetPublicKeyPem(keyLabel KeyLabel) (string, error) { // pem string, error
	name := keyLabel.String()

	req := &kmspb.GetPublicKeyRequest{
		Name: name,
	}

	result, err := s.kmc.GetPublicKey(*s.ctx, req)
	if err != nil {
		return "", err
	}

	pem := result.Pem
	key := []byte(pem)
	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	if int64(crc32c(key)) != result.PemCrc32C.Value {
		return "", fmt.Errorf("getPublicKey: response corrupted in-transit")
	}

	return pem, nil
}

func (s *SDK) GetECDSAPublicKey(keyLabel KeyLabel) (*ecdsa.PublicKey, error) {
	if keyLabel.Algorithm == Secp256k1 {
		return nil, fmt.Errorf("function GetECDSAPublicKey does not support %s. use GetECDSAPublicKeyForSecp256k1 instead", keyLabel.Algorithm.String())
	}
	if keyLabel.Algorithm == RSA2048 {
		return nil, fmt.Errorf("function GetECDSAPublicKey does not support %s. use GetRSAPublicKey instead", keyLabel.Algorithm.String())
	}
	pemString, err := s.GetPublicKeyPem(keyLabel)
	if err != nil {
		return nil, err
	}

	key := []byte(pemString)
	block, _ := pem.Decode(key)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not elliptic curve")
	}
	return ecKey, nil

	return ecKey, nil
}

func (s *SDK) GetECDSAPublicKeyForSecp256k1(keyLabel KeyLabel) (*ecdsa.PublicKey, error) {
	if keyLabel.Algorithm != Secp256k1 {
		return nil, fmt.Errorf("function GetECDSAPublicKeyForSecp256k1 does not support %s. use GetECDSAPublicKey or GetRSAPublicKey instead", keyLabel.Algorithm.String())
	}

	pemString, err := s.GetPublicKeyPem(keyLabel)
	if err != nil {
		return nil, err
	}

	return PemToPubkey(pemString)
}

func (s *SDK) GetRSAPublicKey(keyLabel KeyLabel) (*rsa.PublicKey, error) {
	if keyLabel.Algorithm != RSA2048 {
		return nil, fmt.Errorf("function GetRSAPublicKey does not support %s. use GetECDSAPublicKey or GetECDSAPublicKeyForSecp256k1 instead", keyLabel.Algorithm.String())
	}
	pemString, err := s.GetPublicKeyPem(keyLabel)
	if err != nil {
		return nil, err
	}

	key := []byte(pemString)
	block, _ := pem.Decode(key)
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func (s *SDK) SignString(keyLabel KeyLabel, message string) ([]byte, []byte, error) { // sig, hash, error
	hash, err := SecureHash(message)
	if err != nil {
		return nil, nil, err
	}

	sig, err := s.Sign(keyLabel, hash)
	if err != nil {
		return nil, hash, err
	}

	return sig, hash, nil
}

func (s *SDK) Sign(keyLabel KeyLabel, input []byte) ([]byte, error) {
	// Optional but recommended: Compute digest's CRC32C.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)

	}
	digestCRC32C := crc32c(input)

	// Build the signing request.
	req := &kmspb.AsymmetricSignRequest{
		Name: keyLabel.String(),
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: input,
			},
		},
		DigestCrc32C: wrapperspb.Int64(int64(digestCRC32C)),
	}

	// Call the API.
	result, err := s.kmc.AsymmetricSign(*s.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to sign digest: %v", err)
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if result.VerifiedDigestCrc32C == false {
		return nil, fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	}
	// TODO Comment out when Name is populated by the server
	if result.Name != req.Name {
		return nil, fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	}
	if int64(crc32c(result.Signature)) != result.SignatureCrc32C.Value {
		return nil, fmt.Errorf("AsymmetricSign: response corrupted in-transit")
	}

	return result.Signature, nil
}
