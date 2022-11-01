package btctran_azkv

import (
	"errors"
	"fmt"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"github.com/cxyzhang0/wallet-go/btctran_azkv/test/config"
)

var conf = config.ParseConfig()
var bcyAPI = gobcy.API{Token: conf.Blockcypher.Token, Coin: conf.Blockcypher.Coin, Chain: conf.Blockcypher.Chain}
var Conf = conf
var BcyAPI = bcyAPI

func GetLegacyAddressPubKey(keyLabel kmssdk.KeyLabel, sdk *kmssdk.SDK, networkParams *chaincfg.Params) (*btcutil.AddressPubKey, string, error) { // address pub key, address, error
	pubkey, err := sdk.GetECDSAPublicKey(keyLabel)
	if err != nil {
		return nil, "", err
	}

	key := btcec.PublicKey(*pubkey)
	addrPubKey, err := btcutil.NewAddressPubKey(key.SerializeUncompressed(), networkParams)
	if err != nil {
		return nil, "", err
	}

	addr := addrPubKey.EncodeAddress()

	return addrPubKey, addr, nil
}

func GetBech32AddressPublicKey(keyLabel kmssdk.KeyLabel, sdk *kmssdk.SDK, networkParams *chaincfg.Params) (*btcutil.AddressPubKey, string, error) { // address pub key, address, error
	pubkey, err := sdk.GetECDSAPublicKey(keyLabel)
	if err != nil {
		return nil, "", err
	}

	key := btcec.PublicKey(*pubkey)

	addrPubKey, err := btcutil.NewAddressPubKey(key.SerializeUncompressed(), networkParams)

	witnessProg := btcutil.Hash160(key.SerializeCompressed())
	addrWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, networkParams)
	if err != nil {
		return nil, "", err
	}

	addr := addrWitnessPubKeyHash.EncodeAddress()

	return addrPubKey, addr, nil
}

// 2 of 3 multisig P2SH script - non Taproot
func GetP2SHRedeemScript(keyLabel1 kmssdk.KeyLabel, keyLabel2 kmssdk.KeyLabel, keyLabel3 kmssdk.KeyLabel, sdk *kmssdk.SDK, networkParams *chaincfg.Params) ([]byte, []byte, string, error) { // redeemScript, scriptPubKey, address, error
	pubkey1, err := sdk.GetECDSAPublicKey(keyLabel1)
	if err != nil {
		return nil, nil, "", err
	}

	pubkey2, err := sdk.GetECDSAPublicKey(keyLabel2)
	if err != nil {
		return nil, nil, "", err
	}

	pubkey3, err := sdk.GetECDSAPublicKey(keyLabel3)
	if err != nil {
		return nil, nil, "", err
	}

	key1 := btcec.PublicKey(*pubkey1)
	key2 := btcec.PublicKey(*pubkey2)
	key3 := btcec.PublicKey(*pubkey3)

	// create redeem script for 2 of 3 multisig
	builder := txscript.NewScriptBuilder()
	// OP 2 required sigs
	builder.AddOp(txscript.OP_2)
	// 3 pub keys
	builder.AddData(key1.SerializeCompressed()).AddData(key2.SerializeCompressed()).AddData(key3.SerializeCompressed())
	// OP 3 pub keys
	builder.AddOp(txscript.OP_3)
	// OP check multisig
	builder.AddOp(txscript.OP_CHECKMULTISIG)
	// redeem script
	redeemScript, err := builder.Script()
	if err != nil {
		return nil, nil, "", err
	}
	redeemScriptHash := btcutil.Hash160(redeemScript)
	address, err := btcutil.NewAddressScriptHashFromHash(redeemScriptHash, networkParams)
	if err != nil {
		return nil, nil, "", err
	}

	scriptPubKey, err := txscript.PayToAddrScript(address)
	if err != nil {
		return nil, nil, "", err
	}

	addr := address.EncodeAddress() // base58

	// encode redeemScriptHash to get 1627a5a154b1cc8543353fdccdf395f7aaf2b1c0 as in SCRIPTPUBKEY (ASM)
	// e.g., OP_HASH160 OP_PUSHBYTES_20 1627a5a154b1cc8543353fdccdf395f7aaf2b1c0 OP_EQUAL
	/*
		str := hex.EncodeToString(redeemScriptHash) // this one shows up in SCRIPTPUBKEY (ASM)
		fmt.Sprintf(str)
	*/

	// from add we can get redeemScriptHash as follows
	/*
		address, err := btcutil.DecodeAddress(addr, networkParams)
		if err != nil {
		}
		hash := address.ScriptAddress()  // this is the same as redeemScriptHash
		fmt.Sprintf("%x", hash)
	*/
	return redeemScript, scriptPubKey, addr, nil
}

