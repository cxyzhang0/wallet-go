package sdk

import "crypto/sha256"

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
	Key       string
	Version   string
	Algorithm CryptographAlgorithm
}

func SecureHash(message string) ([]byte, error) {
	plainText := []byte(message)
	digest := sha256.New()
	if _, err := digest.Write(plainText); err != nil {
		return nil, err
	}

	return digest.Sum(nil), nil
}
