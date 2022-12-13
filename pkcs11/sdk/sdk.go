package sdk

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"encoding/asn1"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	p11 "github.com/miekg/pkcs11"
	"github.com/pkg/errors"
	"math/big"
	"runtime"
)

var SupportedAlgorithms map[CryptographAlgorithm][][]*p11.Attribute = map[CryptographAlgorithm][][]*p11.Attribute{ // curve to [publicKeyTemplate, privateKeyTemplate]
	Secp256k1: {
		{
			p11.NewAttribute(p11.CKA_EC_PARAMS, []byte{
				0x06, 0x05, 0x2B, 0x81, 0x04, 0x00, 0x0A,
			}), // OID 1.3.132.0.10 Secp256k1
			p11.NewAttribute(p11.CKA_VERIFY, true),
		},
		{
			p11.NewAttribute(p11.CKA_SIGN, true),
			p11.NewAttribute(p11.CKA_SENSITIVE, true),
			p11.NewAttribute(p11.CKA_EXTRACTABLE, true),
		},
	},
	Ed25519: {
		{
			p11.NewAttribute(p11.CKA_EC_PARAMS, []byte{
				0x06, 0x03, 0x2B, 0x65, 0x70, // from oid 1.3.101.112 for id-Ed25519
				//0x13, 0x0c, 0x65, 0x64, 0x77, 0x61, 0x72, 0x64, 0x73, 0x32, 0x35, 0x35, 0x31, 0x39, // is it for Edwards25519? https://stackoverflow.com/questions/70740280/how-to-export-ed25519-public-key-from-pkcs11-attribute
				//0x06, 0x09, 0x2B, 0x06, 0x01, 0x04, 0x01, 0xDA, 0x47, 0x0F, 0x01, // did not work: from oid 1.3.6.1.4.1.11591.15.1 for ed25519
			}),
			p11.NewAttribute(p11.CKA_VERIFY, true),
		},
		{
			p11.NewAttribute(p11.CKA_SIGN, true),
			p11.NewAttribute(p11.CKA_SENSITIVE, true),
			p11.NewAttribute(p11.CKA_EXTRACTABLE, true),
		},
	},
	Secp256r1: {
		{
			p11.NewAttribute(p11.CKA_EC_PARAMS, []byte{
				0x06, 0x08, 0x2A, 0x86, 0x48, 0xCE, 0x3D, 0x03, 0x01, 0x07,
			}), // OID 1.2.840.10045.3.1.7 Secp256r1/v1, etc
			p11.NewAttribute(p11.CKA_VERIFY, true),
		},
		{
			p11.NewAttribute(p11.CKA_SIGN, true),
			p11.NewAttribute(p11.CKA_SENSITIVE, true),
			p11.NewAttribute(p11.CKA_EXTRACTABLE, true),
		},
	},
	RSA2048: {
		{
			p11.NewAttribute(p11.CKA_CLASS, p11.CKO_PUBLIC_KEY),
			p11.NewAttribute(p11.CKA_KEY_TYPE, p11.CKK_RSA),
			p11.NewAttribute(p11.CKA_PUBLIC_EXPONENT, []byte{1, 0, 1}),
			p11.NewAttribute(p11.CKA_MODULUS_BITS, 2048),
		},
		{
			p11.NewAttribute(p11.CKA_SIGN, true),
			p11.NewAttribute(p11.CKA_SENSITIVE, true),
			p11.NewAttribute(p11.CKA_EXTRACTABLE, true),
		}},
}

type SDK struct {
	P       *p11.Ctx
	Session p11.SessionHandle
	//supportedCurves map[CryptographAlgorithm][][]*p11.Attribute
}

func NewSDK(lib, slotLabel, pin string) (*SDK, error) {
	p := p11.New(lib)
	if p == nil {
		return nil, fmt.Errorf("failed to init lib")
	}

	if e := p.Initialize(); e != nil {
		return nil, e
	}

	slotID, e := GetSlot(p, slotLabel)
	if e != nil {
		return nil, e
	}

	session, e := p.OpenSession(slotID, p11.CKF_SERIAL_SESSION|p11.CKF_RW_SESSION)
	if e != nil {
		return nil, e
	}

	if e = p.Login(session, p11.CKU_USER, pin); e != nil {
		return nil, e
	}

	s := &SDK{
		P:       p,
		Session: session,
		//supportedCurves: curvesMap,
	}

	runtime.SetFinalizer(s, func(s *SDK) {
		s.P.Logout(s.Session)
		s.P.CloseSession(s.Session)
		s.P.Finalize()
		s.P.Destroy()
	})

	return s, nil
}

