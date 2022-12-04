package test

import (
	tran "github.com/cxyzhang0/wallet-go/avaxtran_azkv"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
)

// TestBuildDeployERC20Contract
/*
soft executor address1: 0xBfDdb90EFD9323db8E9cB136cE27E5b38fE29Da8
hsm executor address 1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
burnable, decimals: 2; soft executor
contract address: 0xED1dDBcF4246688A82F7da8Ea5015f76c4D13CB4
decimals: 2; hsm executor
contract address: 0xdbb353c5D6c3F7987b98B6aaFFeE129fA0b542E0
decimals: 18; hsm executor
contract address: 0x0079d779C4736b190E6568818659e4B23DE11D9c
*/
func TestBuildDeployERC20Contract(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key:       keyName,
		Version:   "0eab9a0cc2e84018be05f90e5d914142",
		Algorithm: kmssdk.Secp256k1,
	}
	/*
		from := kmssdk.KeyLabel{
			Key:     "secp256k1-hsm-1",
			Version: "cb848fb15e3a40b49bc41cbe957ea438",
			//Version: "0179a6204ed7491ea5b27a87b541d5cb",
			//Version:   "b6aec266b6a147f7a1c40fe842504650",
			Algorithm: kmssdk.Secp256k1,
		}
	*/
	req := tran.ERC20DeployTxReq{
		From:     from,
		GasLimit: 2000000, // TODO: how to calc it?
	}

	fromAddress, _, err := tran.GetAddressPubKey(req.From, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")
	req.FromAddress = fromAddress

	client, err := ethclient.Dial(avaxURL)
	FailOnErr(t, err, "FonDial")

	gasPrice, err := client.SuggestGasPrice(avavCtx)
	FailOnErr(t, err, "FonSuggestGasPrice")
	req.GasPrice = gasPrice

	nonce, err := client.PendingNonceAt(avavCtx, *fromAddress)
	FailOnErr(t, err, "FonPendingNonceAt")
	req.Nonce = nonce

	contractAddress, signedTx, _contract, err := tran.BuildDeployERC20ContractTx(req, _sdk, client, avaxCChainID)
	FailOnErr(t, err, "FonBuildDeployERC20ContractTx")

	raw, err := signedTx.MarshalBinary()
	FailOnErr(t, err, "FonMarshalBinary")

	t.Logf("contract address: %s\nsigned tx: %s\ntx hash: %s\ncontract: %+v", contractAddress.String(), hexutil.Encode(raw), signedTx.Hash().String(), _contract)
}

