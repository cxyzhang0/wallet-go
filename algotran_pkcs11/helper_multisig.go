package algotran_pkcs11

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
	kmssdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
	"github.com/pkg/errors"
)

func GetMultisigAddress(keyLabels []kmssdk.KeyLabel, m uint8, sdk *kmssdk.SDK) (*crypto.MultisigAccount, *types.Address, []ed25519.PublicKey, string, error) {
	n := len(keyLabels)
	if n < int(m) {
		return nil, nil, nil, "", errors.Errorf("requre n %d >= m %d", n, m)
	}
	addresses := make([]types.Address, n)
	pubKeys := make([]ed25519.PublicKey, n)
	for i, keyLabel := range keyLabels {
		pubKey, address, _, err := GetAddressPubKey(keyLabel, sdk)
		if err != nil {
			return nil, nil, nil, "", errors.Errorf("GetAddressPubKey failed for i %d keyLabel %+v", i, keyLabel)
		}
		addresses[i] = *address
		pubKeys[i] = *pubKey
	}

	ma, err := crypto.MultisigAccountWithParams(1, m, addresses)
	if err != nil {
		return nil, nil, nil, "", errors.WithMessage(err, "MultisigAccountWithParams failed")
	}
	a, err := ma.Address()
	if err != nil {
		return nil, nil, nil, "", errors.WithMessage(err, "ma.Address() failed")
	}
	return &ma, &a, pubKeys, a.String(), nil
}

// SignMultisigTransaction signs the given transaction, and multisig preimage, with the
// private key, returning the bytes of a signed transaction with the multisig field
// partially populated, ready to be passed to other multisig signers to sign or broadcast.
func SignMultisigTransaction(keyLabel kmssdk.KeyLabel, pubKey ed25519.PublicKey, ma crypto.MultisigAccount, tx types.Transaction, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	err = ma.Validate()
	if err != nil {
		return
	}

	// this signer signs a transaction and sets txid from the closure
	customSigner := func() (rawSig types.Signature, err error) {
		rawSig, txid, err = rawSignTransaction(keyLabel, tx, sdk)
		return rawSig, err
	}

	sig, _, err := multisigSingle(pubKey, ma, customSigner)
	if err != nil {
		return
	}

	// Encode the signedTxn
	stx := types.SignedTxn{
		Msig: sig,
		Txn:  tx,
	}

	maAddress, err := ma.Address()
	if err != nil {
		return
	}

	if stx.Txn.Sender != maAddress {
		stx.AuthAddr = maAddress
	}

	stxBytes = msgpack.Encode(stx)
	return
}

// Service function to make a single signature in Multisig
func multisigSingle(myPublicKey ed25519.PublicKey, ma crypto.MultisigAccount, customSigner signer) (msig types.MultisigSig, myIndex int, err error) {
	// check that sk.pk exists in the list of public keys in MultisigAccount ma
	myIndex = len(ma.Pks)
	//myPublicKey := req.PubKeys[i] //sk.Public().(ed25519.PublicKey)
	for i := 0; i < len(ma.Pks); i++ {
		if bytes.Equal(myPublicKey, ma.Pks[i]) {
			myIndex = i
		}
	}
	if myIndex == len(ma.Pks) {
		err = errors.New("pub key mismatch")
		return
	}

	// now, create the signed transaction
	msig.Version = ma.Version
	msig.Threshold = ma.Threshold
	msig.Subsigs = make([]types.MultisigSubsig, len(ma.Pks))
	for i := 0; i < len(ma.Pks); i++ {
		c := make([]byte, len(ma.Pks[i]))
		copy(c, ma.Pks[i])
		msig.Subsigs[i].Key = c
	}
	rawSig, err := customSigner()
	if err != nil {
		return
	}
	msig.Subsigs[myIndex].Sig = rawSig
	return
}

type signer func() (signature types.Signature, err error)

// AppendMultisigTransaction appends the signature corresponding to the given private key,
// returning an encoded signed multisig transaction including the signature.
// While we could compute the multisig preimage from the multisig blob, we ask the caller
// to pass it back in, to explicitly check that they know who they are signing as.
func AppendMultisigTransaction(keyLabel kmssdk.KeyLabel, pubKey ed25519.PublicKey, ma crypto.MultisigAccount, preStxBytes []byte, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	preStx := types.SignedTxn{}
	err = msgpack.Decode(preStxBytes, &preStx)
	if err != nil {
		return
	}
	_, partStxBytes, err := SignMultisigTransaction(keyLabel, pubKey, ma, preStx.Txn, sdk)
	if err != nil {
		return
	}
	txid, stxBytes, err = crypto.MergeMultisigTransactions(partStxBytes, preStxBytes)
	return
}

// merge all m of n signatures
func makeMultisig(req MultisigTxReq, ma crypto.MultisigAccount, tx types.Transaction, sdk *kmssdk.SDK) (txid string, stxBytes []byte, err error) {
	txid, stxBytes, err = SignMultisigTransaction(req.From[0], req.PubKeys[0], ma, tx, sdk)
	if err != nil {
		return "", nil, err
	}
	fmt.Printf("partially signed multsig transaction txid %d: %s", 0, txid)

	for i := 1; i < req.M; i++ {
		txid, stxBytes, err = AppendMultisigTransaction(req.From[i], req.PubKeys[i], ma, stxBytes, sdk)
		if err != nil {
			return "", nil, err
		}
		fmt.Printf("partially signed multsig transaction txid %d: %s", i, txid)
	}

	return
}
