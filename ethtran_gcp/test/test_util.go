package test

import (
	"context"
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
	"github.com/ethereum/go-ethereum/params"
	ubi "gitlab.com/Blockdaemon/ubiquity/ubiquity-go-client/v1/pkg/client"
)

var (
	_sdk           *gcpsdk.SDK
	ubiURL         = "https://svc.blockdaemon.com/universal/v1"
	ubiPlatform    = "ethereum"
	ubiNetwork     = "goerli"
	ubiAccessToken = "0CufGocrOhFvDLWTfbX5kHaOCRlmFmD7BjW-TrY0mRiHNs21"
	ubiAPIClient   *ubi.APIClient
	ubiCtx         context.Context
	chainConfig    *params.ChainConfig = params.GoerliChainConfig
)

func init() {
	var err error
	_sdk, err = gcpsdk.NewSDK()
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
