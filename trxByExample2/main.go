package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/cxyzhang0/wallet-go/config"
	"log"
)

var conf = config.ParseConfig()
var bcyAPI = gobcy.API{Token: conf.Blockcypher.Token, Coin: conf.Blockcypher.Coin, Chain: conf.Blockcypher.Chain}

func main() {
	var fromAddrPrivKey string
	var toAddr string
	var amountToSend int64
	var broadcast bool

	flag.StringVar(&fromAddrPrivKey, "k", "", "private key in WIF for fromAddr")
	flag.StringVar(&toAddr, "to", "", "destination addr in base58 encoding")
	flag.Int64Var(&amountToSend, "amount", 0, "amount to send in Satoshi")
	flag.BoolVar(&broadcast, "b", false, "whether to broadcast tx to blockchain")
	flag.Parse()

	networkParams, err := GetBTCNetworkParams()
	if err != nil {
		fmt.Println(err)
	}

	fromAddr, wif, err := GetAddrFromPrivKey(fromAddrPrivKey, networkParams)
	if err != nil {
		fmt.Println(err)
	}

	if signedTxReq, txHash, err := CreateTxAndHash(toAddr, fromAddr, amountToSend, wif, networkParams); err != nil {
		fmt.Println(err)
		return
	} else {
		log.Println("signedTxReq: ", signedTxReq)
		log.Println("txHash: ", txHash)
		if broadcast {
			bcyAPI.PushTX(signedTxReq)
		}
	}
}

func CreateTxAndHash(toAddr string, fromAddr string, amountToSend int64, wif *btcutil.WIF, networkParams *chaincfg.Params) (string, string, error) {
	msgTx, err := NewTx()
	if err != nil {
		return "", "", err
	}

	// coin selection for inputs
	amountSelected, gasFee, err := SelectCoins(msgTx, toAddr, fromAddr, amountToSend, networkParams)
	if err != nil {
		return "", "", err
	}

	// outputs
	fromAddrScript, err := AddOutputs(msgTx, toAddr, fromAddr, amountToSend, amountSelected, gasFee, networkParams)
	if err != nil {
		return "", "", err
	}

	// signature
	err = SignTx(msgTx, wif, fromAddrScript)
	if err != nil {
		return "", "", err
	}

	// verify
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

func SelectCoins(msgTx *wire.MsgTx /*wif *btcutil.WIF,*/, toAddr string, fromAddr string, amountToSend int64, networkParams *chaincfg.Params) (int64, int, error) {
	// check the balance of fromAddr
	addr, err := bcyAPI.GetAddrBal(fromAddr, nil)
	if err != nil {
		return -1, -1, err
	}

	// initial sufficient balance check
	if int64(addr.Balance) <= amountToSend {
		return -1, -1, errors.New(fmt.Sprintf("Insufficent funds (initial check): balance %d <= amountToSend %d", addr.Balance, amountToSend))
	}

	// select coins
	var params = make(map[string]string, 0)
	params["unspentOnly"] = "true"
	params["includeScript"] = "true"
	utxos, err := bcyAPI.GetAddr(fromAddr, params)
	if err != nil {
		return -1, -1, err
	}
	amountSelected, feeAmount, err := SelectInputs(msgTx, utxos, 0, amountToSend)
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

func AddOutputs(msgTx *wire.MsgTx, toAddr string, fromAddr string, amountToSend int64, amountSelected int64, feeAmount int, networkParams *chaincfg.Params) (*[]byte, error) {
	toAddrScript, err := GetPayToAddrScript(toAddr, networkParams)
	if err != nil {
		return nil, err
	}
	msgTx.AddTxOut(wire.NewTxOut(amountToSend, *toAddrScript))

	changeAmount := amountSelected - amountToSend - int64(feeAmount)
	fromAddrScript, err := GetPayToAddrScript(fromAddr, networkParams)
	if err != nil {
		return nil, err
	}
	msgTx.AddTxOut(wire.NewTxOut(changeAmount, *fromAddrScript))
	return fromAddrScript, nil
}

func SignTx(msgTx *wire.MsgTx, wif *btcutil.WIF, fromAddrScript *[]byte) error {
	for i, txIn := range msgTx.TxIn {
		sigScript, err := txscript.SignatureScript(
			msgTx,
			i,
			*fromAddrScript,
			txscript.SigHashAll,
			wif.PrivKey,
			false)
		if err != nil {
			return err
		}
		txIn.SignatureScript = sigScript
	}
	return nil
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

func GetAddrFromPrivKey(fromAddrPrivKey string, networkParams *chaincfg.Params) (string, *btcutil.WIF, error) {
	wif, err := btcutil.DecodeWIF(fromAddrPrivKey)
	if err != nil {
		return "", nil, err
	}

	fromAddrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), networkParams)
	if err != nil {
		return "", wif, err
	}

	fromAddr := fromAddrPubKey.EncodeAddress()
	return fromAddr, wif, nil
	/**
	Note: if fromAddrPrivKey is in hex instead of WIF,
	privateKeyBytes, err := hex.DecodeString(fromAddrPrivKey)
		if err != nil {
			return
		}
	privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	*/
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
