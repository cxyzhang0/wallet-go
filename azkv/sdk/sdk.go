package sdk

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/btcsuite/btcd/btcec"
	"math/big"
	"runtime"
)

type SDK struct {
	kvc *azkeys.Client
}

func NewSDK(client *azkeys.Client) *SDK {
	s := &SDK{
		kvc: client,
	}
	runtime.SetFinalizer(s, func(s *SDK) {
		s.kvc = nil
	})
	return s
}

// GenerateKeyPair
// To reduce cost, use non HSM
func (s *SDK) GenerateKeyPair(keyLabel KeyLabel) (*azkeys.CreateKeyResponse, error) {
	var params azkeys.CreateKeyParameters
	switch keyLabel.Algorithm {
	case Secp256k1:
		params = azkeys.CreateKeyParameters{
			Curve: to.Ptr(azkeys.JSONWebKeyCurveNameP256K),
			Kty:   to.Ptr(azkeys.JSONWebKeyTypeEC),
			//Kty:   to.Ptr(azkeys.JSONWebKeyTypeECHSM),
		}
	case Secp256p:
		params = azkeys.CreateKeyParameters{
			Curve: to.Ptr(azkeys.JSONWebKeyCurveNameP256),
			Kty:   to.Ptr(azkeys.JSONWebKeyTypeEC),
			//Kty:   to.Ptr(azkeys.JSONWebKeyTypeECHSM),
		}
	case RSA2048:
		params = azkeys.CreateKeyParameters{
			KeySize: to.Ptr(int32(2048)),
			Kty:     to.Ptr(azkeys.JSONWebKeyTypeRSA),
			//Kty:     to.Ptr(azkeys.JSONWebKeyTypeRSAHSM),
		}
	}

	resp, err := s.kvc.CreateKey(context.TODO(), keyLabel.Key, params, nil)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) GetKey(keyLabel KeyLabel) (*azkeys.GetKeyResponse, error) {
	resp, err := s.kvc.GetKey(context.TODO(), keyLabel.Key, keyLabel.Version, nil)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) GetECDSAPublicKey(keyLabel KeyLabel) (*ecdsa.PublicKey, error) {
	var curve elliptic.Curve
	switch keyLabel.Algorithm {
	case Secp256k1:
		curve = btcec.S256()
	case Secp256p:
		curve = elliptic.P256()
	case RSA2048:
		return nil, fmt.Errorf("only support EC")
	}

	resp, err := s.GetKey(keyLabel)
	if err != nil {
		return nil, err
	}

	pubKey := ecdsa.PublicKey{
		Curve: curve,
		X:     &big.Int{},
		Y:     &big.Int{},
	}
	pubKey.X = new(big.Int).SetBytes(resp.Key.X)
	pubKey.Y = new(big.Int).SetBytes(resp.Key.Y)

	return &pubKey, nil
}

func (s *SDK) Sign(keyLabel KeyLabel, input []byte) (*azkeys.SignResponse, error) {
	var algo *azkeys.JSONWebKeySignatureAlgorithm
	switch keyLabel.Algorithm {
	case Secp256k1:
		algo = to.Ptr(azkeys.JSONWebKeySignatureAlgorithmES256K)
	case Secp256p:
		algo = to.Ptr(azkeys.JSONWebKeySignatureAlgorithmES256)
	case RSA2048:
		algo = to.Ptr(azkeys.JSONWebKeySignatureAlgorithmRS256)
	}

	params := azkeys.SignParameters{
		Algorithm: algo,
		Value:     input,
	}
	resp, err := s.kvc.Sign(context.TODO(), keyLabel.Key, keyLabel.Version, params, nil)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetChainSignature returns the serialized signature ready for btc and eth. EC only
// input is normally the hash to be signed.
func (s *SDK) GetChainSignature(keyLabel KeyLabel, input []byte) ([]byte, error) {
	if keyLabel.Algorithm == RSA2048 {
		return nil, fmt.Errorf("only EC is supported")
	}

	resp, err := s.Sign(keyLabel, input)
	if err != nil {
		return nil, err
	}

	len := len(resp.KeyOperationResult.Result)
	signature := btcec.Signature{
		R: &big.Int{},
		S: &big.Int{},
	}
	rBytes := resp.KeyOperationResult.Result[:len/2]
	sBytes := resp.KeyOperationResult.Result[len/2:]
	signature.R = new(big.Int).SetBytes(rBytes)
	signature.S = new(big.Int).SetBytes(sBytes)

	return signature.Serialize(), nil
}

func (s *SDK) SignToECDSA(keyLabel KeyLabel, input []byte) (*azkeys.SignResponse, error) {
	var algo *azkeys.JSONWebKeySignatureAlgorithm
	switch keyLabel.Algorithm {
	case Secp256k1:
		algo = to.Ptr(azkeys.JSONWebKeySignatureAlgorithmES256K)
	case Secp256p:
		algo = to.Ptr(azkeys.JSONWebKeySignatureAlgorithmES256)
	case RSA2048:
		algo = to.Ptr(azkeys.JSONWebKeySignatureAlgorithmRS256)
	}

	params := azkeys.SignParameters{
		Algorithm: algo,
		Value:     input,
	}
	resp, err := s.kvc.Sign(context.TODO(), keyLabel.Key, keyLabel.Version, params, nil)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *SDK) VerifySig(keyLabel KeyLabel, input []byte, sig []byte) (*azkeys.VerifyResponse, error) {
	var algo *azkeys.JSONWebKeySignatureAlgorithm
	switch keyLabel.Algorithm {
	case Secp256k1:
		algo = to.Ptr(azkeys.JSONWebKeySignatureAlgorithmES256K)
	case Secp256p:
		algo = to.Ptr(azkeys.JSONWebKeySignatureAlgorithmES256)
	case RSA2048:
		algo = to.Ptr(azkeys.JSONWebKeySignatureAlgorithmRS256)
	}
	params := azkeys.VerifyParameters{
		Algorithm: algo,
		Digest:    input,
		Signature: sig,
	}
	resp, err := s.kvc.Verify(context.TODO(), keyLabel.Key, keyLabel.Version, params, nil)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
