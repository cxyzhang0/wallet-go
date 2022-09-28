package btctran_gcp

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	tsx "github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	gcpsdk "github.com/cxyzhang0/wallet-go/gcp/sdk"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func SignTx(msgTx *wire.MsgTx, req *TxReq, fromAddrScript *[]byte, sdk *gcpsdk.SDK) error {
	for i, txIn := range msgTx.TxIn {
		sigScript, err := _SignatureScript(
			msgTx,
			i,
			*fromAddrScript,
			tsx.SigHashAll,
			req,
			false,
			sdk)
		if err != nil {
			return err
		}
		txIn.SignatureScript = sigScript
	}
	return nil
}

func _SignatureScript(tx *wire.MsgTx, idx int, subscript []byte, hashType tsx.SigHashType, req *TxReq, compress bool, sdk *gcpsdk.SDK) ([]byte, error) {
	sig, err := _RawTxInSignature(tx, idx, subscript, hashType, req, sdk)
	if err != nil {
		return nil, err
	}

	pk := req.FromAddrPubKey.PubKey()
	//pk := (*btcec.PublicKey)(&privKey.PublicKey)
	var pkData []byte
	if compress {
		pkData = pk.SerializeCompressed()
	} else {
		pkData = pk.SerializeUncompressed()
	}

	return tsx.NewScriptBuilder().AddData(sig).AddData(pkData).Script()
}

func _RawTxInSignature(tx *wire.MsgTx, idx int, subScript []byte,
	hashType tsx.SigHashType, req *TxReq, sdk *gcpsdk.SDK) ([]byte, error) {

	hash, err := tsx.CalcSignatureHash(subScript, hashType, tx, idx)
	if err != nil {
		return nil, err
	}

	//signature, err := key.Sign(hash)
	sig, err := sdk.Sign(req.From, hash)
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx input: %s", err)
	}

	signature, err := btcec.ParseSignature(sig, secp256k1.S256())
	if err != nil {
		return nil, err
	}

	return append(signature.Serialize(), byte(hashType)), nil
}
