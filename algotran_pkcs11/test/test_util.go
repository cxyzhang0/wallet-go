package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	kmssdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"testing"
)

var (
	//module          = "/Users/johnz/futurex/wf_pkcs11_mac/libfxpkcs11-debug.dylib"
	//tokenLabel      = "us01hsm01test.virtucrypt.com:592"
	//pin             = "safest"
	module     = "/usr/local/lib/softhsm/libsofthsm2.so"
	tokenLabel = "Slot Token 0"
	pin        = "5678"
	_sdk       *kmssdk.SDK
	/**
	prerequisite: start the Algorand sandbox containers
	in /Users/johnz/Project/wf-innovation/algorand/sandbox,
	./sandbox up
	*/
	algodAddress = "http://localhost:4001"
	algodToken   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	algodClient  *algod.Client
)

func init() {
	_sdk = getNewSDK()

	client, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		panic(err)
	}
	algodClient = client
}

func getNewSDK() *kmssdk.SDK {
	s, err := kmssdk.NewSDK(module, tokenLabel, pin)
	if err != nil {
		fmt.Errorf("could get new SDK for module %s; token label %s; pin %s: %+v", module, tokenLabel, pin, err)
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

// prettyPrint prints Go structs
func prettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

// printAssetHolding utility to print asset holding for account
func printAssetHolding(assetID uint64, account string, client *algod.Client) {

	act, err := client.AccountInformation(account).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to get account information: %s\n", err)
		return
	}
	for _, assetholding := range act.Assets {
		if assetID == assetholding.AssetId {
			prettyPrint(assetholding)
			break
		}
	}
}

// printCreatedAsset utility to print created assert for account
func printCreatedAsset(assetID uint64, account string, client *algod.Client) {

	act, err := client.AccountInformation(account).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to get account information: %s\n", err)
		return
	}
	for _, asset := range act.CreatedAssets {
		if assetID == asset.Index {
			prettyPrint(asset)
			break
		}
	}
}