// GetSegWitMultisigRedeemScript returns the SegWit multisig (P2WSH)
// redeemScript (witnessScript): OP_2 pubkey1 pubkey2 pubkey2 OP_3 OP_CHECKMULTISIG
// scriptPubKey: OP_0 hash32(redeemScript)
// address: Bech32 encoded address string from scriptPubKey
// Note: hash32(redeemScript) is also known as witness program
func GetSegWitMultisigRedeemScript(keyLabels []kmssdk.KeyLabel, m int, sdk *kmssdk.SDK, networkParams *chaincfg.Params) ([]byte, []byte, string, error) { // address pub key, address, error
	addressPubKeys := make([]*btcutil.AddressPubKey, len(keyLabels))

	for i, keylabel := range keyLabels {
		pubkey, err := sdk.GetECDSAPublicKey(keylabel)
		key := btcec.PublicKey(*pubkey)
		if err != nil {
			return nil, nil, "", err
		}

		addressPubKey, err := btcutil.NewAddressPubKey(key.SerializeCompressed(), networkParams)
		if err != nil {
			return nil, nil, "", err
		}

		addressPubKeys[i] = addressPubKey
	}

	redeemScript, err := txscript.MultiSigScript(addressPubKeys, m)
	if err != nil {
		return nil, nil, "", err
	}

	// P2WSH
	redeemScriptHash, err := kmssdk.SecureHashByteArray(redeemScript)
	if err != nil {
		return nil, nil, "", err
	}

	address, err := btcutil.NewAddressWitnessScriptHash(redeemScriptHash, networkParams)
	if err != nil {
		return nil, nil, "", err
	}

	scriptPubKey, err := txscript.PayToAddrScript(address)
	if err != nil {
		return nil, nil, "", err
	}

	addr := address.EncodeAddress() // base58

	return redeemScript, scriptPubKey, addr, nil
}

// It uses txscript.MultiSigScript(...)
func GetMultisigRedeemScript(keyLabels []kmssdk.KeyLabel, m int, sdk *kmssdk.SDK, networkParams *chaincfg.Params) ([]byte, []byte, string, error) { // redeemScript, scriptPubKey, address, error
	addressPubKeys := make([]*btcutil.AddressPubKey, len(keyLabels))

	for i, keylabel := range keyLabels {
		pubkey, err := sdk.GetECDSAPublicKey(keylabel)
		key := btcec.PublicKey(*pubkey)
		if err != nil {
			return nil, nil, "", err
		}

		addressPubKey, err := btcutil.NewAddressPubKey(key.SerializeCompressed(), networkParams)
		if err != nil {
			return nil, nil, "", err
		}

		addressPubKeys[i] = addressPubKey
	}

	redeemScript, err := txscript.MultiSigScript(addressPubKeys, m)
	if err != nil {
		return nil, nil, "", err
	}

	redeemScriptHash := btcutil.Hash160(redeemScript)
	address, err := btcutil.NewAddressScriptHashFromHash(redeemScriptHash, networkParams)
	if err != nil {
		return nil, nil, "", err
	}

	scriptPubKey, err := txscript.PayToAddrScript(address)
	if err != nil {
		return nil, nil, "", err
	}

	addr := address.EncodeAddress() // base58

	// encode redeemScriptHash to get 1627a5a154b1cc8543353fdccdf395f7aaf2b1c0 as in SCRIPTPUBKEY (ASM)
	// e.g., OP_HASH160 OP_PUSHBYTES_20 1627a5a154b1cc8543353fdccdf395f7aaf2b1c0 OP_EQUAL
	/*
		str := hex.EncodeToString(redeemScriptHash) // this one shows up in SCRIPTPUBKEY (ASM)
		fmt.Sprintf(str)
	*/

	// from add we can get redeemScriptHash as follows
	/*
		address, err := btcutil.DecodeAddress(addr, networkParams)
		if err != nil {
		}
		hash := address.ScriptAddress()  // this is the same as redeemScriptHash
		fmt.Sprintf("%x", hash)
	*/
	return redeemScript, scriptPubKey, addr, nil
}

