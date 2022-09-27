package test

import "github.com/cxyzhang0/wallet-go/gcp/sdk"

var (
	_sdk *sdk.SDK
)

func init() {
	var err error
	_sdk, err = sdk.NewSDK()
	if err != nil {
		panic(err)
	}
}
