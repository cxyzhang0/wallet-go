package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func main() {
	// privKey and dest are from walletgenerator: Private Key (Wallet Import Format), Public Address
	// Why is this tx paying to the same wallet?
	// privKey: "92eU8tJWm4nYj1NuUnDg5KQT73BvEiSrjnBMEpCkeZ4A5eu97aF", address: "n4VVN4CZhm2rdNkh7mXGzWJLaRePVitWhT"  - it has balance
	// privKey: "91bcGqKhkDW699hMGF5s9wf1adqV2pMFbp7NTpQAYJgcK2jPLaY", address: "mwdumvCrxnMAXyME4NaZEA5YWpDttpXPuT"
	rawTx, txHash, err := CreateTxAndHash("92eU8tJWm4nYj1NuUnDg5KQT73BvEiSrjnBMEpCkeZ4A5eu97aF", "mwdumvCrxnMAXyME4NaZEA5YWpDttpXPuT", 6000)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Raw Transaction: ", rawTx)
	fmt.Println("Transaction Hash: ", txHash)
}

func NewTx() (*wire.MsgTx, error) {
	return wire.NewMsgTx(wire.TxVersion), nil
}

// coin selection
func GetUTXO(address string) (string, int64, string, error) {
	var prevTxId string = "ce4ae4fc38b962452d5e8ff0889da4440ecb53649fc02bfed4cb3bc1e40a1111" // this is from faucet when requesting to fill Public Address n4VVN4CZhm2rdNkh7mXGzWJLaRePVitWhT
	//var prevTxId string = "9f617169bfa12ca485e10aad2c7d620f7c3fc99f31c9c3bdff822a4eb32969f7" // this is another faucet fill to the same above addr
	var balance int64 = 62000
	var pubKeyScript string = "76a914fc03ff06d1d5b21733b3a507cfa7c1cdfa74b80a88ac" // this is from BlockCypher: transactions for address n4VVN4CZhm2rdNkh7mXGzWJLaRePVitWhT
	return prevTxId, balance, pubKeyScript, nil
}

func CreateTx(privKey string, dest string, amount int64) (string, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)

	txid, balance, pubKeyScript, err := GetUTXO(addrPubKey.EncodeAddress())
	if err != nil {
		return "", err
	}

	if balance < amount {
		return "", fmt.Errorf("Insufficient Balance in wallet")
	}

	destAddr, err := btcutil.DecodeAddress(dest, &chaincfg.TestNet3Params) // dest = destAddr.AddressScript() or String()
	if err != nil {
		return "", err
	}
	destAddrScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		return "", err
	}

	redeemTx, err := NewTx() // empty tx
	if err != nil {
		return "", err
	}

	utxoHash, err := chainhash.NewHashFromStr(txid) // txid = utxoHash.String()
	if err != nil {
		return "", err
	}

	outPoint := wire.NewOutPoint(utxoHash, 0) // need coin selection
	txIn := wire.NewTxIn(outPoint, nil, nil)  // Q: should signatureScript be pubKeyScript? A: nil seems to work as is
	redeemTx.AddTxIn(txIn)

	txOut := wire.NewTxOut(amount, destAddrScript)
	//txOut := wire.NewTxOut(amount, destAddr.ScriptAddress())
	redeemTx.AddTxOut(txOut)

	finalRawTx, err := SignTx(privKey, pubKeyScript, redeemTx)
	if err != nil {
		return "", err
	}

	return finalRawTx, nil
}

func CreateTxAndHash(privKey string, dest string, amount int64) (string, string, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", "", err
	}

	addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.TestNet3Params)
	//addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)

	txid, balance, pubKeyScript, err := GetUTXO(addrPubKey.EncodeAddress())
	if err != nil {
		return "", "", err
	}

	if balance < amount {
		return "", "", fmt.Errorf("Insufficient Balance in wallet")
	}

	destAddr, err := btcutil.DecodeAddress(dest, &chaincfg.TestNet3Params) // dest = destAddr.AddressScript() or String()
	if err != nil {
		return "", "", err
	}
	destAddrScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		return "", "", err
	}

	redeemTx, err := NewTx() // empty tx
	if err != nil {
		return "", "", err
	}

	utxoHash, err := chainhash.NewHashFromStr(txid) // txid = utxoHash.String()
	if err != nil {
		return "", "", err
	}

	outPoint := wire.NewOutPoint(utxoHash, 0) // need coin selection
	txIn := wire.NewTxIn(outPoint, nil, nil)  // Q: should signatureScript be pubKeyScript? A: nil seems to work as is
	redeemTx.AddTxIn(txIn)

	txOut := wire.NewTxOut(amount, destAddrScript)
	//txOut := wire.NewTxOut(amount, destAddr.ScriptAddress())
	redeemTx.AddTxOut(txOut)

	finalRawTx, err := SignTx(privKey, pubKeyScript, redeemTx)
	if err != nil {
		return "", "", err
	}

	return finalRawTx, redeemTx.TxHash().String(), nil
}

func SignTx(privKey string, pubKeyScript string, redeemTx *wire.MsgTx) (string, error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	sourcePubKeyScript, err := hex.DecodeString(pubKeyScript)
	if err != nil {
		return "", err
	}

	signature, err := txscript.SignatureScript(redeemTx, 0, sourcePubKeyScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return "", err
	}
	redeemTx.TxIn[0].SignatureScript = signature

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)

	btcutil.Hash160(signedTx.Bytes())
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	return hexSignedTx, nil
}
