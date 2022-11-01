package test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/Azure/go-autorest/autorest/azure"

	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"github.com/ethereum/go-ethereum/params"
	ubi "gitlab.com/Blockdaemon/ubiquity/ubiquity-go-client/v1/pkg/client"
	"os"
)

var (
	clientId     = "3757699e-7d6d-4ba0-a584-bbe9d00fcfeb"
	clientSecret = "bcy8Q~N6aYbXRV7Yrf9Dyrvy.JXcJH7d57Y5Abm7"
	tenantId     = "b0c970c0-191d-4289-9971-e961c7b6e8d2"
	vaultName    = "szkv1"
	vaultURI     string // https://szkv1.vault.azure.net/
	_sdk         *kmssdk.SDK

	ubiURL         = "https://svc.blockdaemon.com/universal/v1"
	ubiNativeURL   = "https://svc.blockdaemon.com"
	ubiPlatform    = "ethereum"
	ubiNetwork     = "goerli"
	ubiAccessToken = "0CufGocrOhFvDLWTfbX5kHaOCRlmFmD7BjW-TrY0mRiHNs21"
	ubiAPIClient   *ubi.APIClient
	ubiCtx         context.Context
	chainConfig    *params.ChainConfig = params.GoerliChainConfig
	quicknodeURL                       = "https://dark-cosmopolitan-seed.ethereum-goerli.discover.quiknode.pro/dc6e17a2cfbc338c5e59511eed170c97cc7cfa15/"
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
