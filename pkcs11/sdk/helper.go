package sdk

import (
	"fmt"
	p11 "github.com/miekg/pkcs11"
	"regexp"
	"strconv"
	"strings"
)

//TODO: Should we include SlotId/Token?
type KeyLabel struct {
	Prefix  string
	KeyRing string
	Key     string
	Version uint
	Curve   Curve
}

func (l *KeyLabel) Label() string {
	return fmt.Sprintf("%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%d", l.Prefix, l.KeyRing, l.Key, l.Version)
}
func (l *KeyLabel) ShortLabel() string {
	return fmt.Sprintf("%s/%s/%s/%d", l.Prefix, l.KeyRing, l.Key, l.Version)
}
func (l *KeyLabel) String() string {
	return l.Label()
}
func (l *KeyLabel) ShortString() string {
	return l.ShortLabel()
}
func (l *KeyLabel) Next() KeyLabel {
	return KeyLabel{
		l.Prefix,
		l.Key,
		l.Key,
		l.Version + 1,
		l.Curve,
	}
}
func StringToKeyLabel(labelStr string) (*KeyLabel, error) {
	pattern := ".*/keyRings/.{1,}/cryptoKeys/.{1,}/cryptoKeyVersions/\\d{1,}"
	match, err := regexp.MatchString(pattern, labelStr)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, fmt.Errorf("%s does not match %s", labelStr, pattern)
	}

	slash := regexp.MustCompile(`/`)
	parts := slash.Split(labelStr, -1)
	n := len(parts)
	version, _ := strconv.Atoi(parts[n-1])
	pf := strings.Join(parts[:n-6], "/")
	return &KeyLabel{
		Prefix:  pf, //parts[n-7],
		KeyRing: parts[n-5],
		Key:     parts[n-3],
		Version: uint(version),
	}, nil
}
func ShortStringToKeyLabel(labelStr string) (*KeyLabel, error) {
	pattern := ".*/.{1,}/.{1,}/\\d{1,}"
	match, err := regexp.MatchString(pattern, labelStr)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, fmt.Errorf("%s does not match %s", labelStr, pattern)
	}

	slash := regexp.MustCompile(`/`)
	parts := slash.Split(labelStr, -1)
	n := len(parts)
	version, _ := strconv.Atoi(parts[n-1])
	pf := strings.Join(parts[:n-3], "/")
	return &KeyLabel{
		Prefix:  pf, //parts[n-4],
		KeyRing: parts[n-3],
		Key:     parts[n-2],
		Version: uint(version),
	}, nil
}

type Curve uint

const ( // enum
	Secp256k1 Curve = iota
	Secp256r1
	RSA
)

func (c Curve) String() string {
	switch c {
	case Secp256k1:
		return "Secp256k1"
	case Secp256r1:
		return "Secp256r1"
	case RSA:
		return "RSA"
	}
	return "Undefined"
}

func ListAllSlots(p *p11.Ctx) {
	slots, err := p.GetSlotList(true)
	if err != nil {
		fmt.Printf("error get slot list: %+v", err)
		return
	}
	for _, slot := range slots {
		slotInfo, err := p.GetSlotInfo(slot)
		if err != nil {
			fmt.Printf("error get slot info for slot %d: %+v", slot, err)
			return
		}
		tokenInfo, err := p.GetTokenInfo(slot)
		if err != nil {
			fmt.Printf("error get token info for slot %d: %+v", slot, err)
			return
		}
		fmt.Printf("slotInfo for slot %d: %+v \n", slot, slotInfo)
		fmt.Printf("tokenInfo for slot %d: %+v \n", slot, tokenInfo)
	}
}

func GetSlot(p *p11.Ctx, slotLabel string) (uint, error) {
	slots, err := p.GetSlotList(true)
	if err != nil {
		return 0, err
	}
	for _, slot := range slots {
		_, err := p.GetSlotInfo(slot)
		if err != nil {
			return 0, err
		}
		tokenInfo, err := p.GetTokenInfo(slot)
		if err != nil {
			return 0, err
		}
		if tokenInfo.Label == slotLabel {
			return slot, nil
		}
	}
	return 0, fmt.Errorf("slot not found: %s", slotLabel)
}

func DigestSHA256(p *p11.Ctx, session p11.SessionHandle, message string) ([]byte, error) {
	e := p.DigestInit(session, []*p11.Mechanism{p11.NewMechanism(p11.CKM_SHA256, nil)})
	if e != nil {
		return nil, e
	}

	hash, e := p.Digest(session, []byte(message))
	if e != nil {
		return nil, e
	}

	//P.DigestFinal(Session)
	return hash, nil
}