// It uses ScriptBuilder with the same result as above
func GetMultisigRedeemScript0(keyLabels []kmssdk.KeyLabel, m int, sdk *kmssdk.SDK, networkParams *chaincfg.Params) ([]byte, string, error) { // address pub key, address, error
	n := len(keyLabels)
	if m > n {
		return nil, "", fmt.Errorf("m %d must be <= n %d", m, n)
	}

	if n > 16 {
		return nil, "", fmt.Errorf("n %d must be <= 16", n)
	}

	// create redeem script for m of n multisig
	//opM := (m - 2) + 82
	//opN := (n - 2) + 82

	builder := txscript.NewScriptBuilder()
	// required sigs
	builder.AddInt64(int64(m))
	//builder.AddOp(byte(opM))
	// add pub keys
	for _, keylabel := range keyLabels {
		pubkey, err := sdk.GetECDSAPublicKey(keylabel)
		if err != nil {
			return nil, "", err
		}

		key := btcec.PublicKey(*pubkey)
		builder.AddData(key.SerializeCompressed())
	}

	// OP n pub keys
	builder.AddInt64(int64(n))
	//builder.AddOp(byte(opN))
	// OP check multisig
	builder.AddOp(txscript.OP_CHECKMULTISIG)
	// redeem script
	redeemScript, err := builder.Script()
	if err != nil {
		return nil, "", err
	}
	redeemScriptHash := btcutil.Hash160(redeemScript)
	addrScriptHash, err := btcutil.NewAddressScriptHashFromHash(redeemScriptHash, networkParams)
	if err != nil {
		return nil, "", err
	}

	addr := addrScriptHash.EncodeAddress() // base58

	// encode redeemScriptHash to get 1627a5a154b1cc8543353fdccdf395f7aaf2b1c0 as in SCRIPTPUBKEY (ASM)
	// e.g., OP_HASH160 OP_PUSHBYTES_20 1627a5a154b1cc8543353fdccdf395f7aaf2b1c0 OP_EQUAL
	/*
		str := hex.EncodeToString(redeemScriptHash) // this one shows up in SCRIPTPUBKEY (ASM)
		fmt.Sprintf(str)
	*/

	// from add we can get redeemScriptHash as follows
	/*
		address, err := btcutil.DecodeAddress(addr, networkParams)
		if err != nil {
		}
		hash := address.ScriptAddress()  // this is the same as redeemScriptHash
		fmt.Sprintf("%x", hash)
	*/
	return redeemScript, addr, nil
}

func NewTx() (*wire.MsgTx, error) {
	return wire.NewMsgTx(wire.TxVersion), nil
}

// ref: https://bitzuma.com/posts/making-sense-of-bitcoin-transaction-fees/
func EstimateFee(msgTx *wire.MsgTx, addr string, networkParams *chaincfg.Params) (int, error) {
	address, err := btcutil.DecodeAddress(addr, networkParams)
	if err != nil {
		return -1, err
	}

	scriptPubKey, err := txscript.PayToAddrScript(address)
	if err != nil {
		return -1, err
	}
	inputCount := len(msgTx.TxIn)
	outputCount := len(msgTx.TxOut)
	if outputCount < 2 {
		outputCount = 2
	}
	var txSize int
	switch {
	case txscript.IsPayToWitnessPubKeyHash(scriptPubKey):
		txSize = (42 + 272*inputCount + 128*outputCount) / 4
	case txscript.IsPayToWitnessScriptHash(scriptPubKey):
		txSize = (42+272*inputCount+128*outputCount)/4 + 50 // TODO: 50 is arbitrary for P2WSH
	case txscript.IsPayToScriptHash(scriptPubKey):
		txSize = inputCount*180 + outputCount*34 + 10 + 150 // TODO: 150 is arbitrary for P2SH
	default:
		txSize = inputCount*180 + outputCount*34 + 10
	}

	//txSize = inputCount*180 + outputCount*34 + 10 + 300 // 300 is arbitrary for P2SH
	return conf.Blockcypher.FeeRate * txSize, nil
}

func GetBTCNetworkParams() (*chaincfg.Params, error) {
	if conf.Blockcypher.Coin != "btc" {
		return nil, errors.New("Coin not supported: " + conf.Blockcypher.Coin)
	}
	if conf.Blockcypher.Chain == "test3" {
		return &chaincfg.TestNet3Params, nil
	} else {
		return &chaincfg.MainNetParams, nil
	}
}

