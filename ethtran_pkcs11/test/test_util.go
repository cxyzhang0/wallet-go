package test

import (
	"context"
	"fmt"

	kmssdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"github.com/ethereum/go-ethereum/params"
	ubi "gitlab.com/Blockdaemon/ubiquity/ubiquity-go-client/v1/pkg/client"
	"os"
)

var (
	//module          = "/Users/johnz/futurex/wf_pkcs11_mac/libfxpkcs11-debug.dylib"
	//tokenLabel      = "us01hsm01test.virtucrypt.com:592"
	//privateKeyLabel = "projects/quantum-pilot-360000/locations/us-west1/keyRings/WIM-test/cryptoKeys/secp256k1-hsm-1/cryptoKeyVersions/1"
	//pin             = "safest"
	module     = "/usr/local/lib/softhsm/libsofthsm2.so"
	tokenLabel = "Slot Token 0"
	pin        = "5678"
	_sdk       *kmssdk.SDK
	//networkParams *chaincfg.Params

	ubiURL         = "https://svc.blockdaemon.com/universal/v1"
	ubiPlatform    = "ethereum"
	ubiNetwork     = "goerli"
	ubiAccessToken = "0CufGocrOhFvDLWTfbX5kHaOCRlmFmD7BjW-TrY0mRiHNs21"
	ubiAPIClient   *ubi.APIClient
	ubiCtx         context.Context
	chainConfig    *params.ChainConfig = params.GoerliChainConfig
)

func init() {

	os.Setenv("SOFTHSM2_CONF", "/Users/johnz/Project/wf-innovation/wallet-go/btctran_pkcs11/test/softhsm2.conf")

	os.Setenv("FXPKCS11_CFG", "/Users/johnz/futurex/wf_pkcs11_mac/fxpkcs11.cfg")

	var err error
	_sdk = getNewSDK()
	if err != nil {
		panic(err)
	}

	ubiConfig := ubi.NewConfiguration()
	ubiConfig.Servers = ubi.ServerConfigurations{
		{
			URL:         ubiURL,
			Description: "Production endpoint",
		},
	}
	ubiAPIClient = ubi.NewAPIClient(ubiConfig)
	ubiCtx = context.WithValue(context.Background(), ubi.ContextAccessToken, ubiAccessToken)
	//fmt.Printf("ubi ctx %+v", ubiCtx)

}

func getNewSDK() *kmssdk.SDK {
	s, err := kmssdk.NewSDK(module, tokenLabel, pin)
	if err != nil {
		fmt.Errorf("could get new SDK for module %s; token label %s; pin %s: %+v", module, tokenLabel, pin, err)
		panic(err)
	}
	return s
}