func (s *SDK) GenerateKeyPair(keyLabel KeyLabel, tokenPersistent bool) (p11.ObjectHandle, p11.ObjectHandle, error) {
	t, ok := SupportedAlgorithms[keyLabel.Algorithm]
	if !ok {
		return 0, 0, fmt.Errorf("curve %s is not supported", keyLabel.Algorithm.String())
	}

	publicKeyTemplate := t[0]
	publicKeyTemplate = append(publicKeyTemplate,
		[]*p11.Attribute{
			p11.NewAttribute(p11.CKA_VERIFY, true),
			p11.NewAttribute(p11.CKA_TOKEN, tokenPersistent),
			p11.NewAttribute(p11.CKA_LABEL, keyLabel.ShortLabel()),
		}...,
	)

	privateKeyTemplate := t[1]
	privateKeyTemplate = append(privateKeyTemplate,
		[]*p11.Attribute{
			p11.NewAttribute(p11.CKA_SIGN, true),
			p11.NewAttribute(p11.CKA_SENSITIVE, true),
			p11.NewAttribute(p11.CKA_EXTRACTABLE, true),
			p11.NewAttribute(p11.CKA_TOKEN, tokenPersistent),
			p11.NewAttribute(p11.CKA_LABEL, keyLabel.ShortLabel()),
		}...)

	var keyGenMechanism []*p11.Mechanism
	switch keyLabel.Algorithm {
	case RSA2048:
		keyGenMechanism = []*p11.Mechanism{p11.NewMechanism(p11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)}
	case Ed25519:
		keyGenMechanism = []*p11.Mechanism{p11.NewMechanism(p11.CKM_EC_EDWARDS_KEY_PAIR_GEN, nil)}
	default:
		keyGenMechanism = []*p11.Mechanism{p11.NewMechanism(p11.CKM_EC_KEY_PAIR_GEN, nil)}
	}

	pbk, pvk, err := s.P.GenerateKeyPair(s.Session, keyGenMechanism, publicKeyTemplate, privateKeyTemplate)
	if err != nil {
		return 0, 0, err
	}

	return pbk, pvk, nil
}

func (s *SDK) Sign(keyLabel KeyLabel, input []byte) ([]byte, error) {
	prvk, err := s.GetPrivateKeyHandle(keyLabel)
	if err != nil {
		return nil, err
	}

	var mechanism []*p11.Mechanism
	switch keyLabel.Algorithm {
	case RSA2048:
		mechanism = []*p11.Mechanism{
			p11.NewMechanism(p11.CKM_RSA_PKCS, nil),
		}
	case Ed25519:
		mechanism = []*p11.Mechanism{
			p11.NewMechanism(p11.CKM_EDDSA, nil),
		}
	default:
		mechanism = []*p11.Mechanism{
			//p11.NewMechanism(p11.CKM_ECDSA_SHA256, nil), // this gives CKR_MECHANISM_INVALID
			//p11.NewMechanism(p11.CKM_ECDSA_SHA1, nil), // this gives CKR_MECHANISM_INVALID
			p11.NewMechanism(p11.CKM_ECDSA, nil),
		}
	}

	if err := s.P.SignInit(s.Session, mechanism, prvk); err != nil {
		return nil, err
	}

	sig, err := s.P.Sign(s.Session, input)
	if err != nil {
		return nil, err
	}

	if err := s.VerifySig(keyLabel, input, sig); err != nil {
		return nil, err
	}

	return sig, nil
}

func (s *SDK) VerifySig(keyLabel KeyLabel, input []byte, sig []byte) error {
	pubk, err := s.GetPublicKeyHandle(keyLabel)
	if err != nil {
		return err
	}

	var mechanism []*p11.Mechanism
	switch keyLabel.Algorithm {
	case RSA2048:
		mechanism = []*p11.Mechanism{
			p11.NewMechanism(p11.CKM_RSA_PKCS, nil),
		}
	case Ed25519:
		mechanism = []*p11.Mechanism{
			p11.NewMechanism(p11.CKM_EDDSA, nil),
		}
	default:
		mechanism = []*p11.Mechanism{
			p11.NewMechanism(p11.CKM_ECDSA, nil),
		}
	}

	// verify
	if err := s.P.VerifyInit(s.Session, mechanism, pubk); err != nil {
		return err
	}

	if err := s.P.Verify(s.Session, input, sig); err != nil {
		return err
	}

	// It is not multi-part signature
	//if err := s.P.VerifyFinal(s.Session, sig); err != nil {
	//	return nil, err
	//}

	return nil
}

