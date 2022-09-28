package btctran_gcp

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
)

type TxReq struct {
	From   gcpsdk.KeyLabel
	To     gcpsdk.KeyLabel
	Amount int64
	// the following field are populated by BuildTx
	FromAddrPubKey *btcutil.AddressPubKey
	FromAddr       string
	ToAddrPubKey   *btcutil.AddressPubKey
	ToAddr         string
}

func BuildTx(req *TxReq, sdk *gcpsdk.SDK, networkParams *chaincfg.Params) (string, string, error) {
	fromAddrPubKey, fromAddr, err := GetAddressPubKey(req.From, sdk, networkParams)
	if err != nil {
		return "", "", err
	}
	req.FromAddrPubKey = fromAddrPubKey
	req.FromAddr = fromAddr

	toAddrPubKey, toAddr, err := GetAddressPubKey(req.To, sdk, networkParams)
	if err != nil {
		return "", "", err
	}
	req.ToAddrPubKey = toAddrPubKey
	req.ToAddr = toAddr

	msgTx, err := NewTx()
	if err != nil {
		return "", "", err
	}

	// coin selection for inputs
	amountSelected, gasFee, err := SelectCoins(msgTx, req /*fromAddr, req.Amount*/)
	if err != nil {
		return "", "", err
	}

	// outputs
	fromAddrScript, err := AddOutputs(msgTx, req, /*toAddr, fromAddr, req.Amount,*/ amountSelected, gasFee, networkParams)
	if err != nil {
		return "", "", err
	}

	// signature
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

	var buf bytes.Buffer
	msgTx.Serialize(&buf)
	// the following should work too
	//buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	//msgTx.Serialize(buf)

	signedTxReq := hex.EncodeToString(buf.Bytes())

	return signedTxReq, msgTx.TxHash().String(), nil
}
