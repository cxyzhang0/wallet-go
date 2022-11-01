package ethtran_azkv

import (
	"bytes"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strconv"
)

func GetAddressPubKey(keyLabel kmssdk.KeyLabel, sdk *kmssdk.SDK) (*common.Address, string, error) { // address pub key, address, error
	pubkey, err := sdk.GetECDSAPublicKey(keyLabel)
	if err != nil {
		return nil, "", err
	}

	addrPubKey := crypto.PubkeyToAddress(*pubkey)

	addr := addrPubKey.Hex()

	return &addrPubKey, addr, nil
}

func StringToBigInt(str string) (*big.Int, error) {
	bigFloat, _, err := big.ParseFloat(str, 10, 0, big.ToZero)
	if err != nil {
		return nil, fmt.Errorf("failed to pars big int: %+v", err)
	}

	bigInt, _ := bigFloat.Int(nil)

	return bigInt, nil
}

func GetSignedTx(signature []byte, tx *types.Transaction, signer types.Signer, txHash common.Hash, fromAddrPubKey *common.Address) (string, string, error) { // signed raw tx, tx hash, error
	var pubKeyAddr func([]byte) common.Address // TODO: is it more efficient to have it as a standard external func?
	pubKeyAddr = func(bytes []byte) common.Address {
		digest := crypto.Keccak256(bytes[1:])
		var addr common.Address
		copy(addr[:], digest[12:])
		return addr
	}

	// parse sig
	var params struct{ R, S *big.Int }
	_, err := asn1.Unmarshal(signature, &params)
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

func GetCompleteSignature(signature []byte, txHash common.Hash, fromAddrPubKey *common.Address) ([]byte, error) { // signed raw tx, tx hash, error
	var pubKeyAddr func([]byte) common.Address // TODO: is it more efficient to have it as a standard external func?
	pubKeyAddr = func(bytes []byte) common.Address {
		digest := crypto.Keccak256(bytes[1:])
		var addr common.Address
		copy(addr[:], digest[12:])
		return addr
	}

	// parse sig
	var params struct{ R, S *big.Int }
	_, err := asn1.Unmarshal(signature, &params)
	if err != nil {
		return nil, fmt.Errorf("asymmetric signature encoding: %w", err)
	}
	var rLen, sLen int // byte size
	if params.R != nil {
		rLen = (params.R.BitLen() + 7) / 8
	}
	if params.S != nil {
		sLen = (params.S.BitLen() + 7) / 8
	}
	if rLen == 0 || rLen > 32 || sLen == 0 || sLen > 32 {
		return nil, fmt.Errorf("asymmetric signature with %d-byte r and %d-byte s denied on size", rLen, sLen)
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

			if err != nil {
				return nil, fmt.Errorf("unable to encode signed tx to bytes %+v", err)
			}
			return sig[1:], nil
		}
	}
	// recoverErr can be nil, but that's OK
	return nil, fmt.Errorf("asymmetric signature address recovery mis: %w", recoverErr)
}

const (
	txtypeHash           = "0x3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d7"
	nameHash             = "0xb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e6"
	versionHash          = "0xc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6"
	eip712DomaintypeHash = "0xd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac56472"
	salt                 = "0x251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc0"
	//chainiD              = 1
	allZero = "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" //做padding用
)

