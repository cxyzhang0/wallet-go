package algotran_pkcs11

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/base32"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
	kmssdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"github.com/pkg/errors"
)

// txidPrefix is prepended to a transaction when computing its txid
var txidPrefix = []byte("TX")

func GetAddressPubKey(keyLabel kmssdk.KeyLabel, sdk *kmssdk.SDK) (*ed25519.PublicKey, *types.Address, string, error) {
	pubKey, err := sdk.GetEdDSAPublicKey(keyLabel)
	if err != nil {
		return nil, nil, "", err
	}
	var a types.Address
	copy(a[:], *pubKey)

	return pubKey, &a, a.String(), nil
}

// SignTransaction signs a transaction by hsm, and returns the
// bytes of a signed transaction ready to be broadcasted to the network
// If the signerKeyLabel's corresponding address is different than the txn sender's, the signerKeyLabel's
// corresponding address will be assigned as AuthAddr
func SignTransaction(req TxReq, tx types.Transaction, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	s, txid, err := rawSignTransaction(req.From, tx, sdk)
	if err != nil {
		return
	}
	// Construct the SignedTxn
	stx := types.SignedTxn{
		Sig: s,
		Txn: tx,
	}

	if stx.Txn.Sender != *req.FromAddress {
		stx.AuthAddr = *req.FromAddress
	}

	// Encode the SignedTxn
	stxBytes = msgpack.Encode(stx)
	return
}

// rawSignTransaction signs the msgpack-encoded tx (with prepended "TX" prefix), and returns the sig and txid
func rawSignTransaction(keyLabel kmssdk.KeyLabel, tx types.Transaction, sdk *kmssdk.SDK) (s types.Signature, txid string, err error) {
	toBeSigned := rawTransactionBytesToSign(tx)

	// Sign the encoded transaction
	signature, err := sdk.Sign(keyLabel, toBeSigned)
	if err != nil {
		return
	}

	//TODO: verify fails but the signed tx is successfully submitted. why?
	//verified := crypto.VerifyBytes(*req.FromPubKey, toBeSigned, signature)
	//if !verified {
	//	err = errors.New("signature verification failed")
	//}

	// Copy the resulting signature into a Signature, and check that it's
	// the expected length
	n := copy(s[:], signature)
	if n != len(s) {
		err = errors.New("pkcs11 sdk ed25519 returned an invalid ed25519 signature")
		return
	}
	// Populate txID
	txid = txIDFromRawTxnBytesToSign(toBeSigned)
	return
}

// rawTransactionBytesToSign returns the byte form of the tx that we actually sign
// and compute txID from.
func rawTransactionBytesToSign(tx types.Transaction) []byte {
	// Encode the transaction as msgpack
	encodedTx := msgpack.Encode(tx)

	// Prepend the hashable prefix
	msgParts := [][]byte{txidPrefix, encodedTx}
	return bytes.Join(msgParts, nil)
}

// txID computes a transaction id base32 string from raw transaction bytes
func txIDFromRawTxnBytesToSign(toBeSigned []byte) (txid string) {
	txidBytes := sha512.Sum512_256(toBeSigned)
	txid = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(txidBytes[:])
	return
}