func TestTotalSupply(t *testing.T) {
	contractAddress := common.HexToAddress("0xED1dDBcF4246688A82F7da8Ea5015f76c4D13CB4")
	//contractAddress := common.HexToAddress("0xdbb353c5D6c3F7987b98B6aaFFeE129fA0b542E0")

	executor := kmssdk.KeyLabel{
		Key:     keyName,
		Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version: "0ff7adfdbe0a4b69881c4dac6b0f81f4",
		//Version:   "485248105ef047aaa1f33cf0baaf9a03",
		Algorithm: kmssdk.Secp256k1,
	}

	/*
		executor := kmssdk.KeyLabel{
			Key:     "secp256k1-hsm-1",
			Version: "cb848fb15e3a40b49bc41cbe957ea438",
			//Version: "0179a6204ed7491ea5b27a87b541d5cb",
			//Version:   "b6aec266b6a147f7a1c40fe842504650",
			Algorithm: kmssdk.Secp256k1,
		}
	*/
	executorAddress, _, err := tran.GetAddressPubKey(executor, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")

	req := tran.ERC20TxReq{
		ContractAddress: &contractAddress,
		Executor:        executor,
		ExecutorAddress: executorAddress,
	}

	client, err := ethclient.Dial(avaxURL)
	FailOnErr(t, err, "FonDial")

	totalSupply, err := tran.TotalSupply(req, client)
	FailOnErr(t, err, "FonTotalSupply")

	t.Logf("total supply at contract %s: %.2f", contractAddress.String(), totalSupply)
}

func TestBalanceOf(t *testing.T) {
	contractAddress := common.HexToAddress("0xED1dDBcF4246688A82F7da8Ea5015f76c4D13CB4")
	//contractAddress := common.HexToAddress("0xdbb353c5D6c3F7987b98B6aaFFeE129fA0b542E0")

	executor := kmssdk.KeyLabel{
		Key: keyName,
		//Version: "0eab9a0cc2e84018be05f90e5d914142",
		//Version: "0ff7adfdbe0a4b69881c4dac6b0f81f4",
		Version:   "485248105ef047aaa1f33cf0baaf9a03",
		Algorithm: kmssdk.Secp256k1,
	}
	/*
		executor := kmssdk.KeyLabel{
			Key:     "secp256k1-hsm-1",
			Version: "cb848fb15e3a40b49bc41cbe957ea438",
			//Version: "0179a6204ed7491ea5b27a87b541d5cb",
			//Version:   "b6aec266b6a147f7a1c40fe842504650",
			Algorithm: kmssdk.Secp256k1,
		}
	*/
	executorAddress, _, err := tran.GetAddressPubKey(executor, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")

	req := tran.ERC20TxReq{
		ContractAddress: &contractAddress,
		Executor:        executor,
		ExecutorAddress: executorAddress,
	}

	client, err := ethclient.Dial(avaxURL)
	FailOnErr(t, err, "FonDial")

	balance, err := tran.BalanceOf(req, client)
	FailOnErr(t, err, "FonBalanceOf")

	t.Logf("balance of %s at contract %s: %.2f", req.ExecutorAddress.String(), contractAddress.String(), balance)
}

func TestMint(t *testing.T) {
	contractAddress := common.HexToAddress("0xED1dDBcF4246688A82F7da8Ea5015f76c4D13CB4")
	//contractAddress := common.HexToAddress("0xdbb353c5D6c3F7987b98B6aaFFeE129fA0b542E0")

	executor := kmssdk.KeyLabel{
		Key:       keyName,
		Version:   "0eab9a0cc2e84018be05f90e5d914142",
		Algorithm: kmssdk.Secp256k1,
	}

	/*
		executor := kmssdk.KeyLabel{
			Key:     "secp256k1-hsm-1",
			Version: "cb848fb15e3a40b49bc41cbe957ea438",
			//Version: "0179a6204ed7491ea5b27a87b541d5cb",
			//Version:   "b6aec266b6a147f7a1c40fe842504650",
			Algorithm: kmssdk.Secp256k1,
		}
	*/
	executorAddress, _, err := tran.GetAddressPubKey(executor, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")

	to := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version: "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}
	toAddress, _, err := tran.GetAddressPubKey(to, _sdk)

	req := tran.ERC20TxReq{
		ContractAddress: &contractAddress,
		Executor:        executor,
		ExecutorAddress: executorAddress,
		To:              &to,
		ToAddress:       toAddress,
		Amount:          big.NewInt(1000e2), // USD pennies
		GasLimit:        1500000,            // TODO: how to calc it?
		Data:            []byte(""),
	}

	client, err := ethclient.Dial(avaxURL)
	FailOnErr(t, err, "FonDial")

	gasPrice, err := client.SuggestGasPrice(avavCtx)
	FailOnErr(t, err, "FonSuggestGasPrice")
	req.GasPrice = gasPrice

	nonce, err := client.PendingNonceAt(avavCtx, *req.ExecutorAddress)
	FailOnErr(t, err, "FonPendingNonceAt")
	req.ExecutorNonce = nonce

	signedTx, err := tran.BuildMintTx(req, _sdk, client, avaxCChainID)
	FailOnErr(t, err, "FonBuildMintTx")

	raw, err := signedTx.MarshalBinary()
	FailOnErr(t, err, "FonMarshalBinary")

	t.Logf("contract address: %s\nsigned tx: %s\ntx hash: %s", req.ContractAddress.String(), hexutil.Encode(raw), signedTx.Hash().String())

}

func TestTransfer(t *testing.T) {
	contractAddress := common.HexToAddress("0xED1dDBcF4246688A82F7da8Ea5015f76c4D13CB4")
	//contractAddress := common.HexToAddress("0xdbb353c5D6c3F7987b98B6aaFFeE129fA0b542E0")

	executor := kmssdk.KeyLabel{
		Key: keyName,
		//Version: "0eab9a0cc2e84018be05f90e5d914142",
		Version: "0ff7adfdbe0a4b69881c4dac6b0f81f4",
		//Version:   "485248105ef047aaa1f33cf0baaf9a03",
		Algorithm: kmssdk.Secp256k1,
	}

	/*
		executor := kmssdk.KeyLabel{
			Key:     "secp256k1-hsm-1",
			Version: "cb848fb15e3a40b49bc41cbe957ea438",
			//Version: "0179a6204ed7491ea5b27a87b541d5cb",
			//Version:   "b6aec266b6a147f7a1c40fe842504650",
			Algorithm: kmssdk.Secp256k1,
		}
	*/
	executorAddress, _, err := tran.GetAddressPubKey(executor, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")

	to := kmssdk.KeyLabel{
		Key: keyName,
		//Version: "0eab9a0cc2e84018be05f90e5d914142",
		Version: "0ff7adfdbe0a4b69881c4dac6b0f81f4",
		//Version:   "485248105ef047aaa1f33cf0baaf9a03",
		Algorithm: kmssdk.Secp256k1,
	}
	/*
		to := kmssdk.KeyLabel{
			Key: "secp256k1-hsm-1",
			//Version: "cb848fb15e3a40b49bc41cbe957ea438",
			//Version: "0179a6204ed7491ea5b27a87b541d5cb",
			Version:   "b6aec266b6a147f7a1c40fe842504650",
			Algorithm: kmssdk.Secp256k1,
		}
	*/
	//toAddress, _, err := tran.GetAddressPubKey(to, _sdk)
	toAddress := common.HexToAddress("0xBF272Dc09F22fa85e14D2f28c3a2B4227E28F6ec")
	req := tran.ERC20TxReq{
		ContractAddress: &contractAddress,
		Executor:        executor,
		ExecutorAddress: executorAddress,
		To:              &to,
		ToAddress:       &toAddress,
		Amount:          big.NewInt(100e2), // USD pennies
		GasLimit:        1500000,           // TODO: how to calc it?
		Data:            []byte(""),
	}

	client, err := ethclient.Dial(avaxURL)
	FailOnErr(t, err, "FonDial")

	gasPrice, err := client.SuggestGasPrice(avavCtx)
	FailOnErr(t, err, "FonSuggestGasPrice")
	req.GasPrice = gasPrice

	nonce, err := client.PendingNonceAt(avavCtx, *req.ExecutorAddress)
	FailOnErr(t, err, "FonPendingNonceAt")
	req.ExecutorNonce = nonce

	signedTx, err := tran.BuildTransferTx(req, _sdk, client, avaxCChainID)
	FailOnErr(t, err, "FonBuildTransferTx")

	raw, err := signedTx.MarshalBinary()
	FailOnErr(t, err, "FonMarshalBinary")

	t.Logf("contract address: %s\nsigned tx: %s\ntx hash: %s", req.ContractAddress.String(), hexutil.Encode(raw), signedTx.Hash().String())

}

func TestBurn(t *testing.T) {
	contractAddress := common.HexToAddress("0xED1dDBcF4246688A82F7da8Ea5015f76c4D13CB4")

	executor := kmssdk.KeyLabel{
		Key: keyName,
		//Version: "0eab9a0cc2e84018be05f90e5d914142",
		Version: "0ff7adfdbe0a4b69881c4dac6b0f81f4",
		//Version:   "485248105ef047aaa1f33cf0baaf9a03",
		Algorithm: kmssdk.Secp256k1,
	}

	/*
		executor := kmssdk.KeyLabel{
			Key:     "secp256k1-hsm-1",
			Version: "cb848fb15e3a40b49bc41cbe957ea438",
			//Version: "0179a6204ed7491ea5b27a87b541d5cb",
			//Version:   "b6aec266b6a147f7a1c40fe842504650",
			Algorithm: kmssdk.Secp256k1,
		}
	*/
	executorAddress, _, err := tran.GetAddressPubKey(executor, _sdk)
	FailOnErr(t, err, "FonGetAddressPubKey")

	req := tran.ERC20TxReq{
		ContractAddress: &contractAddress,
		Executor:        executor,
		ExecutorAddress: executorAddress,
		Amount:          big.NewInt(100e2), // USD pennies
		GasLimit:        1500000,           // TODO: how to calc it?
		Data:            []byte(""),
	}

	client, err := ethclient.Dial(avaxURL)
	FailOnErr(t, err, "FonDial")

	gasPrice, err := client.SuggestGasPrice(avavCtx)
	FailOnErr(t, err, "FonSuggestGasPrice")
	req.GasPrice = gasPrice

	nonce, err := client.PendingNonceAt(avavCtx, *req.ExecutorAddress)
	FailOnErr(t, err, "FonPendingNonceAt")
	req.ExecutorNonce = nonce

	signedTx, err := tran.BuildBurnTx(req, _sdk, client, avaxCChainID)
	FailOnErr(t, err, "FonBuildBurnTx")

	raw, err := signedTx.MarshalBinary()
	FailOnErr(t, err, "FonMarshalBinary")

	t.Logf("contract address: %s\nsigned tx: %s\ntx hash: %s", req.ContractAddress.String(), hexutil.Encode(raw), signedTx.Hash().String())

}
