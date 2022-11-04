package lemon

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/core/types"

	contracts "github.com/cxyzhang0/wallet-go/ethtran_azkv/contract1"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// DeploySimpleMultiSigContract 创建部署一个多签合约，返回合约地址
// rpcHost: rpc host,address: 签名的参与者， mRequired: m-n 至少多少人签名生效
// return: address_hex,txid
func DeploySimpleMultiSigContract(rpcHost string, chainID big.Int, privkHex string, hexAddress []string, mRequired uint8) (string, error) {
	//TODO 参数验证:rpchost格式，chainID枚举,privkHex格式,hexAddress格式和排序和长度
	var (
		client *ethclient.Client
		err    error
		privk  *ecdsa.PrivateKey
		owners []common.Address
	)

	{ // 预处理
		//owners需要保持incr序，合约内的简单查重依赖于排序性
		sort.Slice(hexAddress, func(i, j int) bool {
			return hexAddress[i] < hexAddress[j]
		})
	}

	{ // init vars
		client, err = ethclient.Dial(rpcHost)
		if err != nil {
			return "", fmt.Errorf("无法连接到rpcHost, %v", err)
		}
		defer client.Close()

		privk, err = crypto.HexToECDSA(privkHex)
		if err != nil {
			return "", fmt.Errorf("转换私钥时发生错误,%v", err)
		}

		for _, ha := range hexAddress {
			owners = append(owners, common.HexToAddress(ha))
		}
	}

	{ // 部署多签合约
		auth := bind.NewKeyedTransactor(privk)
		//TODO set gas,gasPrice
		add, tx, w, err := contracts.DeployContract(auth, client, big.NewInt(int64(mRequired)), owners, &chainID)
		if err != nil {
			return "", fmt.Errorf("部署多签合约失败, %v", err)
		}
		time.Sleep(time.Second)

		ver, err := w.GetVersion(&bind.CallOpts{Pending: true})
		if err != nil {
			return "", fmt.Errorf("无法调用合约函数获取version, %v", err)
		}
		// return ver + "," + add.Hex() + "," + tx.Hash().Hex(), nil
		info := add.Hex() + "," + tx.Hash().Hex() + ", ver:" + ver
		fmt.Println("deploy info", info)

		return add.Hex(), nil
	}
}

// TxParams 交易参数
type TxParams struct {
	rpcHost                              string     //rpc host
	sigV                                 []uint8    //签名
	sigR, sigS                           [][32]byte //签名
	privkHex                             string
	multisigContractAddress, fromAddress string //多签合约地址，发起地址
	destination, executor                string //toAddress
	value, gasLimit                      int64
	data                                 []byte
}

