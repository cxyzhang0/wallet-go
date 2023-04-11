package test

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/Azure/go-autorest/autorest/azure"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"os"
	"testing"
)

// subscriptionId: dc934097-fb2c-4e79-ab45-4427a2363216
// "tenantId": "b0c970c0-191d-4289-9971-e961c7b6e8d2"
// directoryId: b0c970c0-191d-4289-9971-e961c7b6e8d2
// "name": "cxyzhang@gmail.com"
// vaultURI: https://szkv1.vault.azure.net/
// keyValutName: szkv1
// user principal name: cxyzhang_gmail.com#EXT#@cxyzhanggmail.onmicrosoft.com
// object id: ea599c48-bd9c-46a5-aa73-7a2d6effaadf // is this the principal id?
// key id: https://szkv1.vault.azure.net/keys/secp256k1-hsm-1
// az ad sp create-for-rbac --sdk-auth > azureauth.json
var (
	//azureAuthLocation = "/Users/johnz/Project/wf-innovation/wallet-go/azkv/sdk/azureauth.json"
	clientId     = "3757699e-7d6d-4ba0-a584-bbe9d00fcfeb"
	clientSecret = "gbg8Q~AU3ShJ8C7ymlz6Ze.peBpzLBPFLFgRkch8"
	//clientSecret = "bcy8Q~N6aYbXRV7Yrf9Dyrvy.JXcJH7d57Y5Abm7"  gbg8Q~AU3ShJ8C7ymlz6Ze.peBpzLBPFLFgRkch8
	tenantId  = "b0c970c0-191d-4289-9971-e961c7b6e8d2"
	vaultName = "szkv1"
	keyName   = "secp256k1-soft-1"
	//keyName      = "secp256k1-hsm-1"
	vaultURI string // https://szkv1.vault.azure.net/
	_sdk     *kmssdk.SDK
)

func init() {
	//os.Setenv("AZURE_AUTH_LOCATION", azureAuthLocation)
	os.Setenv("AZURE_TENANT_ID", tenantId)
	os.Setenv("AZURE_CLIENT_ID", clientId)
	os.Setenv("AZURE_CLIENT_SECRET", clientSecret)

	vaultURI = fmt.Sprintf("https://%s.%s", vaultName, azure.PublicCloud.KeyVaultDNSSuffix)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client := azkeys.NewClient(vaultURI, cred, nil)

	//client := keyvault.New()
	//authorizer, err := auth.NewAuthorizerFromFile()
	//if err != nil {
	//	panic(err)
	//}
	//
	//client.Authorizer = authorizer

	_sdk = kmssdk.NewSDK(client)
}

func FailOnErr(t *testing.T, e error, msg string) {
	if e != nil {
		t.Fatalf("Fatal on error, %s, %v", msg, e)
	}
}

func FailOnFlag(t *testing.T, flag bool, params ...interface{}) {
	if flag {
		t.Fatalf("Fail on falg, %v", params)
	}
}
