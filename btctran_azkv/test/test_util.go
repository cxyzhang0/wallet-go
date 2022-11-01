package test

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/btcsuite/btcd/chaincfg"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"github.com/cxyzhang0/wallet-go/btctran_azkv"
	"os"
)

var (
	clientId      = "3757699e-7d6d-4ba0-a584-bbe9d00fcfeb"
	clientSecret  = "bcy8Q~N6aYbXRV7Yrf9Dyrvy.JXcJH7d57Y5Abm7"
	tenantId      = "b0c970c0-191d-4289-9971-e961c7b6e8d2"
	vaultName     = "szkv1"
	vaultURI      string // https://szkv1.vault.azure.net/
	_sdk          *kmssdk.SDK
	networkParams *chaincfg.Params

	//wsHost = "socket.blockcypher.com"
	//wsPath = "v1/btc/test3?token=e905d13ae51748e2b618da1ba4ce0458"
)

func init() {
	os.Setenv("AAZURE_TENANT_ID", tenantId)
	os.Setenv("AZURE_CLIENT_ID", clientId)
	os.Setenv("AZURE_CLIENT_SECRET", clientSecret)

	vaultURI = fmt.Sprintf("https://%s.%s", vaultName, azure.PublicCloud.KeyVaultDNSSuffix)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client := azkeys.NewClient(vaultURI, cred, nil)

	_sdk = kmssdk.NewSDK(client)
	if err != nil {
		panic(err)
	}

	networkParams, err = btctran_azkv.GetBTCNetworkParams()
	if err != nil {
		panic(err)
	}

}