// GetChainSignature returns the serialized signature ready for btc and eth. EC only
func (s *SDK) GetChainSignature(keyLabel KeyLabel, input []byte) ([]byte, error) {
	if keyLabel.Algorithm == RSA2048 {
		return nil, fmt.Errorf("only EC is supported")
	}

	sig, err := s.Sign(keyLabel, input)
	if err != nil {
		return nil, err
	}
	len := len(sig)
	signature := btcec.Signature{
		R: &big.Int{},
		S: &big.Int{},
	}

	rBytes := sig[:len/2]
	sBytes := sig[len/2:]
	signature.R = new(big.Int).SetBytes(rBytes)
	signature.S = new(big.Int).SetBytes(sBytes)

	return signature.Serialize(), nil
}

func (s *SDK) findObjects(template []*p11.Attribute, args ...int) ([]p11.ObjectHandle, error) {
	max := 100 // default
	if len(args) > 0 {
		max = args[0]
	}

	if err := s.P.FindObjectsInit(s.Session, template); err != nil {
		return nil, err
	}

	objs, _, err := s.P.FindObjects(s.Session, max)
	if err != nil {
		return nil, err
	}

	if err = s.P.FindObjectsFinal(s.Session); err != nil {
		return nil, err
	}

	return objs, nil
}

// GetPrivateKeyHandle pkcs11 alllows multiple private keys with the same label
// Application needs to enforce uniqueness.
// This function will return error if uniqueness is violated.
func (s *SDK) GetPrivateKeyHandle(keyLabel KeyLabel) (p11.ObjectHandle, error) {
	label := keyLabel.ShortLabel()
	var noKey p11.ObjectHandle
	template := []*p11.Attribute{
		p11.NewAttribute(p11.CKA_CLASS, p11.CKO_PRIVATE_KEY),
		p11.NewAttribute(p11.CKA_LABEL, label),
	}

	objs, err := s.findObjects(template)
	if err != nil {
		return noKey, err
	}
	//if err := s.P.FindObjectsInit(s.Session, template); err != nil {
	//	return noKey, err
	//}
	//
	//objs, _, err := s.P.FindObjects(s.Session, 2)
	//if err != nil {
	//	return noKey, err
	//}
	//
	//if err = s.P.FindObjectsFinal(s.Session); err != nil {
	//	return noKey, err
	//}

	if len(objs) == 0 {
		err = fmt.Errorf("private key not found")
		return noKey, err
	}

	if len(objs) > 1 {
		err = fmt.Errorf("more than 1 private key is found")
		return noKey, err
	}

	return objs[0], nil
}

// GetPublicKeyHandle return unique public key with the given keyLabel
func (s *SDK) GetPublicKeyHandle(keyLabel KeyLabel) (p11.ObjectHandle, error) {
	label := keyLabel.ShortLabel()
	var noKey p11.ObjectHandle
	template := []*p11.Attribute{
		p11.NewAttribute(p11.CKA_CLASS, p11.CKO_PUBLIC_KEY),
		p11.NewAttribute(p11.CKA_LABEL, label),
	}

	objs, err := s.findObjects(template)
	if err != nil {
		return noKey, err
	}

	if len(objs) == 0 {
		err = fmt.Errorf("public key not found")
		return noKey, err
	}

	if len(objs) > 1 {
		err = fmt.Errorf("more than 1 public key is found")
		return noKey, err
	}

	return objs[0], nil
}

//GetPublicKeyECPoint
/**
SoftHSM does not calculate CKA_PUBLIC_KEY_INFO
https://github.com/opendnssec/SoftHSMv2/blob/f82d4eda55401a4d23e647d85a00a8b0c8ccf712/src/lib/P11Objects.cpp
Neither does FX
So we use CKA_EC_POINT instead.
*/
func (s *SDK) GetPublicKeyECPoint(keyLabel KeyLabel) (*p11.Attribute, error) {
	if keyLabel.Algorithm == RSA2048 {
		return nil, fmt.Errorf("only EC curve is supported.")
	}

	pubk, err := s.GetPublicKeyHandle(keyLabel)
	if err != nil {
		return nil, err
	}

	template := []*p11.Attribute{
		p11.NewAttribute(p11.CKA_EC_POINT, nil),
		//p11.NewAttribute(p11.CKA_VALUE, nil),
		//p11.NewAttribute(p11.CKA_PUBLIC_KEY_INFO, nil),
	}

	attr, err := s.P.GetAttributeValue(s.Session, p11.ObjectHandle(pubk), template)
	if err != nil {
		return nil, err
	}

	if len(attr) != 1 {
		return nil, fmt.Errorf("got %d attributes. expect 1", len(attr))
	}
	return attr[0], nil
}

