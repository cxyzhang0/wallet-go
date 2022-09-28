package btctran_gcp

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
	"github.com/cxyzhang0/wallet-go/btctran_gcp/test/config"
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
)

var conf = config.ParseConfig()
var bcyAPI = gobcy.API{Token: conf.Blockcypher.Token, Coin: conf.Blockcypher.Coin, Chain: conf.Blockcypher.Chain}

func GetAddressPubKey(keyLabel gcpsdk.KeyLabel, sdk *gcpsdk.SDK, networkParams *chaincfg.Params) (*btcutil.AddressPubKey, string, error) { // address pub key, address, error
	pubkey, err := sdk.GetECDSAPublicKeyForSecp256k1(keyLabel)
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

func NewTx() (*wire.MsgTx, error) {
	return wire.NewMsgTx(wire.TxVersion), nil
}

func EstimateFee(msgTx *wire.MsgTx) int {
	inputCount := len(msgTx.TxIn)
	outputCount := len(msgTx.TxOut)
	if outputCount < 2 {
		outputCount = 2
	}
	txSize := inputCount*180 + outputCount*34 + 10
	return conf.Blockcypher.FeeRate * txSize
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

func SelectCoins(msgTx *wire.MsgTx, req *TxReq /*fromAddr string, amountToSend int64*/) (int64, int, error) {
	// check the balance of fromAddr
	addr, err := bcyAPI.GetAddrBal(req.FromAddr, nil)
	if err != nil {
		return -1, -1, err
	}

	// initial sufficient balance check
	if int64(addr.Balance) <= req.Amount {
		return -1, -1, errors.New(fmt.Sprintf("Insufficent funds (initial check): balance %d <= amountToSend %d", addr.Balance, req.Amount))
	}

	// select coins
	var params = make(map[string]string, 0)
	params["unspentOnly"] = "true"
	params["includeScript"] = "true"
	utxos, err := bcyAPI.GetAddr(req.FromAddr, params)
	if err != nil {
		return -1, -1, err
	}
	amountSelected, feeAmount, err := SelectInputs(msgTx, utxos, 0, req.Amount)
	if err != nil {
		return -1, -1, err
	}

	return amountSelected, feeAmount, nil
}

func SelectInputs(msgTx *wire.MsgTx, utxos gobcy.Addr, amountSelected int64, amountToSend int64) (int64, int, error) {
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
		feeAmount = EstimateFee(msgTx)
		if amountSelected >= amountToSend+int64(feeAmount) {
			return amountSelected, feeAmount, nil
		}
	}
	return -1, -1, errors.New(fmt.Sprintf("Insufficent funds: amountSelected %d < amountToSend %d + gassFee %d", amountSelected, amountToSend, feeAmount))
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

func GetPayToAddrScript(addrStr string, networkParams *chaincfg.Params) (*[]byte, error) {
	addr, err := btcutil.DecodeAddress(addrStr, networkParams)
	if err != nil {
		return nil, err
	}

	addrScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	return &addrScript, nil
}
