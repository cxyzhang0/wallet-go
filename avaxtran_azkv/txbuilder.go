package ethtran_azkv

import (
	"context"
	"encoding/hex"
	"fmt"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	contract "github.com/cxyzhang0/wallet-go/ethtran_azkv/contract1"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	log "github.com/sirupsen/logrus"
	"math/big"
	"sort"
)

type TxReq struct {
	From      kmssdk.KeyLabel
	To        kmssdk.KeyLabel
	ToAddr    string
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

	var toAddrPubKey *common.Address
	if req.ToAddr != "" {
		addr := common.HexToAddress(req.ToAddr)
		toAddrPubKey = &addr
		//b, err := hexutil.Decode(req.ToAddr)
		//if err != nil {
		//	return "", "", err
		//}
		//toAddrPubKey = (*common.Address)(b)
	} else {
		addr, _, err := GetAddressPubKey(req.To, sdk)
		if err != nil {
			return "", "", err
		}
		toAddrPubKey = addr
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

	sig, err := GetCompleteSignature(signature, txHash[:], fromAddrPubKey)
	if err != nil {
		return "", "", err
	}

	signedTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return "", "", err
	}

	raw, err := signedTx.MarshalBinary()
	//raw, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", "", fmt.Errorf("unable to encode signed tx to bytes %+v", err)
	}
	return hexutil.Encode(raw), fmt.Sprintf("0x%x", signedTx.Hash().Bytes()), nil

	//return GetSignedTx(signature, tx, signer, txHash, fromAddrPubKey)

	/*
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
	*/
}

type AddressInfo struct {
	KeyLabel kmssdk.KeyLabel
	Address  common.Address
}

type MultisigDeployTxReq struct {
	From         kmssdk.KeyLabel // deploy from
	FromAddress  *common.Address
	MultisigFrom []kmssdk.KeyLabel // for multisig wallet owners
	M            uint32
	Nonce        uint64
	GasLimit     uint64
	GasPrice     *big.Int
}

type MultisigTxReq struct {
	Executor            kmssdk.KeyLabel
	ExecutorAddress     *common.Address
	ContractAddress     *common.Address
	MultisigAddressInfo []*AddressInfo // for multisig wallet owners, to be sorted by Address
	//MultisigFrom    []kmssdk.KeyLabel // for multisig wallet owners
	M                     uint32
	To                    kmssdk.KeyLabel
	ToAddress             *common.Address
	Amount                *big.Int
	ExecutorNonce         uint64
	ContractNonce         uint64
	ContractVariableNonce uint64
	GasLimit              uint64
	GasPrice              *big.Int
	Data                  []byte
}

// BuildDeployContractTx
// Build a tx to deploy multisig contract and deploy it.
func BuildDeployContractTx(req MultisigDeployTxReq, sdk *kmssdk.SDK, client *ethclient.Client, chainID *big.Int) (common.Address, *types.Transaction, *contract.Contract, error) {
	nilAddress := common.Address{}
	//fromAddrPubKey := req.FromAddress

	auth, err := NewKeyedTransactorWithChainID(req.From, sdk /*fromAddrPubKey,*/, chainID)
	if err != nil {
		return nilAddress, nil, nil, err
	}

	auth.Nonce = big.NewInt(int64(req.Nonce))
	auth.Value = big.NewInt(0) // 0 for deploy contract tx
	auth.GasLimit = req.GasLimit
	auth.GasPrice = req.GasPrice

	multisigAddresses := make([]common.Address, len(req.MultisigFrom))
	for i, keyLabel := range req.MultisigFrom {
		addrPubKey, _, err := GetAddressPubKey(keyLabel, sdk)
		if err != nil {
			return nilAddress, nil, nil, err
		}

		multisigAddresses[i] = *addrPubKey
	}

	// multisigAddresses in strictly increasing order
	sort.Slice(multisigAddresses, func(i, j int) bool {
		return multisigAddresses[i].String() < multisigAddresses[j].String()
		//return strings.ToLower(multisigAddresses[i].String()) < strings.ToLower(multisigAddresses[j].String())
	})

	return contract.DeployContract(auth, client, big.NewInt(int64(req.M)), multisigAddresses, chainID)
}

