package ethtran_azkv

import (
	"encoding/asn1"
	"fmt"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

type TxReq struct {
	From      kmssdk.KeyLabel
	To        kmssdk.KeyLabel
	Amount    *big.Int
	Nonce     uint64
	GasTipCap *big.Int
	GasFeeCap *big.Int
	Gas       uint64
}

// BuildTx builds a transfer transaction and returns a signed ethereum transaction request and its hash,
// ready to be submitted to the blockchain as determined by chainConfig (e.g., params.GoerliChainConfig).
// req provides the usual from, to and amount as well as nonce and gas fee related parameters.
// See TestBuildTx for how the rest of parameters are calc'd on results from querying Blockdaemon ubiquity API.
// TODO: A future refactoring of this func into CoreTx may call ubiquity API directly from CoreTx, in which case,
// the nonce and gas fees need not come from the caller of this func.
// Also, the next iteration may not need chainConfig as a parameter - the caller just needs to pass in the
// network (mainnet or one of the testnets) so the ubiquity call and this func are consistent re network.
func BuildTx(req TxReq, sdk *kmssdk.SDK, chainConfig *params.ChainConfig) (string, string, error) { // signed raw tx, tx hash, error
	fromAddrPubKey, _, err := GetAddressPubKey(req.From, sdk)
	if err != nil {
		return "", "", err
	}

	toAddrPubKey, _, err := GetAddressPubKey(req.To, sdk)
	if err != nil {
		return "", "", err
	}

	var data []byte

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     req.Nonce,
		GasTipCap: req.GasTipCap,
		GasFeeCap: req.GasFeeCap,
		Gas:       req.Gas,
		To:        toAddrPubKey,
		Value:     req.Amount,
		Data:      data,
	})

	signer := types.MakeSigner(chainConfig, chainConfig.LondonBlock)

	txHash := signer.Hash(tx)

	signature, err := sdk.GetChainSignature(req.From, txHash.Bytes())
	if err != nil {
		return "", "", err
	}

	var pubKeyAddr func([]byte) common.Address // TODO: is it more efficient to have it as a standard external func?
	pubKeyAddr = func(bytes []byte) common.Address {
		digest := crypto.Keccak256(bytes[1:])
		var addr common.Address
		copy(addr[:], digest[12:])
		return addr
	}

	// parse sig
	var params struct{ R, S *big.Int }
	_, err = asn1.Unmarshal(signature, &params)
	if err != nil {
		return "", "", fmt.Errorf("asymmetric signature encoding: %w", err)
	}
	var rLen, sLen int // byte size
	if params.R != nil {
		rLen = (params.R.BitLen() + 7) / 8
	}
	if params.S != nil {
		sLen = (params.S.BitLen() + 7) / 8
	}
	if rLen == 0 || rLen > 32 || sLen == 0 || sLen > 32 {
		return "", "", fmt.Errorf("asymmetric signature with %d-byte r and %d-byte s denied on size", rLen, sLen)
	}

	// Need uncompressed signature with "recovery ID" at end:
	// https://bitcointalk.org/index.php?topic=5249677.0
	// https://ethereum.stackexchange.com/a/53182/39582
	var sig [66]byte // + 1-byte header + 1-byte tailer
	params.R.FillBytes(sig[33-rLen : 33])
	params.S.FillBytes(sig[65-sLen : 65])

	// brute force try includes KMS verification
	var recoverErr error
	for recoveryID := byte(0); recoveryID < 2; recoveryID++ {
		sig[0] = recoveryID + 27 // BitCoin header
		btcsig := sig[:65]       // exclude Ethereum 'v' parameter
		pubKey, _, err := btcecdsa.RecoverCompact(btcsig, txHash[:])
		if err != nil {
			recoverErr = err
			continue
		}

		if pubKeyAddr(pubKey.SerializeUncompressed()) == *fromAddrPubKey {
			// sign the transaction
			sig[65] = recoveryID // Ethereum 'v' parameter

			signedTx, err := tx.WithSignature(signer, sig[1:])
			if err != nil {
				return "", "", err
			}

			raw, err := signedTx.MarshalBinary()
			//raw, err := rlp.EncodeToBytes(signedTx)
			if err != nil {
				return "", "", fmt.Errorf("unable to encode signed tx to bytes %+v", err)
			}
			return hexutil.Encode(raw), fmt.Sprintf("0x%x", signedTx.Hash().Bytes()), nil
			//return fmt.Sprintf("0x%x", raw), fmt.Sprintf("0x%x", signedTx.Hash().Bytes()), nil
			//return hex.EncodeToString(raw), signedTx.Hash().String(), nil

			//// sign the transaction
			//sig[65] = recoveryID // Ethereum 'v' parameter
			//etcsig := sig[1:]    // exclude BitCoin header
			//signedTx, err := tx.WithSignature(signer, etcsig)
			//if err == nil {
			//	return "", "", err
			//}
			//
			//raw, err := rlp.EncodeToBytes(signedTx)
			//if err != nil {
			//	return "", "", fmt.Errorf("unable to encode signed tx to bytes %+v", err)
			//}
			//return hex.EncodeToString(raw), signedTx.Hash().String(), nil
		}
	}
	// recoverErr can be nil, but that's OK
	return "", "", fmt.Errorf("asymmetric signature address recovery mis: %w", recoverErr)

}
