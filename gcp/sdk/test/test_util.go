package test

import (
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
	"testing"
)

var (
	gcpCredentials         = "GOOGLE_APPLICATION_CREDENTIALS"
	gcpCredentialsLocation = "/Users/johnz/Project/wf-innovation/certs/coreblock-367317-554ec1b320cb.json" // on mac studio
	//gcpCredentialsLocation = "/Users/johnz/Project/wf-innovation/certs/quantum-pilot-360000-df0aea9a2be4.json"
	//gcpCredentialsLocation = "/Users/johnz/Project/wf-innovation/certs/coreblock-367317-cba9932853e6.json" // on mac studio
	//gcpCredentialsLocation = "/Users/johnz/Project/wf-innovation/certs/coreblock-367317-95abf9200b1d.json"  // on mac pro
	//gcpCredentialsLocation = "/Users/johnz/Project/wf-innovation/certs/coreblock-367317-81d89892117b.json"
	gcpProject  = "coreblock-367317"
	gcpLocation = "us-west1"
	gcpKeyRing  = "sean-1"
	//gcpProject  = "quantum-pilot-360000"
	//gcpLocation = "us-west2"
	//gcpKeyRing  = "WIM-test-2"
	//gcpKeyRing  = "sean-hsm"
	_sdk *gcpsdk.SDK
)

func init() {
	var err error
	_sdk, err = gcpsdk.NewSDK()
	if err != nil {
		panic(err)
	}

	// TODO: this process level ENV does not work.
	// 	So need export $gcpCredentials=$gcpCredentialsLocation in ~/.zshenv,
	// 	e.g., export GOOGLE_APPLICATION_CREDENTIALS=/Users/johnz/Project/wf-innovation/certs/coreblock-367317-554ec1b320cb.json
	//os.Setenv(gcpCredentials, gcpCredentialsLocation)
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
