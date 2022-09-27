package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	tsx "github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"golang.org/x/crypto/cryptobyte"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"hash/crc32"
	"io"
	"math/big"

	"context"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

var (
	OidPublicKeyECDSA = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	OidSecp256k1      = asn1.ObjectIdentifier{1, 3, 132, 0, 10}
)

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

// pkcs1PublicKey reflects the ASN.1 structure of a PKCS #1 public key.
type pkcs1PublicKey struct {
	N *big.Int
	E int
}

func SignTx(msgTx *wire.MsgTx, wif *btcutil.WIF, fromAddrScript *[]byte) error {
	for i, txIn := range msgTx.TxIn {
		sigScript, err := _SignatureScript(
			msgTx,
			i,
			*fromAddrScript,
			tsx.SigHashAll,
			wif.PrivKey,
			false)
		if err != nil {
			return err
		}
		txIn.SignatureScript = sigScript
	}
	return nil
}
func _SignatureScript(tx *wire.MsgTx, idx int, subscript []byte, hashType tsx.SigHashType, privKey *btcec.PrivateKey, compress bool) ([]byte, error) {
	sig, err := _RawTxInSignature(tx, idx, subscript, hashType, privKey)
	if err != nil {
		return nil, err
	}

	pk := (*btcec.PublicKey)(&privKey.PublicKey)
	var pkData []byte
	if compress {
		pkData = pk.SerializeCompressed()
	} else {
		pkData = pk.SerializeUncompressed()
	}

	return tsx.NewScriptBuilder().AddData(sig).AddData(pkData).Script()
}

func _RawTxInSignature(tx *wire.MsgTx, idx int, subScript []byte,
	hashType tsx.SigHashType, key *btcec.PrivateKey) ([]byte, error) {

	hash, err := tsx.CalcSignatureHash(subScript, hashType, tx, idx)
	if err != nil {
		return nil, err
	}
	signature, err := key.Sign(hash)
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx input: %s", err)
	}

	return append(signature.Serialize(), byte(hashType)), nil
}

// signAsymmetric will sign a plaintext message using a saved asymmetric private
// key stored in Cloud KMS.
func signAsymmetric(w io.Writer, name string, message string) error {
	//name := "projects/quantum-pilot-360000/locations/us-west1/keyRings/WIM-test/cryptoKeys/secp256k1-hsm-1/cryptoKeyVersions/1"
	// message := "my message"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}
	defer client.Close()

	// Convert the message into bytes. Cryptographic plaintexts and
	// ciphertexts are always byte arrays.
	plaintext := []byte(message)

	// Calculate the digest of the message.
	digest := sha256.New()
	if _, err := digest.Write(plaintext); err != nil {
		return fmt.Errorf("failed to create digest: %v", err)
	}

	// Optional but recommended: Compute digest's CRC32C.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)

	}
	digestCRC32C := crc32c(digest.Sum(nil))

	// Build the signing request.
	//
	// Note: Key algorithms will require a varying hash function. For example,
	// EC_SIGN_P384_SHA384 requires SHA-384.
	req := &kmspb.AsymmetricSignRequest{
		Name: name,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest.Sum(nil),
			},
		},
		DigestCrc32C: wrapperspb.Int64(int64(digestCRC32C)),
	}

	// Call the API.
	result, err := client.AsymmetricSign(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to sign digest: %v", err)
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if result.VerifiedDigestCrc32C == false {
		return fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	}
	// TODO(iamtamjam) Uncomment when this field is populated by the server
	// if result.Name != req.Name {
	//      return fmt.Errorf("AsymmetricSign: request corrupted in-transit")
	// }
	if int64(crc32c(result.Signature)) != result.SignatureCrc32C.Value {
		return fmt.Errorf("AsymmetricSign: response corrupted in-transit")
	}

	fmt.Fprintf(w, "Signed digest: %s", result.Signature)
	return nil
}

func getPublicKey(name string) (*ecdsa.PublicKey, error) {
	//name := "projects/quantum-pilot-360000/locations/us-west1/keyRings/WIM-test/cryptoKeys/secp256k1-hsm-1/cryptoKeyVersions/1"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.GetPublicKeyRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetPublicKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %v", err)
	}

	// Extract and print the labels.
	fmt.Printf("name=%s; pem=%s; str=%s\n", result.Name, result.Pem, result.String())

	return PemToPubkey(result.Pem)
}

// this does not work with secp256k1
func getPublicKey_old(name string) (*ecdsa.PublicKey, error) {
	//name := "projects/quantum-pilot-360000/locations/us-west1/keyRings/WIM-test/cryptoKeys/secp256k1-hsm-1/cryptoKeyVersions/1"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create kms client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.GetPublicKeyRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetPublicKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %v", err)
	}

	// Extract and print the labels.
	fmt.Printf("name=%s; pem=%s; str=%s\n", result.Name, result.Pem, result.String())
	// Parse the public key. Note, this example assumes the public key is in the
	// ECDSA format.
	block, _ := pem.Decode([]byte(result.Pem))
	//publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)  // RSA
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)

	//publicKey, err := secp.ParsePubKey(block.Bytes)

	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}
	ecKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not elliptic curve")
	}
	return ecKey, nil
}

// other useful funcs:  https://github.com/vanhallio/go-secp256k1-pem/blob/main/s256Pem.go
func PemToPubkey(pemString string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	derBytes := block.Bytes
	var pki publicKeyInfo
	if rest, err := asn1.Unmarshal(derBytes, &pki); err != nil {
		if _, err := asn1.Unmarshal(derBytes, &pkcs1PublicKey{}); err == nil {
			return nil, errors.New("x509: failed to parse public key (use ParsePKCS1PublicKey instead for this key format)")
		}
		return nil, err
	} else if len(rest) != 0 {
		return nil, errors.New("x509: trailing data after ASN.1 of public-key")
	}

	if !pki.Algorithm.Algorithm.Equal(OidPublicKeyECDSA) {
		return nil, errors.New("x509: not a ECDSA public key")
	}

	der := cryptobyte.String(pki.PublicKey.RightAlign())
	paramsDer := cryptobyte.String(pki.Algorithm.Parameters.FullBytes)
	namedCurveOID := new(asn1.ObjectIdentifier)
	if !paramsDer.ReadASN1ObjectIdentifier(namedCurveOID) {
		return nil, errors.New("x509: invalid ECDSA parameters")
	}

	if !namedCurveOID.Equal(OidSecp256k1) {
		return nil, errors.New("x509: not a secp256k1 curve")
	}
	curve := btcec.S256()
	//curve := crypto.S256()

	x, y := elliptic.Unmarshal(curve, der)
	if x == nil {
		return nil, errors.New("x509: failed to unmarshal secp256k1 curve point")
	}
	pub := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	return pub, nil
}

// getKeyLabels fetches the labels on a KMS key.
func getKeyLabelsW(w io.Writer, name string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.GetCryptoKeyRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetCryptoKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get key: %v", err)
	}

	// Extract and print the labels.
	for k, v := range result.Labels {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return nil
}

func getKeyLabels(name string) error {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.GetCryptoKeyRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.GetCryptoKey(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get key: %v", err)
	}

	// Extract and print the labels.
	for k, v := range result.Labels {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