// ExecuteTX .
func ExecuteTX(txp *TxParams) (string, error) {
	var (
		client           *ethclient.Client
		err              error
		privk            *ecdsa.PrivateKey
		multisigContract *contracts.Contract
	)

	{ // init vars
		client, err = ethclient.Dial(txp.rpcHost)
		if err != nil {
			return "", fmt.Errorf("无法连接到rpcHost, %v", err)
		}
		defer client.Close()

		privk, err = crypto.HexToECDSA(txp.privkHex)
		if err != nil {
			return "", err
		}

		multisigContract, err = contracts.NewContract(common.HexToAddress(txp.multisigContractAddress), client)
		if err != nil {
			return "", fmt.Errorf("构建多签合约调用时异常,检查合约地址和rpc server,%v", err)
		}
	}

	//TODO 参数校验
	{ // 调用合约方法
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(txp.multisigContractAddress))
		if err != nil {
			return "", fmt.Errorf("获取多签地址nonce时发生错误, %v", err)
		}

		var signerFn bind.SignerFn = func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			// if address != keyAddr {
			// 	return nil, errors.New("not authorized to sign this account")
			// }
			signer := types.LatestSignerForChainID(big.NewInt(chainiD))

			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), privk)
			if err != nil {
				return nil, err
			}
			// return tx.WithSignature(signer, append(signature, byte(12))) 这里瞎写的，测试失败的情况，可以删掉
			return tx.WithSignature(signer, signature)
		}
		tx, err := multisigContract.Execute(&bind.TransactOpts{
			From:   common.HexToAddress(txp.fromAddress),
			Nonce:  big.NewInt(int64(nonce)),
			Signer: signerFn,
			/*func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				// if address != keyAddr {
				// 	return nil, errors.New("not authorized to sign this account")
				// }
				signature, err := crypto.Sign(signer.Hash(tx).Bytes(), privk)
				if err != nil {
					return nil, err
				}
				// return tx.WithSignature(signer, append(signature, byte(12))) 这里瞎写的，测试失败的情况，可以删掉
				return tx.WithSignature(signer, signature)
			},*/
		},
			txp.sigV,
			txp.sigR,
			txp.sigS,
			common.HexToAddress(txp.destination),
			big.NewInt(int64(txp.value)),
			txp.data,
			common.HexToAddress(txp.executor),
			big.NewInt(txp.gasLimit))
		if err != nil {
			return "", fmt.Errorf("调用合约交易方法时发生错误, %v", err)
		}

		{ // [DEBUG] 尝试读取事件日志
			go func() {
				ito, err := multisigContract.FilterExecuteStart(&bind.FilterOpts{
					Start: 0,
				})
				if err != nil {
					log.Error("无法过滤reLog", err)
					return
				}
				for {
					if !ito.Next() {
						log.Info("事件迭代no more next")
						break
					}
					evt := ito.Event
					log.Infof("evt execute start: %+v", evt)
				}
				log.Info("事件日志over")
			}()

			go func() {
				ito, err := multisigContract.FilterExecuteVerifySender(&bind.FilterOpts{
					Start: 0,
				})
				if err != nil {
					log.Error("无法过滤reLog", err)
					return
				}
				for {
					if !ito.Next() {
						log.Info("事件迭代no more next")
						break
					}
					evt := ito.Event
					log.Infof("evt execute verify sender: %s", evt.Sender.Hex())
				}
				log.Info("事件日志over")
			}()

			go func() {
				ito, err := multisigContract.FilterExecuteLog(&bind.FilterOpts{
					Start: 0,
				})
				if err != nil {
					log.Error("无法过滤executeLog", err)
					return
				}
				for {
					if !ito.Next() {
						log.Info("事件迭代no more next")
						break
					}
					evt := ito.Event
					log.Info("evt seperator:", hex.EncodeToString(evt.Sperator[:]))
					log.Info("evt TxInputHash:", hex.EncodeToString(evt.TxInputHash[:]))
					log.Info("evt TotalHash:", hex.EncodeToString(evt.TotalHash[:]))
				}
				log.Info("事件日志over")
			}()

			go func() {
				ito, err := multisigContract.FilterRecoverStart(&bind.FilterOpts{
					Start: 0,
				})
				if err != nil {
					log.Error("无法过滤reLog", err)
					return
				}
				for {
					if !ito.Next() {
						log.Info("事件迭代no more next")
						break
					}
					evt := ito.Event
					log.Infof("evt recover start with index: %s", evt.I.String())
				}
				log.Info("事件日志over")
			}()

			go func() {
				ito, err := multisigContract.FilterRecoverVerify(&bind.FilterOpts{
					Start: 0,
				})
				if err != nil {
					log.Error("无法过滤reLog", err)
					return
				}
				for {
					if !ito.Next() {
						log.Info("事件迭代no more next")
						break
					}
					evt := ito.Event
					log.Infof("evt recover verify for addr index: %s; address: %s", evt.I.String(), evt.Addr.Hex())
				}
				log.Info("事件日志over")
			}()

			go func() {
				ito, err := multisigContract.FilterRecoverdAddr(&bind.FilterOpts{
					Start: 0,
				})
				if err != nil {
					log.Error("无法过滤reLog", err)
					return
				}
				for {
					if !ito.Next() {
						log.Info("事件迭代no more next")
						break
					}
					evt := ito.Event
					log.Info("evt recover addr: ", evt.I.String(), " ", evt.Addr.Hex())
				}
				log.Info("事件日志over")
			}()
		}

		return tx.Hash().Hex(), nil
	}

}

const (
	txtypeHash           = "0x3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d7"
	nameHash             = "0xb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e6"
	versionHash          = "0xc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6"
	eip712DomaintypeHash = "0xd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac56472"
	salt                 = "0x251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc0"
	chainiD              = 1
	allZero              = "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" //做padding用
)

//return v,r,s
func createSig(signerPrivkHex string, multisigAddr, destinationAddr, executor string, nonce, value, gasLimit int64, data []byte) (uint8, [32]byte, [32]byte, error) {
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

// GetContractInfo 获取多签合约地址内的地址和签名人数
func GetContractInfo(rpcHost, contractAddress string) (addrs []string, mRequired int64, err error) {
	var client *ethclient.Client
	var multisigContract *contracts.ContractCaller

	{ // init vars
		client, err = ethclient.Dial(rpcHost)
		if err != nil {
			return nil, 0, fmt.Errorf("无法连接到rpcHost, %v", err)
		}
		defer client.Close()

		multisigContract, err = contracts.NewContractCaller(common.HexToAddress(contractAddress), client)
		if err != nil {
			return nil, 0, fmt.Errorf("构建多签合约调用时异常,检查合约地址和rpc server,%v", err)
		}
	}
	callOpts := &bind.CallOpts{}
	{ //获取地址
		length, err := multisigContract.GetOwersLength(callOpts)
		if err != nil {
			return nil, 0, fmt.Errorf("获取持有人数量时发生异常, %v", err)
		}

		for i := int64(0); i < length.Int64(); i++ {
			addr, err := multisigContract.OwnersArr(callOpts, big.NewInt(i))
			if err != nil {
				return addrs, 0, fmt.Errorf("获取地址[%d/%d]时发生错误, %v", i+1, length.Int64(), err)
			}
			addrs = append(addrs, addr.Hex())
		}
	}
	{ // 获取mRequired
		mInt, err := multisigContract.Threshold(callOpts)
		if err != nil {
			return nil, 0, fmt.Errorf("获取签名数时发生错误，%v", err)
		}
		mRequired = mInt.Int64()
	}
	return
}
