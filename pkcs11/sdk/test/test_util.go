package test

import (
	"fmt"
	"github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	p11 "github.com/miekg/pkcs11"
	"os"
	"testing"
)

var (
	//// FX HSM
	module     = "/Users/johnz/futurex/wf_pkcs11_mac_447/libfxpkcs11-debug.dylib"
	tokenLabel = "us01crypto01test.virtucrypt.com:"
	//module     = "/Users/johnz/futurex/wf_pkcs11_mac/libfxpkcs11-debug.dylib"
	//tokenLabel = "us01hsm01test.virtucrypt.com:592"
	pin = "safest"
	////privateKeyLabel = "projects/quantum-pilot-360000/locations/us-west1/keyRings/WIM-test/cryptoKeys/secp256k1-hsm-1/cryptoKeyVersions/1"

	//// softhsm
	//module     = "/usr/local/lib/softhsm/libsofthsm2.so"
	//tokenLabel = "Slot Token 0"
	//pin        = "5678"
	////privateKeyLabel = "projects/quantum-pilot-360000/locations/us-west1/keyRings/WIM-test/cryptoKeys/secp256k1-hsm-1/cryptoKeyVersions/1"

	_sdk *sdk.SDK
)

func init() {
	if x := os.Getenv("SOFTHSM_LIB"); x != "" {
		module = x
	}
	if x := os.Getenv("SOFTHSM_TOKENLABEL"); x != "" {
		tokenLabel = x
	}
	//if x := os.Getenv("SOFTHSM_PRIVKEYLABEL"); x != "" {
	//	privateKeyLabel = x
	//}
	if x := os.Getenv("SOFTHSM_PIN"); x != "" {
		pin = x
	}
	wd, _ := os.Getwd()
	//os.Setenv("SOFTHSM_CONF", wd+"/softhsm.conf")
	os.Setenv("SOFTHSM2_CONF", wd+"/softhsm2.conf")

	os.Setenv("FXPKCS11_CFG", "/Users/johnz/futurex/wf_pkcs11_mac_447/fxpkcs11.cfg")
	//os.Setenv("FXPKCS11_CFG", "/Users/johnz/futurex/wf_pkcs11_mac/fxpkcs11.cfg")
	_sdk = getNewSDK()
}

func initPKCS11Context(modulePath string) (*p11.Ctx, error) {
	context := p11.New(modulePath)

	if context == nil {
		return nil, fmt.Errorf("unable to load PKCS#11 module")
	}

	// May need to run it once
	//err := context.InitToken(0, pin, "") // should use Slot0 as label
	//if err != nil {
	//	return nil, err
	//}

	err := context.Initialize()
	return context, err
}

func getNewSDK() *sdk.SDK {
	s, err := sdk.NewSDK(module, tokenLabel, pin)
	if err != nil {
		fmt.Errorf("could get new SDK for module: %s; token label: %s; pin: %s", module, tokenLabel, pin)
		panic(err)
	}
	return s
}

func FailOnErr(t *testing.T, e error, msg string) {
	if e != nil {
		t.Fatalf("Fatal on error, %s, %+v", msg, e)
	}
}

func FailOnFlag(t *testing.T, flag bool, params ...interface{}) {
	if flag {
		t.Fatalf("Fail on falg, %v", params)
	}
}
