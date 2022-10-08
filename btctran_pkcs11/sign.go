package btctran_pkcs11

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	tsx "github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	pkcs11sdk "github.com/cxyzhang0/wallet-go/pkcs11/sdk"
)

const (
	sigHashMask = 0x1f
)

func SignTx(msgTx *wire.MsgTx, req *TxReq, fromAddrScript *[]byte, sdk *pkcs11sdk.SDK) error {
	if !tsx.IsPayToWitnessPubKeyHash(*fromAddrScript) {
		for i, txIn := range msgTx.TxIn {
			sigScript, err := SignatureScript(
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
	} else {
		txSignHashes := tsx.NewTxSigHashes(msgTx)
		for i, txIn := range msgTx.TxIn {
			witness, err := WitnessSignature(
				msgTx,
				txSignHashes,
				i,
				req.InputValues[i],
				*fromAddrScript,
				tsx.SigHashAll,
				req,
				true,
				sdk)
			if err != nil {
				return err
			}
			txIn.Witness = witness
		}
	}
	return nil
}

func SignatureScript(tx *wire.MsgTx, idx int, subscript []byte, hashType tsx.SigHashType, req *TxReq, compress bool, sdk *pkcs11sdk.SDK) ([]byte, error) {
	sig, err := RawTxInSignature(tx, idx, subscript, hashType, req, sdk)
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

func RawTxInSignature(tx *wire.MsgTx, idx int, subScript []byte,
	hashType tsx.SigHashType, req *TxReq, sdk *pkcs11sdk.SDK) ([]byte, error) {

	hash, err := tsx.CalcSignatureHash(subScript, hashType, tx, idx)
	if err != nil {
		return nil, err
	}

	sig, err := sdk.GetChainSignature(req.From, hash)
	//sig, err := sdk.Sign(req.From, hash)
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx input: %s", err)
	}

	signature, err := btcec.ParseDERSignature(sig, btcec.S256())
	//signature, err := btcec.ParseSignature(sig, btcec.S256()) // this one works too
	//signature, err := btcec.ParseSignature(sig, secp256k1.S256()) // secp256k1.S256() is from go-ethereum; use btcec instead.
	if err != nil {
		return nil, err
	}

	return append(signature.Serialize(), byte(hashType)), nil
}

// WitnessSignature creates an input witness stack for tx to spend BTC sent
// from a previous output to the owner of privKey using the p2wkh script
// template. The passed transaction must contain all the inputs and outputs as
// dictated by the passed hashType. The signature generated observes the new
// transaction digest algorithm defined within BIP0143.
func WitnessSignature(tx *wire.MsgTx, sigHashes *tsx.TxSigHashes, idx int, amt int64,
	subscript []byte, hashType tsx.SigHashType, req *TxReq,
	compress bool, sdk *pkcs11sdk.SDK) (wire.TxWitness, error) {

	sig, err := RawTxInWitnessSignature(tx, sigHashes, idx, amt, subscript,
		hashType, req, sdk)
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

	// A witness script is actually a stack, so we return an array of byte
	// slices here, rather than a single byte slice.
	return wire.TxWitness{sig, pkData}, nil
}

// RawTxInWitnessSignature returns the serialized ECDA signature for the input
// idx of the given transaction, with the hashType appended to it. This
// function is identical to RawTxInSignature, however the signature generated
// signs a new sighash digest defined in BIP0143.
func RawTxInWitnessSignature(tx *wire.MsgTx, sigHashes *tsx.TxSigHashes, idx int,
	amt int64, subScript []byte, hashType tsx.SigHashType,
	req *TxReq, sdk *pkcs11sdk.SDK) ([]byte, error) {

	parsedScript, err := parseScript(subScript)
	if err != nil {
		return nil, fmt.Errorf("cannot parse output script: %v", err)
	}

	hash, err := calcWitnessSignatureHash(parsedScript, sigHashes, hashType, tx,
		idx, amt)
	if err != nil {
		return nil, err
	}

	sig, err := sdk.GetChainSignature(req.From, hash)
	//sig, err := sdk.Sign(req.From, hash)
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx input: %s", err)
	}

	signature, err := btcec.ParseDERSignature(sig, btcec.S256())
	//signature, err := btcec.ParseSignature(sig, btcec.S256()) // this one works too
	if err != nil {
		return nil, err
	}

	return append(signature.Serialize(), byte(hashType)), nil
}

// calcWitnessSignatureHash computes the sighash digest of a transaction's
// segwit input using the new, optimized digest calculation algorithm defined
// in BIP0143: https://github.com/bitcoin/bips/blob/master/bip-0143.mediawiki.
// This function makes use of pre-calculated sighash fragments stored within
// the passed HashCache to eliminate duplicate hashing computations when
// calculating the final digest, reducing the complexity from O(N^2) to O(N).
// Additionally, signatures now cover the input value of the referenced unspent
// output. This allows offline, or hardware wallets to compute the exact amount
// being spent, in addition to the final transaction fee. In the case the
// wallet if fed an invalid input amount, the real sighash will differ causing
// the produced signature to be invalid.
func calcWitnessSignatureHash(subScript []parsedOpcode, sigHashes *tsx.TxSigHashes,
	hashType tsx.SigHashType, tx *wire.MsgTx, idx int, amt int64) ([]byte, error) {

	// As a sanity check, ensure the passed input index for the transaction
	// is valid.
	if idx > len(tx.TxIn)-1 {
		return nil, fmt.Errorf("idx %d but %d txins", idx, len(tx.TxIn))
	}

	// We'll utilize this buffer throughout to incrementally calculate
	// the signature hash for this transaction.
	var sigHash bytes.Buffer

	// First write out, then encode the transaction's version number.
	var bVersion [4]byte
	binary.LittleEndian.PutUint32(bVersion[:], uint32(tx.Version))
	sigHash.Write(bVersion[:])

	// Next write out the possibly pre-calculated hashes for the sequence
	// numbers of all inputs, and the hashes of the previous outs for all
	// outputs.
	var zeroHash chainhash.Hash

	// If anyone can pay isn't active, then we can use the cached
	// hashPrevOuts, otherwise we just write zeroes for the prev outs.
	if hashType&tsx.SigHashAnyOneCanPay == 0 {
		sigHash.Write(sigHashes.HashPrevOuts[:])
	} else {
		sigHash.Write(zeroHash[:])
	}

	// If the sighash isn't anyone can pay, single, or none, the use the
	// cached hash sequences, otherwise write all zeroes for the
	// hashSequence.
	if hashType&tsx.SigHashAnyOneCanPay == 0 &&
		hashType&sigHashMask != tsx.SigHashSingle &&
		hashType&sigHashMask != tsx.SigHashNone {
		sigHash.Write(sigHashes.HashSequence[:])
	} else {
		sigHash.Write(zeroHash[:])
	}

	txIn := tx.TxIn[idx]

	// Next, write the outpoint being spent.
	sigHash.Write(txIn.PreviousOutPoint.Hash[:])
	var bIndex [4]byte
	binary.LittleEndian.PutUint32(bIndex[:], txIn.PreviousOutPoint.Index)
	sigHash.Write(bIndex[:])

	if isWitnessPubKeyHash(subScript) {
		// The script code for a p2wkh is a length prefix varint for
		// the next 25 bytes, followed by a re-creation of the original
		// p2pkh pk script.
		sigHash.Write([]byte{0x19})
		sigHash.Write([]byte{tsx.OP_DUP})
		sigHash.Write([]byte{tsx.OP_HASH160})
		sigHash.Write([]byte{tsx.OP_DATA_20})
		sigHash.Write(subScript[1].data)
		sigHash.Write([]byte{tsx.OP_EQUALVERIFY})
		sigHash.Write([]byte{tsx.OP_CHECKSIG})
	} else {
		// For p2wsh outputs, and future outputs, the script code is
		// the original script, with all code separators removed,
		// serialized with a var int length prefix.
		rawScript, _ := unparseScript(subScript)
		wire.WriteVarBytes(&sigHash, 0, rawScript)
	}

	// Next, add the input amount, and sequence number of the input being
	// signed.
	var bAmount [8]byte
	binary.LittleEndian.PutUint64(bAmount[:], uint64(amt))
	sigHash.Write(bAmount[:])
	var bSequence [4]byte
	binary.LittleEndian.PutUint32(bSequence[:], txIn.Sequence)
	sigHash.Write(bSequence[:])

	// If the current signature mode isn't single, or none, then we can
	// re-use the pre-generated hashoutputs sighash fragment. Otherwise,
	// we'll serialize and add only the target output index to the signature
	// pre-image.
	if hashType&tsx.SigHashSingle != tsx.SigHashSingle &&
		hashType&tsx.SigHashNone != tsx.SigHashNone {
		sigHash.Write(sigHashes.HashOutputs[:])
	} else if hashType&sigHashMask == tsx.SigHashSingle && idx < len(tx.TxOut) {
		var b bytes.Buffer
		wire.WriteTxOut(&b, 0, 0, tx.TxOut[idx])
		sigHash.Write(chainhash.DoubleHashB(b.Bytes()))
	} else {
		sigHash.Write(zeroHash[:])
	}

	// Finally, write out the transaction's locktime, and the sig hash
	// type.
	var bLockTime [4]byte
	binary.LittleEndian.PutUint32(bLockTime[:], tx.LockTime)
	sigHash.Write(bLockTime[:])
	var bHashType [4]byte
	binary.LittleEndian.PutUint32(bHashType[:], uint32(hashType))
	sigHash.Write(bHashType[:])

	return chainhash.DoubleHashB(sigHash.Bytes()), nil
}

// calcSignatureHash will, given a script and hash type for the current script
// engine instance, calculate the signature hash to be used for signing and
// verification.
func calcSignatureHash(script []parsedOpcode, hashType tsx.SigHashType, tx *wire.MsgTx, idx int) []byte {
	// The SigHashSingle signature type signs only the corresponding input
	// and output (the output with the same index number as the input).
	//
	// Since transactions can have more inputs than outputs, this means it
	// is improper to use SigHashSingle on input indices that don't have a
	// corresponding output.
	//
	// A bug in the original Satoshi client implementation means specifying
	// an index that is out of range results in a signature hash of 1 (as a
	// uint256 little endian).  The original intent appeared to be to
	// indicate failure, but unfortunately, it was never checked and thus is
	// treated as the actual signature hash.  This buggy behavior is now
	// part of the consensus and a hard fork would be required to fix it.
	//
	// Due to this, care must be taken by software that creates transactions
	// which make use of SigHashSingle because it can lead to an extremely
	// dangerous situation where the invalid inputs will end up signing a
	// hash of 1.  This in turn presents an opportunity for attackers to
	// cleverly construct transactions which can steal those coins provided
	// they can reuse signatures.
	if hashType&sigHashMask == tsx.SigHashSingle && idx >= len(tx.TxOut) {
		var hash chainhash.Hash
		hash[0] = 0x01
		return hash[:]
	}

	// Remove all instances of OP_CODESEPARATOR from the script.
	script = removeOpcode(script, OP_CODESEPARATOR)

	// Make a shallow copy of the transaction, zeroing out the script for
	// all inputs that are not currently being processed.
	txCopy := shallowCopyTx(tx)
	for i := range txCopy.TxIn {
		if i == idx {
			// UnparseScript cannot fail here because removeOpcode
			// above only returns a valid script.
			sigScript, _ := unparseScript(script)
			txCopy.TxIn[idx].SignatureScript = sigScript
		} else {
			txCopy.TxIn[i].SignatureScript = nil
		}
	}

	switch hashType & sigHashMask {
	case tsx.SigHashNone:
		txCopy.TxOut = txCopy.TxOut[0:0] // Empty slice.
		for i := range txCopy.TxIn {
			if i != idx {
				txCopy.TxIn[i].Sequence = 0
			}
		}

	case tsx.SigHashSingle:
		// Resize output array to up to and including requested index.
		txCopy.TxOut = txCopy.TxOut[:idx+1]

		// All but current output get zeroed out.
		for i := 0; i < idx; i++ {
			txCopy.TxOut[i].Value = -1
			txCopy.TxOut[i].PkScript = nil
		}

		// Sequence on all other inputs is 0, too.
		for i := range txCopy.TxIn {
			if i != idx {
				txCopy.TxIn[i].Sequence = 0
			}
		}

	default:
		// Consensus treats undefined hashtypes like normal SigHashAll
		// for purposes of hash generation.
		fallthrough
	case tsx.SigHashOld:
		fallthrough
	case tsx.SigHashAll:
		// Nothing special here.
	}
	if hashType&tsx.SigHashAnyOneCanPay != 0 {
		txCopy.TxIn = txCopy.TxIn[idx : idx+1]
	}

	// The final hash is the double sha256 of both the serialized modified
	// transaction and the hash type (encoded as a 4-byte little-endian
	// value) appended.
	wbuf := bytes.NewBuffer(make([]byte, 0, txCopy.SerializeSizeStripped()+4))
	txCopy.SerializeNoWitness(wbuf)
	binary.Write(wbuf, binary.LittleEndian, hashType)
	return chainhash.DoubleHashB(wbuf.Bytes())
}

// isWitnessPubKeyHash returns true if the passed script is a
// pay-to-witness-pubkey-hash, and false otherwise.
func isWitnessPubKeyHash(pops []parsedOpcode) bool {
	return len(pops) == 2 &&
		pops[0].opcode.value == tsx.OP_0 &&
		pops[1].opcode.value == tsx.OP_DATA_20
}

// parseScript preparses the script in bytes into a list of parsedOpcodes while
// applying a number of sanity checks.
func parseScript(script []byte) ([]parsedOpcode, error) {
	return parseScriptTemplate(script, &opcodeArray)
}

// unparseScript reversed the action of parseScript and returns the
// parsedOpcodes as a list of bytes
func unparseScript(pops []parsedOpcode) ([]byte, error) {
	script := make([]byte, 0, len(pops))
	for _, pop := range pops {
		b, err := pop.bytes()
		if err != nil {
			return nil, err
		}
		script = append(script, b...)
	}
	return script, nil
}

// parseScriptTemplate is the same as parseScript but allows the passing of the
// template list for testing purposes.  When there are parse errors, it returns
// the list of parsed opcodes up to the point of failure along with the error.
func parseScriptTemplate(script []byte, opcodes *[256]opcode) ([]parsedOpcode, error) {
	retScript := make([]parsedOpcode, 0, len(script))
	var err error
	for i := 0; i < len(script); {
		instr := script[i]
		op := &opcodes[instr]
		pop := parsedOpcode{opcode: op}
		i, err = pop.checkParseableInScript(script, i)
		if err != nil {
			return retScript, err
		}

		retScript = append(retScript, pop)
	}

	return retScript, nil
}

// removeOpcodeByData will return the script minus any opcodes that would push
// the passed data to the stack.
func removeOpcodeByData(pkscript []parsedOpcode, data []byte) []parsedOpcode {
	retScript := make([]parsedOpcode, 0, len(pkscript))
	for _, pop := range pkscript {
		if !canonicalPush(pop) || !bytes.Contains(pop.data, data) {
			retScript = append(retScript, pop)
		}
	}
	return retScript

}

// removeOpcode will remove any opcode matching ``opcode'' from the opcode
// stream in pkscript
func removeOpcode(pkscript []parsedOpcode, opcode byte) []parsedOpcode {
	retScript := make([]parsedOpcode, 0, len(pkscript))
	for _, pop := range pkscript {
		if pop.opcode.value != opcode {
			retScript = append(retScript, pop)
		}
	}
	return retScript
}

// canonicalPush returns true if the object is either not a push instruction
// or the push instruction contained wherein is matches the canonical form
// or using the smallest instruction to do the job. False otherwise.
func canonicalPush(pop parsedOpcode) bool {
	opcode := pop.opcode.value
	data := pop.data
	dataLen := len(pop.data)
	if opcode > OP_16 {
		return true
	}

	if opcode < OP_PUSHDATA1 && opcode > OP_0 && (dataLen == 1 && data[0] <= 16) {
		return false
	}
	if opcode == OP_PUSHDATA1 && dataLen < OP_PUSHDATA1 {
		return false
	}
	if opcode == OP_PUSHDATA2 && dataLen <= 0xff {
		return false
	}
	if opcode == OP_PUSHDATA4 && dataLen <= 0xffff {
		return false
	}
	return true
}

// shallowCopyTx creates a shallow copy of the transaction for use when
// calculating the signature hash.  It is used over the Copy method on the
// transaction itself since that is a deep copy and therefore does more work and
// allocates much more space than needed.
func shallowCopyTx(tx *wire.MsgTx) wire.MsgTx {
	// As an additional memory optimization, use contiguous backing arrays
	// for the copied inputs and outputs and point the final slice of
	// pointers into the contiguous arrays.  This avoids a lot of small
	// allocations.
	txCopy := wire.MsgTx{
		Version:  tx.Version,
		TxIn:     make([]*wire.TxIn, len(tx.TxIn)),
		TxOut:    make([]*wire.TxOut, len(tx.TxOut)),
		LockTime: tx.LockTime,
	}
	txIns := make([]wire.TxIn, len(tx.TxIn))
	for i, oldTxIn := range tx.TxIn {
		txIns[i] = *oldTxIn
		txCopy.TxIn[i] = &txIns[i]
	}
	txOuts := make([]wire.TxOut, len(tx.TxOut))
	for i, oldTxOut := range tx.TxOut {
		txOuts[i] = *oldTxOut
		txCopy.TxOut[i] = &txOuts[i]
	}
	return txCopy
}

// payToPubKeyHashScript creates a new script to pay a transaction
// output to a 20-byte pubkey hash. It is expected that the input is a valid
// hash.
func payToPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return tsx.NewScriptBuilder().AddOp(OP_DUP).AddOp(OP_HASH160).
		AddData(pubKeyHash).AddOp(OP_EQUALVERIFY).AddOp(OP_CHECKSIG).
		Script()
}