//return v,r,s
func createSig(signerPrivkHex string, multisigAddr, destinationAddr, executor string, nonce, value, gasLimit int64, data []byte, chainiD int64) (uint8, [32]byte, [32]byte, error) {
	// log.WithFields(log.Fields{
	// 	"multiAddr": multisigAddr,
	// 	"destAddr":  destinationAddr,
	// 	"executor":  executor,
	// 	"nonce":     nonce,
	// 	"gasLmt":    gasLimit,
	// 	"value":     value,
	// 	"data":      data,
	// }).Info("createSig")

	privk, err := crypto.HexToECDSA(signerPrivkHex)
	if err != nil {
		panic(err)
	}

	leftPad2Str := func(str string) string { // 将小于64位的字符串(hex编码的)填充至64位（64位转为byte即32位，对应32*8=256 bit）
		needed := 64 - len(str)
		return allZero[:needed] + str
	}
	i2hex := func(i int64) string { //转为16进制字符串
		return strconv.FormatInt(i, 16)
		// return fmt.Sprintf("%x", i)
	}
	hexToKeccak256ThenHex := func(byts []byte) string { // 将hex编码的字符串的字节串decode为字节串，然后进行keccak256Hash,返回hex输出
		if bytes.Index(byts, []byte("0x")) == 0 {
			byts = byts[2:]
		}

		decodedData, err := hex.DecodeString(string(byts))
		if err != nil {
			fmt.Println("byts:", string(byts))
			panic(err)
		}
		return crypto.Keccak256Hash([]byte(decodedData)).Hex()
	}
	localKeccak256 := func(byts []byte) []byte {
		if bytes.Index(byts, []byte("0x")) == 0 {
			byts = byts[2:]
		}

		decodedData, err := hex.DecodeString(string(byts))
		if err != nil {
			fmt.Println("byts:", string(byts))
			panic(err)
		}
		return crypto.Keccak256([]byte(decodedData))
	}

	domainData := eip712DomaintypeHash + nameHash[2:] + versionHash[2:] + leftPad2Str(i2hex(chainiD)) + leftPad2Str(multisigAddr[2:]) + salt[2:]
	domainSeparatorHashHex := hexToKeccak256ThenHex([]byte(domainData))
	txInput := txtypeHash + leftPad2Str(destinationAddr[2:]) + leftPad2Str(i2hex(value)) + hexToKeccak256ThenHex(data)[2:] + leftPad2Str(i2hex(nonce)) + leftPad2Str(executor[2:]) + leftPad2Str(i2hex(gasLimit))
	// fmt.Println("[DBG](txInput)", txInput)
	txInputHashHex := hexToKeccak256ThenHex([]byte(txInput))

	input := "0x19" + "01" + domainSeparatorHashHex[2:] + txInputHashHex[2:]
	// log.Info("[DBG](txInputHashHex,input)", txInputHashHex, input)
	// log.Info("domainData:  ", domainData)
	// log.Info("domainSeperator: ", domainSeparatorHashHex)
	// log.Info("txInput:  ", txInput)
	// log.Info("txInputHashHex:  ", txInputHashHex)
	// log.Info("input:  ", input)

	hashBytes := localKeccak256([]byte(input))
	// log.Info("totalHash:", hex.EncodeToString(hashBytes))

	sig, err := crypto.Sign(hashBytes, privk)
	// fmt.Println("len of sig: ", len(sig))
	if err != nil {
		panic(err)
		// return 0, nil, nil , fmt.Errorf("签名失败,%v", err)
	}
	r, s, v := sig[:32], sig[32:64], uint8(int(sig[64]))+27

	{ //【调试用】做内部的ecrecover验证,可移除
		go func() {
			rePub, err := crypto.SigToPub(hashBytes, sig)
			// rePub, err := crypto.Ecrecover([]byte(hash), sig)
			if err != nil {
				panic(fmt.Errorf("ecrecover err: %v", err))
			}
			reAddr := crypto.PubkeyToAddress(*rePub)
			addrFromPriv := crypto.PubkeyToAddress(privk.PublicKey)
			fmt.Println("addrFromPrivKey vs recoverdAddr")
			fmt.Println(addrFromPriv.Hex())
			fmt.Println(reAddr.Hex())
		}()
	}
	toBytes32 := func(b []byte) [32]byte {
		b32 := new([32]byte)
		if len(b) <= 32 {
			copy(b32[:], b)
		} else {
			panic(fmt.Sprintf("overflow byte(32),actual: %d", len(b)))
		}
		return *b32
	}
	return v, toBytes32(r), toBytes32(s), nil
}