// BuildMultisigTx
// Build a multisig tx and deploy it
func BuildMultisigTx(req MultisigTxReq, sdk *kmssdk.SDK, client *ethclient.Client, chainID *big.Int) (*types.Transaction, error) {
	auth, err := NewKeyedTransactorWithChainID(req.Executor, sdk, chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(req.ExecutorNonce))
	//auth.Value = req.Amount // NOTE: don't set value
	auth.GasLimit = req.GasLimit
	auth.GasPrice = req.GasPrice

	//multisigAddresses := make([]common.Address, len(req.MultisigAddressInfo))
	//for i, addressInfo := range req.MultisigAddressInfo {
	//	addrPubKey, _, err := GetAddressPubKey(addressInfo.KeyLabel, sdk)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	multisigAddresses[i] = *addrPubKey
	//}
	//
	//// multisigAddresses in strictly increasing order
	//sort.Slice(multisigAddresses, func(i, j int) bool {
	//	return strings.ToLower(multisigAddresses[i].String()) < strings.ToLower(multisigAddresses[j].String())
	//})

	sort.Slice(req.MultisigAddressInfo, func(i, j int) bool {
		return req.MultisigAddressInfo[i].Address.String() < req.MultisigAddressInfo[j].Address.String()
		//return strings.ToLower(req.MultisigAddressInfo[i].Address.String()) < strings.ToLower(req.MultisigAddressInfo[j].Address.String())
	})
	// contract instance
	instance, err := contract.NewContract(*req.ContractAddress, client)
	if err != nil {
		return nil, err
	}

	var (
		sigV = make([]uint8, req.M)
		sigR = make([][32]byte, req.M)
		sigS = make([][32]byte, req.M)
	)

	contractDomainSeparatorHash, contractInputHash, contractTotalHash, err := instance.GetHashes(nil, *req.ToAddress, req.Amount, req.Data, *req.ExecutorAddress, big.NewInt(int64(req.GasLimit)))
	log.Infof("contract domainSeperator hex: %s", hex.EncodeToString(contractDomainSeparatorHash[:]))
	log.Infof("contract txInputHash hex: %s", hex.EncodeToString(contractInputHash[:]))
	log.Infof("contract totalHash hex: %s", hex.EncodeToString(contractTotalHash[:]))

	totalHash := getTotalHashForMultisig(req, chainID.Int64())
	for i := 0; i < int(req.M); i++ {
		v, r, s, err := createSig(req, sdk, i, totalHash)
		//v, r, s, err := createSig(req, sdk, i, contractTotalHash[:])
		//v, r, s, err := createSig(req, sdk, i, totalHash)
		if err != nil {
			return nil, err
		}
		sigV[i] = v
		sigR[i] = r
		sigS[i] = s
	}

	//return instance.Execute(auth, sigV, sigR, sigS, *toAddrPubKey, req.Amount, req.Data, *req.ExecutorAddress, big.NewInt(int64(req.GasLimit)))
	//signedTx, err := instance.Execute(auth, sigV, sigR, sigS, *toAddrPubKey, req.Amount, req.Data, *req.ExecutorAddress, big.NewInt(int64(req.GasLimit)))
	signedTx, err := instance.Execute(auth, sigV, sigR, sigS, *req.ToAddress, req.Amount, req.Data, *req.ExecutorAddress, big.NewInt(int64(req.GasLimit)))
	if err != nil {
		return nil, fmt.Errorf("failed to execute: %+v\n", err)
	}

	{ // [DEBUG] experiment with events
		//time.Sleep(time.Minute)
		latestBlock, err := client.BlockNumber(context.Background())
		if err != nil {
			latestBlock = 7880537
		}
		go func() {
			ito, err := instance.FilterExecuteLog(&bind.FilterOpts{
				Start: latestBlock,
			})
			if err != nil {
				log.Errorf("failed to filter execute log: %+v\n", err)
				return
			}
			for {
				if !ito.Next() {
					log.Infof("event no more next")
					break
				}
				evt := ito.Event
				log.Infof("evt seperator: %s\n", hex.EncodeToString(evt.Sperator[:]))
				log.Infof("evt TxInputHash: %s\n", hex.EncodeToString(evt.TxInputHash[:]))
				log.Infof("evt TotalHash: %s\n", hex.EncodeToString(evt.TotalHash[:]))
			}
			log.Infof("event over")
		}()

		go func() {
			ito, err := instance.FilterRecoverdAddr(&bind.FilterOpts{
				Start: latestBlock,
			})
			if err != nil {
				log.Errorf("failed to filter recover addr log: %+v\n", err)
				return
			}
			for {
				if !ito.Next() {
					log.Infof("event no more next")
					break
				}
				evt := ito.Event
				log.Infof("evt recover addr: %s %s\n", evt.I.String(), evt.Addr.Hex())
			}
			log.Infof("event over")
		}()
	}

	return signedTx, nil
	//return contract.DeployContract(auth, client, big.NewInt(int64(req.M)), addresses, chainConfig.ChainID)
}

//func NewKeyedTransactorWithChainID(req MultisigDeployTxReq, sdk *kmssdk.SDK, fromAddrPubKey *common.Address, chainID *big.Int) (*bind.TransactOpts, error) {
func NewKeyedTransactorWithChainID(fromKeyLabel kmssdk.KeyLabel, sdk *kmssdk.SDK /*fromAddrPubKey *common.Address,*/, chainID *big.Int) (*bind.TransactOpts, error) {
	fromAddrPubKey, _, err := GetAddressPubKey(fromKeyLabel, sdk)
	if err != nil {
		return nil, err
	}
	//keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	if chainID == nil {
		return nil, bind.ErrNoChainID
	}
	signer := types.LatestSignerForChainID(chainID)
	return &bind.TransactOpts{
		From: *fromAddrPubKey,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != *fromAddrPubKey {
				return nil, bind.ErrNotAuthorized
			}

			txHash := signer.Hash(tx)

			signature, err := sdk.GetChainSignature(fromKeyLabel, txHash.Bytes())
			if err != nil {
				return nil, err
			}

			sig, err := GetCompleteSignature(signature, txHash[:], fromAddrPubKey)
			if err != nil {
				return nil, err
			}

			//signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			//if err != nil {
			//	return nil, err
			//}
			return tx.WithSignature(signer, sig)
		},
		Context: context.Background(),
	}, nil
}