func SelectCoins(msgTx *wire.MsgTx, req *TxReq, networkParams *chaincfg.Params) (int64, int, error) {
	// check the balance of fromAddr
	addr, err := bcyAPI.GetAddrBal(req.FromAddr, nil)
	if err != nil {
		return -1, -1, err
	}

	// initial sufficient balance check
	if int64(addr.Balance) <= req.Amount {
		return -1, -1, fmt.Errorf("Insufficent funds (initial check): balance %d <= amountToSend %d", addr.Balance, req.Amount)
		//return -1, -1, errors.New(fmt.Sprintf("Insufficent funds (initial check): balance %d <= amountToSend %d", addr.Balance, req.Amount))
	}

	// select coins
	var params = make(map[string]string, 0)
	params["unspentOnly"] = "true"
	params["includeScript"] = "true"
	utxos, err := bcyAPI.GetAddr(req.FromAddr, params)
	if err != nil {
		return -1, -1, err
	}
	amountSelected, feeAmount, err := SelectInputs(msgTx, utxos, 0, req, networkParams)
	if err != nil {
		return -1, -1, err
	}

	return amountSelected, feeAmount, nil
}

func SelectCoinsForMultisig(msgTx *wire.MsgTx, req *MultisigTxReq, networkParams *chaincfg.Params) (int64, int, error) {
	// check the balance of fromAddr
	addr, err := bcyAPI.GetAddrBal(req.FromAddr, nil)
	if err != nil {
		return -1, -1, err
	}

	// initial sufficient balance check
	if int64(addr.Balance) <= req.Amount {
		return -1, -1, fmt.Errorf("Insufficent funds (initial check): balance %d <= amountToSend %d", addr.Balance, req.Amount)
		//return -1, -1, errors.New(fmt.Sprintf("Insufficent funds (initial check): balance %d <= amountToSend %d", addr.Balance, req.Amount))
	}

	// select coins
	var params = make(map[string]string, 0)
	params["unspentOnly"] = "true"
	params["includeScript"] = "true"
	utxos, err := bcyAPI.GetAddr(req.FromAddr, params)
	if err != nil {
		return -1, -1, err
	}
	amountSelected, feeAmount, err := SelectInputsForMultisig(msgTx, utxos, 0, req, networkParams)
	if err != nil {
		return -1, -1, err
	}

	return amountSelected, feeAmount, nil
}

func SelectInputs(msgTx *wire.MsgTx, utxos gobcy.Addr, amountSelected int64, req *TxReq, networkParams *chaincfg.Params) (int64, int, error) {
	var feeAmount int
	for _, txref := range utxos.TXRefs {
		amountSelected += int64(txref.Value)
		utxoHash, err := chainhash.NewHashFromStr(txref.TXHash)
		if err != nil {
			return -1, -1, err
		}
		outPoint := wire.NewOutPoint(utxoHash, uint32(txref.TXOutputN))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		msgTx.AddTxIn(txIn)
		req.InputValues = append(req.InputValues, int64(txref.Value))

		feeAmount, err := EstimateFee(msgTx, req.FromAddr, networkParams)
		if err != nil {
			return -1, -1, err
		}
		if amountSelected >= req.Amount+int64(feeAmount) {
			return amountSelected, feeAmount, nil
		}
	}
	return -1, -1, fmt.Errorf("Insufficent funds: amountSelected %d < amountToSend %d + gassFee %d", amountSelected, req.Amount, feeAmount)
	//return -1, -1, errors.New(fmt.Sprintf("Insufficent funds: amountSelected %d < amountToSend %d + gassFee %d", amountSelected, amountToSend, feeAmount))
}

func SelectInputsForMultisig(msgTx *wire.MsgTx, utxos gobcy.Addr, amountSelected int64, req *MultisigTxReq, networkParams *chaincfg.Params) (int64, int, error) {
	var feeAmount int
	for _, txref := range utxos.TXRefs {
		amountSelected += int64(txref.Value)
		utxoHash, err := chainhash.NewHashFromStr(txref.TXHash)
		if err != nil {
			return -1, -1, err
		}
		outPoint := wire.NewOutPoint(utxoHash, uint32(txref.TXOutputN))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		msgTx.AddTxIn(txIn)
		req.InputValues = append(req.InputValues, int64(txref.Value))

		feeAmount, err := EstimateFee(msgTx, req.FromAddr, networkParams)
		if err != nil {
			return -1, -1, err
		}
		if amountSelected >= req.Amount+int64(feeAmount) {
			return amountSelected, feeAmount, nil
		}
	}
	return -1, -1, fmt.Errorf("Insufficent funds: amountSelected %d < amountToSend %d + gassFee %d", amountSelected, req.Amount, feeAmount)
	//return -1, -1, errors.New(fmt.Sprintf("Insufficent funds: amountSelected %d < amountToSend %d + gassFee %d", amountSelected, amountToSend, feeAmount))
}

