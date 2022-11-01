package btctran_azkv

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
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
// Return: signed raw tx, tx hash, from addr, error
func BuildTx(req *TxReq, sdk *kmssdk.SDK, networkParams *chaincfg.Params) (string, string, string, error) { // (signed raw tx, tx hash, fromAddress, error)
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
		return "", "", "", err
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
		return "", "", "", err
	}
	req.ToAddrPubKey = toAddrPubKey
	req.ToAddr = toAddr

	// new wired BTC tx
	msgTx, err := NewTx()
	if err != nil {
		return "", "", "", err
	}

	// coin selection for inputs and add to tx
	amountSelected, gasFee, err := SelectCoins(msgTx, req /*fromAddr, req.Amount*/, networkParams)
	if err != nil {
		return "", "", "", err
	}

	// add outputs to tx
	fromAddrScript, err := AddOutputs(msgTx, req /*toAddr, fromAddr, req.Amount,*/, amountSelected, gasFee, networkParams)
	if err != nil {
		return "", "", "", err
	}

	// signature for each input by KMS provider per sdk
	err = SignTx(msgTx, req, fromAddrScript, sdk)
	if err != nil {
		return "", "", "", err
	}

	// veryfy
	flags := txscript.StandardVerifyFlags
	for i, _ := range msgTx.TxIn {
		vm, err := txscript.NewEngine(*fromAddrScript, msgTx, i, flags, nil, nil, req.InputValues[i])
		if err != nil {
			return "", "", "", err
		}
		err = vm.Execute()
		if err != nil {
			return "", "", "", err
		}
	}

	// serialize and encode to string
	var buf bytes.Buffer
	msgTx.Serialize(&buf)
	// the following should work too
	//buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	//msgTx.Serialize(buf)

	signedTxReq := hex.EncodeToString(buf.Bytes())

	return signedTxReq, msgTx.TxHash().String(), fromAddr, nil
}

type MultisigTxReq struct {
	From   []kmssdk.KeyLabel // n = 3
	M      int               // m = 2
	To     kmssdk.KeyLabel
	Amount int64
	// the following field are populated by BuildTx
	RedeemScript []byte
	ScriptPubKey []byte
	FromAddr     string
	ToAddrPubKey *btcutil.AddressPubKey
	ToAddr       string
	InputValues  []int64
}

// Return: signed raw tx, tx hash, from addr, error
func BuildMultisigTx(req *MultisigTxReq, sdk *kmssdk.SDK, networkParams *chaincfg.Params) (string, string, string, error) {
	redeemScript, scriptPubKey, fromAddr, err := GetMultisigRedeemScript(req.From, req.M, sdk, networkParams)
	if err != nil {
		return "", "", "", err
	}
	req.RedeemScript = redeemScript
	req.ScriptPubKey = scriptPubKey
	req.FromAddr = fromAddr

	toAddrPubKey, toAddr, err := GetBech32AddressPublicKey(req.To, sdk, networkParams)
	if err != nil {
		return "", "", "", err
	}
	req.ToAddrPubKey = toAddrPubKey
	req.ToAddr = toAddr

	// new wired BTC tx
	msgTx, err := NewTx()
	if err != nil {
		return "", "", "", err
	}

	// coin selection for inputs and add to tx
	amountSelected, gasFee, err := SelectCoinsForMultisig(msgTx, req /*fromAddr, req.Amount*/, networkParams)
	if err != nil {
		return "", "", "", err
	}

	// add outputs to tx
	err = AddOutputsForMultisig(msgTx, req, amountSelected, gasFee, networkParams)
	if err != nil {
		return "", "", "", err
	}

	// sign: for each input, we get the first M signatures
	for i, txIn := range msgTx.TxIn {
		sigBuilder := txscript.NewScriptBuilder()
		sigBuilder.AddOp(txscript.OP_0) // same as OP_0
		//sigBuilder.AddOp(txscript.OP_FALSE) // same as OP_0
		for j := 0; j < req.M; j++ { // j-th signature
			sig, err := RawTxInSignature(msgTx, i, redeemScript, txscript.SigHashAll, req.From[j], sdk)
			if err != nil {
				return "", "", "", fmt.Errorf("failed to get %dth signature: %+v", j, err)
			}
			sigBuilder.AddData(sig)
		}
		sigBuilder.AddData(redeemScript)
		sigScript, err := sigBuilder.Script()
		if err != nil {
			return "", "", "", err
		}
		txIn.SignatureScript = sigScript
	}

	// veryfy
	//flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
	//	txscript.ScriptStrictMultiSig |
	//	txscript.ScriptDiscourageUpgradableNops
	flags := txscript.StandardVerifyFlags
	//lockingScript, err := GetLockingScript(req.FromAddr, networkParams)
	for i, _ := range msgTx.TxIn {
		vm, err := txscript.NewEngine(req.ScriptPubKey, msgTx, i, flags, nil, nil, req.InputValues[i])
		if err != nil {
			return "", "", "", err
		}
		err = vm.Execute()
		if err != nil {
			return "", "", "", fmt.Errorf("failed to verify transaction: %+v", err)
		}
	}
	// serialize and encode to string
	var buf bytes.Buffer
	msgTx.Serialize(&buf)
	// the following should work too
	//buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	//msgTx.Serialize(buf)

	signedTxReq := hex.EncodeToString(buf.Bytes())

	return signedTxReq, msgTx.TxHash().String(), fromAddr, nil
}

