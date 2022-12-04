package test

import (
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
	"os"
	"testing"
)

var (
	gcpCredentials = "GOOGLE_APPLICATION_CREDENTIALS"
	//gcpCredentialsLocation = "/Users/johnz/Project/wf-innovation/certs/quantum-pilot-360000-df0aea9a2be4.json"
	gcpCredentialsLocation = "/Users/johnz/Project/wf-innovation/certs/coreblock-367317-95abf9200b1d.json"
	//gcpCredentialsLocation = "/Users/johnz/Project/wf-innovation/certs/coreblock-367317-81d89892117b.json"
	gcpProject  = "coreblock-367317"
	gcpLocation = "us-west1"
	gcpKeyRing  = "sean-1"
	//gcpKeyRing  = "sean-hsm"
	_sdk *gcpsdk.SDK
)

func init() {
	var err error
	_sdk, err = gcpsdk.NewSDK()
	if err != nil {
		panic(err)
	}

	os.Setenv(gcpCredentials, gcpCredentialsLocation)
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
