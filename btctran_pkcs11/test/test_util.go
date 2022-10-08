package test

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cxyzhang0/wallet-go/btctran_pkcs11"
	pkcs11sdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"os"
)

var (
	//module          = "/Users/johnz/futurex/wf_pkcs11_mac/libfxpkcs11-debug.dylib"
	//tokenLabel      = "us01hsm01test.virtucrypt.com:592"
	//privateKeyLabel = "projects/quantum-pilot-360000/locations/us-west1/keyRings/WIM-test/cryptoKeys/secp256k1-hsm-1/cryptoKeyVersions/1"
	//pin             = "safest"
	module        = "/usr/local/lib/softhsm/libsofthsm2.so"
	tokenLabel    = "Slot Token 0"
	pin           = "5678"
	_sdk          *pkcs11sdk.SDK
	networkParams *chaincfg.Params
)

func init() {

	os.Setenv("SOFTHSM2_CONF", "/Users/johnz/Project/wf-innovation/wallet-go/btctran_pkcs11/test/softhsm2.conf")

	os.Setenv("FXPKCS11_CFG", "/Users/johnz/futurex/wf_pkcs11_mac/fxpkcs11.cfg")

	var err error
	_sdk = getNewSDK()
	if err != nil {
		panic(err)
	}

	networkParams, err = btctran_pkcs11.GetBTCNetworkParams()
	if err != nil {
		panic(err)
	}

}

func getNewSDK() *pkcs11sdk.SDK {
	s, err := pkcs11sdk.NewSDK(module, tokenLabel, pin)
	if err != nil {
		fmt.Errorf("could get new SDK for module %s; token label %s; pin %s: %+v", module, tokenLabel, pin, err)
		panic(err)
	}
	return s
}