// payToWitnessPubKeyHashScript creates a new script to pay to a version 0
// pubkey hash witness program. The passed hash is expected to be valid.
func payToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return tsx.NewScriptBuilder().AddOp(OP_0).AddData(pubKeyHash).Script()
}

// payToScriptHashScript creates a new script to pay a transaction output to a
// script hash. It is expected that the input is a valid hash.
func payToScriptHashScript(scriptHash []byte) ([]byte, error) {
	return tsx.NewScriptBuilder().AddOp(OP_HASH160).AddData(scriptHash).
		AddOp(OP_EQUAL).Script()
}

// payToWitnessPubKeyHashScript creates a new script to pay to a version 0
// script hash witness program. The passed hash is expected to be valid.
func payToWitnessScriptHashScript(scriptHash []byte) ([]byte, error) {
	return tsx.NewScriptBuilder().AddOp(OP_0).AddData(scriptHash).Script()
}

// payToPubkeyScript creates a new script to pay a transaction output to a
// public key. It is expected that the input is a valid pubkey.
func payToPubKeyScript(serializedPubKey []byte) ([]byte, error) {
	return tsx.NewScriptBuilder().AddData(serializedPubKey).
		AddOp(OP_CHECKSIG).Script()
}

// LogClosure is a closure that can be printed with %v to be used to
// generate expensive-to-create data for a detailed log level and avoid doing
// the work if the data isn't printed.
type logClosure func() string

