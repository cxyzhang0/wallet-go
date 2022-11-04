package lemon

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/onrik/ethrpc"
)

func TestSimplemultisig(t *testing.T) {
	//执行该测试需要本地运行在 port:8545 的 ganache
	const (
		rpcHost = "http://localhost:8545"
	)

	//一些无意义的注释（改动以防止test cache）

	var genAddr0, genAddr1, genAddr2 *addrInfo
	var addrs []*addrInfo
	{ //生成3个地址，并排序
		for i := 0; i < 3; i++ {
			addr, err := genNewAddress()
			FailOnErr(t, err, "FonGetNewAddress_")
			addrs = append(addrs, addr)
		}
		sort.Slice(addrs, func(i, j int) bool {
			return addrs[i].address < addrs[j].address
		})
		genAddr0, genAddr1, genAddr2 = addrs[0], addrs[1], addrs[2]
	}
	prepareFunds4address(t, rpcHost, genAddr0.address, int64(3))

	type args struct {
		rpcHost    string
		privkHex   string
		hexAddress []string
		mRequired  uint8
	}
	var (
		contractAddress string
	)
	{ // 部署合约测试
		arg := args{
			rpcHost:    rpcHost,
			privkHex:   genAddr0.privkHex,
			hexAddress: []string{genAddr0.address, genAddr1.address, genAddr2.address},
			mRequired:  2,
		}

		chainID := big.NewInt(1)
		got, err := DeploySimpleMultiSigContract(rpcHost, *chainID, arg.privkHex, arg.hexAddress, arg.mRequired)
		if err != nil {
			t.Errorf("DeployMultiSigWalletContract() error = %v", err)
			t.FailNow()
		}
		fmt.Println("deployMultisigWalletContract got:", got)

		contractAddress = got
		fmt.Println("contractAddress", contractAddress)
	}

	{ // 部署好后验证合约属性（owners/mRequired）
		time.Sleep(time.Millisecond * 200)
		owners, mRequired, err := GetContractInfo(rpcHost, contractAddress)
		FailOnErr(t, err, "Failed to get contract info")
		fmt.Println("contract info", owners, mRequired)
		FailOnFlag(t, len(owners) != 3, "len owners != 3", len(owners))
		FailOnFlag(t, mRequired != 2, "mRequired != 2", mRequired)
	}

	{ //合约部署后往其中转入资金
		prepareFunds4address(t, rpcHost, contractAddress, 2)
	}

	outAddr := getDefaultOutAddress(t, rpcHost)
	{ // 交易测试
		var (
			sigV                                 []uint8    //签名
			sigR, sigS                           [][32]byte //签名
			privkHex                             string
			multisigContractAddress, fromAddress string //多签合约地址，发起地址
			destination, executor                string //toAddress
			value, gasLimit                      int64
			data                                 []byte
			nonce                                = int64(0)
		)
		// 012由0发起，0和2签名, 把钱赚到1的地址上,executor 为0
		privkHex = genAddr0.privkHex
		multisigContractAddress = contractAddress
		fromAddress = genAddr0.address
		destination = outAddr
		executor = genAddr0.address
		value = 1 * E18
		gasLimit = 239963
		data = []byte("")

		for _, add := range []*addrInfo{genAddr0, genAddr2} {
			v, r, s, err := createSig(add.privkHex, contractAddress, destination, executor, nonce, value, gasLimit, data)
			FailOnErr(t, err, "create sig failed")
			sigV = append(sigV, v)
			sigR = append(sigR, r)
			sigS = append(sigS, s)
		}

		txid, err := ExecuteTX(&TxParams{
			rpcHost:                 rpcHost,
			sigV:                    sigV,
			sigR:                    sigR,
			sigS:                    sigS,
			privkHex:                privkHex,
			multisigContractAddress: multisigContractAddress,
			fromAddress:             fromAddress,
			destination:             destination,
			executor:                executor,
			value:                   value,
			gasLimit:                gasLimit,
			data:                    data,
		})
		FailOnErr(t, err, "Execute Failed")
		fmt.Println("execute txid", txid)
	}

	{ // 完了检查确实转账成功
		time.Sleep(time.Second)
		client, err := ethclient.Dial(rpcHost)
		FailOnErr(t, err, "ConnRpcFail")
		bal, err := client.BalanceAt(context.Background(), common.HexToAddress(outAddr), nil)
		FailOnErr(t, err, "FonGetBal")
		fmt.Println("balance of new tx", bal)
	}
	time.Sleep(time.Second)
}

const E18 = 1000000000000000000

// 为addr准备一定量的eth,
func prepareFunds4address(t *testing.T, rpcHost, addr string, funds int64) {
	rpc := ethrpc.New(rpcHost)
	accounts, err := rpc.EthAccounts()
	FailOnErr(t, err, "FonGetAccounts")
	txid, err := rpc.EthSendTransaction(ethrpc.T{
		From:  accounts[0],
		To:    addr,
		Value: big.NewInt(funds * E18),
	})
	FailOnErr(t, err, "FonSendTX")
	fmt.Printf("已为%s准备资金%d, txid:%s\n", addr, funds, txid)
}

// getDefaultOutAddress 获取预置的地址之2（index 0），这样用ganache的话容易直接查看余额
func getDefaultOutAddress(t *testing.T, rpcHost string) string {
	rpc := ethrpc.New(rpcHost)
	accounts, err := rpc.EthAccounts()
	FailOnErr(t, err, "FonGetAccounts")
	return accounts[1]
}

// FailOnErr used in testing assert
func FailOnErr(t *testing.T, e error, msg string) {
	if e != nil {
		t.Fatalf("Fatal on error, %s, %v", msg, e)
	}
}

func FailOnFlag(t *testing.T, flag bool, params ...interface{}) {
	if flag {
		t.Fatalf("Fail on falg, %v", params)
	}
}
