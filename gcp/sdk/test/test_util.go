package test

import (
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
)

var (
	_sdk *gcpsdk.SDK
)

func init() {
	var err error
	_sdk, err = gcpsdk.NewSDK()
	if err != nil {
		panic(err)
	}
}
