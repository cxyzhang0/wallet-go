package sdk

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/btcsuite/btcd/btcec"
	cosmos256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	tendermintbtct "github.com/tendermint/btcd/btcec"

	//"github.com/decred/dcrd/dcrec/secp256k1/v4"
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
			Kty:   to.Ptr(azkeys.JSONWebKeyTypeEC), // non HSM
			//Kty: to.Ptr(azkeys.JSONWebKeyTypeECHSM), // HSM
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

//	func (s *SDK) GetBTCECPublicKey(keyLabel KeyLabel) (*btcec.PublicKey, error) {
//		ecdsaPublicKey, err := s.GetECDSAPublicKey(keyLabel)
//		if err != nil {
//			return nil, err
//		}
//
//		return btcec.PublicKey(ecdsaPublicKey), nil
//
// }
// secp256k1.PublicKey is used by Cosmos
func (s *SDK) GetSECP256K1PublicKey(keyLabel KeyLabel) (*secp256k1.PublicKey, error) {
	if keyLabel.Algorithm != Secp256k1 {
		return nil, fmt.Errorf("only support secp256k1")
	}

	resp, err := s.GetKey(keyLabel)
	if err != nil {
		return nil, err
	}

	//var x_32, y_32 [32]byte
	//copy(x_32[:], resp.Key.X)
	//copy(y_32[:], resp.Key.Y)
	var x, y secp256k1.FieldVal
	//x.SetBytes(&x_32)
	//y.SetBytes(&y_32)
	x.SetBytes((*[32]byte)(resp.Key.X))
	y.SetBytes((*[32]byte)(resp.Key.Y))

	return secp256k1.NewPublicKey(&x, &y), nil
}

// cosmos secp256k1.PubKey is used by Cosmos
func (s *SDK) GetCosmosSECP256K1PubKey(keyLabel KeyLabel) (*cosmos256k1.PubKey, error) {
	//ecdsaPublicKey, err := s.GetECDSAPublicKey(keyLabel)
	publicKey, err := s.GetSECP256K1PublicKey(keyLabel)
	if err != nil {
		return nil, err
	}
	//var publicKey = btcec.PublicKey(*ecdsaPublicKey)

	pk := publicKey.SerializeCompressed()
	return &cosmos256k1.PubKey{Key: pk}, nil
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

// GetCosmosChainSignature returns the serialized signature ready for cosmos. EC only
// input is normally the hash to be signed.
func (s *SDK) GetCosmosChainSignature(keyLabel KeyLabel, input []byte) ([]byte, error) {
	if keyLabel.Algorithm == RSA2048 {
		return nil, fmt.Errorf("only EC is supported")
	}

	resp, err := s.Sign(keyLabel, input)
	if err != nil {
		return nil, err
	}

	len := len(resp.KeyOperationResult.Result)
	signature := tendermintbtct.Signature{
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