func (s *SDK) GetECDSAPublicKey(keyLabel KeyLabel) (*ecdsa.PublicKey, error) {
	var curve elliptic.Curve
	switch keyLabel.Algorithm {
	case Secp256k1:
		curve = btcec.S256()
	case Secp256r1:
		curve = elliptic.P256()
	case Ed25519:
		return nil, fmt.Errorf("ed25519 is not supported")
	case RSA2048:
		return nil, fmt.Errorf("RSA is not supported")
	}

	ecPoint, err := s.GetPublicKeyECPoint(keyLabel)
	if err != nil {
		return nil, err
	}

	var pointBytes []byte
	extra, err := asn1.Unmarshal(ecPoint.Value, &pointBytes) // ecPoint.Value is 67 bytes
	if err != nil {
		return nil, errors.WithMessage(err, "ec point is invalid ASN.1")
	}
	if len(extra) > 0 {
		// We weren't expecting extra data
		return nil, errors.New("unexpected data found when parsing elliptic curve point")
	}

	x, y := elliptic.Unmarshal(curve, pointBytes)

	pubKey := ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	/*
		pubKey := ecdsa.PublicKey{
			Curve: curve,
			X:     &big.Int{},
			Y:     &big.Int{},
		}
		xBytes := ecPoint.Value[3:35] //the point array has 67 bytes, ignore the first 3
		yBytes := ecPoint.Value[35:]

		pubKey.X = new(big.Int).SetBytes(xBytes)
		pubKey.Y = new(big.Int).SetBytes(yBytes)
	*/
	return &pubKey, nil
}

func (s *SDK) GetEdDSAPublicKey(keyLabel KeyLabel) (*ed25519.PublicKey, error) {
	if keyLabel.Algorithm != Ed25519 {
		return nil, fmt.Errorf("on Ed25519 is supported")
	}

	ecPoint, err := s.GetPublicKeyECPoint(keyLabel)
	if err != nil {
		return nil, err
	}

	var pointBytes []byte
	extra, err := asn1.Unmarshal(ecPoint.Value, &pointBytes) // ecPoint.Value is 34 bytes
	if err != nil {
		return nil, errors.WithMessage(err, "ec point is invalid ASN.1")
	}
	if len(extra) > 0 {
		// We weren't expecting extra data
		return nil, errors.New("unexpected data found when parsing elliptic curve point")
	}

	pubKey := ed25519.PublicKey(pointBytes) // pointBytes is 32 bytes

	return &pubKey, nil
}

func (s *SDK) GetPublicKeyAttrFromPrivateKey(keyLabel KeyLabel) (*p11.Attribute, error) {
	privk, err := s.GetPrivateKeyHandle(keyLabel)
	if err != nil {
		return nil, err
	}

	template := []*p11.Attribute{
		p11.NewAttribute(p11.CKA_PUBLIC_KEY_INFO, nil),
	}

	attr, err := s.P.GetAttributeValue(s.Session, p11.ObjectHandle(privk), template)
	if err != nil {
		return nil, err
	}

	if len(attr) != 1 {
		return nil, fmt.Errorf("got %d attributes. expect 1", len(attr))
	}
	return attr[0], nil
}

func (s *SDK) GetAllPrivateKeys() ([]p11.ObjectHandle, error) {
	template := []*p11.Attribute{
		p11.NewAttribute(p11.CKA_CLASS, p11.CKO_PRIVATE_KEY),
	}

	objs, err := s.findObjects(template)
	if err != nil {
		return nil, err
	}

	if len(objs) == 0 {
		err = fmt.Errorf("private keys not found")
		return nil, err
	}
	return objs, nil
}

func (s *SDK) GetAllECKeys() ([]p11.ObjectHandle, error) {
	template := []*p11.Attribute{
		p11.NewAttribute(p11.CKA_KEY_TYPE, p11.CKK_EC),
	}

	objs, err := s.findObjects(template)
	if err != nil {
		return nil, err
	}

	if len(objs) == 0 {
		err = fmt.Errorf("EC keys not found")
		return nil, err
	}
	return objs, nil
}

// GetPubKeyBySig
/**
TODO: ethereum SigToPub does not like the sig (length 65 vs 64)
this discussion may give some clue
https://github.com/celo-org/optics-monorepo/discussions/598
*/
func (s *SDK) GetPubKeyBySig(keyLabel KeyLabel) (*ecdsa.PublicKey, error) {
	if keyLabel.Algorithm != Secp256k1 {
		return nil, fmt.Errorf("only secp256k1 is supported")
	}
	message := "sign me"

	hash, err := SecureHash(message)
	//hash, err := Keccak256Hash(message)
	//hash, err := DigestSHA256(s.P, s.Session, message)
	if err != nil {
		return nil, err
	}

	sig, err := s.Sign(keyLabel, hash)
	if err != nil {
		return nil, err
	}

	return SigToPub(hash, sig)
}

// this func may not be used since NewSDK has SetFinalizer(...)
func (s *SDK) closeSession() {
	s.P.Logout(s.Session)
	s.P.CloseSession(s.Session)
	s.P.Finalize()
	s.P.Destroy()
}