// Return: signed raw tx, tx hash, from addr, error
func BuildSegWitMultisigTx(req *MultisigTxReq, sdk *kmssdk.SDK, networkParams *chaincfg.Params) (string, string, string, error) {
	redeemScript, scriptPubKey, fromAddr, err := GetSegWitMultisigRedeemScript(req.From, req.M, sdk, networkParams)
	if err != nil {
		return "", "", "", err
	}
	req.RedeemScript = redeemScript
	req.ScriptPubKey = scriptPubKey
	req.FromAddr = fromAddr

	toAddrPubKey, toAddr, err := GetBech32AddressPublicKey(req.To, sdk, networkParams)
	if err != nil {
		return "", "", "", err
	}
	req.ToAddrPubKey = toAddrPubKey
	req.ToAddr = toAddr

	// new wired BTC tx
	msgTx, err := NewTx()
	if err != nil {
		return "", "", "", err
	}

	// coin selection for inputs and add to tx
	amountSelected, gasFee, err := SelectCoinsForMultisig(msgTx, req /*fromAddr, req.Amount*/, networkParams)
	if err != nil {
		return "", "", "", err
	}

	// add outputs to tx
	err = AddOutputsForSegWitMultisig(msgTx, req, amountSelected, gasFee, networkParams)
	if err != nil {
		return "", "", "", err
	}

	// sign: for each input, we get the first M signatures
	/*
		lockingScript, err := GetSegWitLockingScript(req.FromAddr, networkParams)
		if err != nil {
			return "", "", "", err
		}
		if !txscript.IsPayToWitnessScriptHash(*lockingScript) {
			return "", "", "", fmt.Errorf("not P2WSH")
		}
	*/

	txSignHashes := txscript.NewTxSigHashes(msgTx)
	for i, txIn := range msgTx.TxIn {
		witness := wire.TxWitness{}
		witness = append(witness, nil) // nil instead of []byte{0}
		for j := 0; j < req.M; j++ {   // j-th signature
			sig, err := RawTxInMultisigWitnessSignature(
				msgTx,
				txSignHashes,
				i,
				req.InputValues[i],
				redeemScript,
				txscript.SigHashAll, req.From[j], sdk)
			if err != nil {
				return "", "", "", fmt.Errorf("failed to get %dth signature: %+v", j, err)
			}
			witness = append(witness, sig)
		}
		witness = append(witness, redeemScript)
		txIn.Witness = witness //wire.TxWitness{sigScript, redeemScript}
	}
	// veryfy
	/*
		flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
			txscript.ScriptStrictMultiSig |
			txscript.ScriptDiscourageUpgradableNops
	*/
	flags := txscript.StandardVerifyFlags
	//lockingScript, err := GetSegWitLockingScript(req.FromAddr, networkParams)
	//if err != nil {
	//	return "", "", "", err
	//}
	for i, _ := range msgTx.TxIn {
		vm, err := txscript.NewEngine(req.ScriptPubKey, msgTx, i, flags, nil, nil, req.InputValues[i])
		//vm, err := txscript.NewEngine(*lockingScript, msgTx, i, flags, nil, nil, req.InputValues[i])
		if err != nil {
			return "", "", "", fmt.Errorf("failed to verify transaction NewEngine: %+v", err)
		}
		err = vm.Execute()
		if err != nil {
			return "", "", "", fmt.Errorf("failed to verify transaction Excute: %+v", err)
		}
	}
	// serialize and encode to string
	var buf bytes.Buffer
	msgTx.Serialize(&buf)
	// the following should work too
	//buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	//msgTx.Serialize(buf)

	signedTxReq := hex.EncodeToString(buf.Bytes())

	return signedTxReq, msgTx.TxHash().String(), fromAddr, nil
}