func (c logClosure) String() string {
	return c()
}

func newLogClosure(c func() string) logClosure {
	return logClosure(c)
}

// isScriptHash returns true if the script passed is a pay-to-script-hash
// transaction, false otherwise.
func isScriptHash(pops []parsedOpcode) bool {
	return len(pops) == 3 &&
		pops[0].opcode.value == OP_HASH160 &&
		pops[1].opcode.value == OP_DATA_20 &&
		pops[2].opcode.value == OP_EQUAL
}

// isPushOnly returns true if the script only pushes data, false otherwise.
func isPushOnly(pops []parsedOpcode) bool {
	// NOTE: This function does NOT verify opcodes directly since it is
	// internal and is only called with parsed opcodes for scripts that did
	// not have any parse errors.  Thus, consensus is properly maintained.

	for _, pop := range pops {
		// All opcodes up to OP_16 are data push instructions.
		// NOTE: This does consider OP_RESERVED to be a data push
		// instruction, but execution of OP_RESERVED will fail anyways
		// and matches the behavior required by consensus.
		if pop.opcode.value > OP_16 {
			return false
		}
	}
	return true
}

// isWitnessProgram returns true if the passed script is a witness program, and
// false otherwise. A witness program MUST adhere to the following constraints:
// there must be exactly two pops (program version and the program itself), the
// first opcode MUST be a small integer (0-16), the push data MUST be
// canonical, and finally the size of the push data must be between 2 and 40
// bytes.
func isWitnessProgram(pops []parsedOpcode) bool {
	return len(pops) == 2 &&
		isSmallInt(pops[0].opcode) &&
		canonicalPush(pops[1]) &&
		(len(pops[1].data) >= 2 && len(pops[1].data) <= 40)
}

// isSmallInt returns whether or not the opcode is considered a small integer,
// which is an OP_0, or OP_1 through OP_16.
func isSmallInt(op *opcode) bool {
	if op.value == OP_0 || (op.value >= OP_1 && op.value <= OP_16) {
		return true
	}
	return false
}
