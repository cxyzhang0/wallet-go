package sdk

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/cryptobyte"
	"math/big"
	"regexp"
	"strconv"
)

type CryptographAlgorithm uint

/**
sha256 is our constant digest function. So we don't name it in the enum.
NOTE: GCP KMS does not support Secp256r1
*/
const ( // enum
	Secp256k1 CryptographAlgorithm = iota
	Secp256p
	RSA2048
)

func (c CryptographAlgorithm) String() string {
	switch c {
	case Secp256k1:
		return "Secp256k1"
	case Secp256p:
		return "Secp256p"
	case RSA2048:
		return "RSA2048"
	}
	return "Undefined"
}

type KeyLabel struct {
	Project   string
	Location  string
	KeyRing   string
	Key       string
	Version   uint
	Algorithm CryptographAlgorithm
}

func (l *KeyLabel) String() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%d", l.Project, l.Location, l.KeyRing, l.Key, l.Version)
}

func (l *KeyLabel) ParentLocation() string {
	return fmt.Sprintf("projects/%s/locations/%s", l.Project, l.Location)
}

func (l *KeyLabel) ParentKeyRing() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", l.Project, l.Location, l.KeyRing)
}

func (l *KeyLabel) ParentKey() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", l.Project, l.Location, l.KeyRing, l.Key)
}

func (l *KeyLabel) Next() KeyLabel {
	return KeyLabel{
		l.Project,
		l.Location,
		l.KeyRing,
		l.Key,
		l.Version + 1,
		l.Algorithm,
	}
}

func StringToKeyLabel(labelStr string) (*KeyLabel, error) {
	pattern := "projects/.{1,}/locations/.{1,}/keyRings/.{1,}/cryptoKeys/.{1,}/cryptoKeyVersions/\\d{1,}"
	match, err := regexp.MatchString(pattern, labelStr)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, fmt.Errorf("%s does not match %s", labelStr, pattern)
	}

	slash := regexp.MustCompile(`/`)
	parts := slash.Split(labelStr, 10)
	n := len(parts)
	version, _ := strconv.Atoi(parts[n-1])
	return &KeyLabel{
		Project:  parts[1],
		Location: parts[3],
		KeyRing:  parts[5],
		Key:      parts[7],
		Version:  uint(version),
	}, nil
}

func SecureHash(message string) ([]byte, error) {
	plainText := []byte(message)
	digest := sha256.New()
	if _, err := digest.Write(plainText); err != nil {
		return nil, err
	}

	return digest.Sum(nil), nil
}

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
