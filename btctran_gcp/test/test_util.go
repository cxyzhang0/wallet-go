package test

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cxyzhang0/wallet-go/btctran_gcp"
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
)

var (
	_sdk          *gcpsdk.SDK
	networkParams *chaincfg.Params
)

func init() {
	var err error
	_sdk, err = gcpsdk.NewSDK()
	if err != nil {
		panic(err)
	}

	networkParams, err = btctran_gcp.GetBTCNetworkParams()
	if err != nil {
		panic(err)
	}

}
