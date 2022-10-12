package btctran_azkv

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
)

type TxReq struct {
	From         kmssdk.KeyLabel
	IsFromLegacy bool
	To           kmssdk.KeyLabel
	IsToLegacy   bool
	Amount       int64
	// the following field are populated by BuildTx
	FromAddrPubKey *btcutil.AddressPubKey
	FromAddr       string
	ToAddrPubKey   *btcutil.AddressPubKey
	ToAddr         string
	InputValues    []int64
}

// BuildTx returns a signed Bitcoin transaction request and its hash, ready to be
// submitted to the blockchain.
// req only needs From, To and Amount. The rest is populated by this func.
// sdk points to a specific KMS SDK, here it is SDK for gcp KMS. The code
// should be the same for pkcs11 SDK. So it is pluggable.
// networkParams points to a specific BTC network: mainnet or testnet.
func BuildTx(req *TxReq, sdk *kmssdk.SDK, networkParams *chaincfg.Params) (string, string, error) { // (signed raw tx, tx hash, error)
	// get from addresses
	var fromAddrPubKey, toAddrPubKey *btcutil.AddressPubKey
	var fromAddr, toAddr string
	var err error
	// get from addresses
	if req.IsFromLegacy {
		fromAddrPubKey, fromAddr, err = GetLegacyAddressPubKey(req.From, sdk, networkParams)
	} else {
		fromAddrPubKey, fromAddr, err = GetBech32AddressPublicKey(req.From, sdk, networkParams)
	}
	if err != nil {
		return "", "", err
	}
	req.FromAddrPubKey = fromAddrPubKey
	req.FromAddr = fromAddr

	// get to addresses
	if req.IsToLegacy {
		toAddrPubKey, toAddr, err = GetLegacyAddressPubKey(req.To, sdk, networkParams)
	} else {
		toAddrPubKey, toAddr, err = GetBech32AddressPublicKey(req.To, sdk, networkParams)
	}
	if err != nil {
		return "", "", err
	}
	req.ToAddrPubKey = toAddrPubKey
	req.ToAddr = toAddr

	// new wired BTC tx
	msgTx, err := NewTx()
	if err != nil {
		return "", "", err
	}

	// coin selection for inputs and add to tx
	amountSelected, gasFee, err := SelectCoins(msgTx, req /*fromAddr, req.Amount*/)
	if err != nil {
		return "", "", err
	}

	// add outputs to tx
	fromAddrScript, err := AddOutputs(msgTx, req /*toAddr, fromAddr, req.Amount,*/, amountSelected, gasFee, networkParams)
	if err != nil {
		return "", "", err
	}

	// signature for each input by KMS provider per sdk
	err = SignTx(msgTx, req, fromAddrScript, sdk)
	if err != nil {
		return "", "", err
	}

	// veryfy
	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(*fromAddrScript, msgTx, len(msgTx.TxIn)-1, flags, nil, nil, amountSelected)
	if err != nil {
		return "", "", err
	}
	err = vm.Execute()
	if err != nil {
		return "", "", err
	}

	// serialize and encode to string
	var buf bytes.Buffer
	msgTx.Serialize(&buf)
	// the following should work too
	//buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	//msgTx.Serialize(buf)

	signedTxReq := hex.EncodeToString(buf.Bytes())

	return signedTxReq, msgTx.TxHash().String(), nil
}