func AddOutputs(msgTx *wire.MsgTx, req *TxReq /*toAddr string, fromAddr string, amountToSend int64,*/, amountSelected int64, feeAmount int, networkParams *chaincfg.Params) (*[]byte, error) {
	toAddrScript, err := GetPayToAddrScript(req.ToAddr, networkParams)
	if err != nil {
		return nil, err
	}
	msgTx.AddTxOut(wire.NewTxOut(req.Amount, *toAddrScript))

	changeAmount := amountSelected - req.Amount - int64(feeAmount)
	fromAddrScript, err := GetPayToAddrScript(req.FromAddr, networkParams)
	if err != nil {
		return nil, err
	}
	msgTx.AddTxOut(wire.NewTxOut(changeAmount, *fromAddrScript))
	return fromAddrScript, nil
}

func AddOutputsForMultisig(msgTx *wire.MsgTx, req *MultisigTxReq /*toAddr string, fromAddr string, amountToSend int64,*/, amountSelected int64, feeAmount int, networkParams *chaincfg.Params) error {
	toAddrScript, err := GetPayToAddrScript(req.ToAddr, networkParams)
	if err != nil {
		return err
	}
	msgTx.AddTxOut(wire.NewTxOut(req.Amount, *toAddrScript))

	changeAmount := amountSelected - req.Amount - int64(feeAmount)
	//fromAddrScript, err := GetLockingScript(req.FromAddr, networkParams)
	//if err != nil {
	//	return err
	//}
	msgTx.AddTxOut(wire.NewTxOut(changeAmount, req.ScriptPubKey))
	//msgTx.AddTxOut(wire.NewTxOut(changeAmount, *fromAddrScript))
	return nil
}

func AddOutputsForSegWitMultisig(msgTx *wire.MsgTx, req *MultisigTxReq /*toAddr string, fromAddr string, amountToSend int64,*/, amountSelected int64, feeAmount int, networkParams *chaincfg.Params) error {
	toAddrScript, err := GetPayToAddrScript(req.ToAddr, networkParams)
	if err != nil {
		return err
	}
	msgTx.AddTxOut(wire.NewTxOut(req.Amount, *toAddrScript))

	changeAmount := amountSelected - req.Amount - int64(feeAmount)
	//fromAddrScript, err := GetSegWitLockingScript(req.FromAddr, networkParams)
	//if err != nil {
	//	return err
	//}
	msgTx.AddTxOut(wire.NewTxOut(changeAmount, req.ScriptPubKey))
	//msgTx.AddTxOut(wire.NewTxOut(changeAmount, *fromAddrScript))
	return nil
}

func GetPayToAddrScript(addrStr string, networkParams *chaincfg.Params) (*[]byte, error) {
	addr, err := btcutil.DecodeAddress(addrStr, networkParams)
	if err != nil {
		return nil, err
	}

	addrScript, err := txscript.PayToAddrScript(addr)
	//addrScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	return &addrScript, nil
}

func GetLockingScript(addrStr string, networkParams *chaincfg.Params) (*[]byte, error) {
	addr, err := btcutil.DecodeAddress(addrStr, networkParams)
	if err != nil {
		return nil, err
	}

	lockingScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	//redeemScriptHash := addr.ScriptAddress() // this is the same as redeemScriptHash
	//builder := txscript.NewScriptBuilder()
	//builder.AddOp(txscript.OP_HASH160).AddData(redeemScriptHash).AddOp(txscript.OP_EQUAL)
	//lockingScript, err := builder.Script()
	//if err != nil {
	//	return nil, err
	//}

	return &lockingScript, nil
}

func GetSegWitLockingScript(addrStr string, networkParams *chaincfg.Params) (*[]byte, error) {
	addr, err := btcutil.DecodeAddress(addrStr, networkParams)
	if err != nil {
		return nil, err
	}

	lockingScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	//redeemScriptHash := addr.ScriptAddress() // this is the same as redeemScriptHash
	//builder := txscript.NewScriptBuilder()
	//builder.AddOp(txscript.OP_0).AddData(redeemScriptHash)
	//lockingScript, err := builder.Script()
	//if err != nil {
	//	return nil, err
	//}

	return &lockingScript, nil
}
