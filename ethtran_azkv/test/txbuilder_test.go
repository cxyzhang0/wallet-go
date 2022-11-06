package test

import (
	"context"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	tran "github.com/cxyzhang0/wallet-go/ethtran_azkv"
	contract "github.com/cxyzhang0/wallet-go/ethtran_azkv/contract1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	ubi "gitlab.com/Blockdaemon/ubiquity/ubiquity-go-client/v1/pkg/client"
	"math/big"
	"strconv"
	"testing"
)

// TestBuildTx
// It uses Blockdaemon ubiquity API to get the nonce, and gas related parameters.
// address1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
// address2: 0x4A2EBB506da083caC4d61f9305dF8967E595D16b
func TestBuildTx(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}
	to := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}
	req := tran.TxReq{
		From: from,
		To:   to,
		//ToAddr: "0x4A2EBB506da083caC4d61f9305dF8967E595D16b",
		Amount: big.NewInt(1e15),
		Gas:    21000, // for standard tx transfering ether from one address to another
		// the following will be set be code
		Nonce:     0,
		GasTipCap: big.NewInt(3200),
		GasFeeCap: big.NewInt(32000),
	}

	_, fromAddr, err := tran.GetAddressPubKey(req.From, _sdk)
	if err != nil {
		t.Fatalf("failed to get address pubkey: %+v", err)
	}

	balances, _, err := ubiAPIClient.AccountsAPI.GetListOfBalancesByAddress(ubiCtx, ubiPlatform, ubiNetwork, fromAddr).Execute()
	if err != nil {
		t.Fatalf("failed to get balance via ubiquity: %+v", err)
	}
	if len(balances) == 0 {
		t.Fatalf("no balance record")
	}
	balance, err := tran.StringToBigInt(balances[0].GetConfirmedBalance())
	if err != nil {
		t.Fatalf("failed to pars big int: %+v", err)
	}

	if balance.Cmp(req.Amount) != 1 {
		t.Fatalf("balance %+v <= req amount %+v", balance, req.Amount)
	}

	req.Nonce = uint64(balances[0].GetConfirmedNonce())

	estimate, _, err := ubiAPIClient.TransactionsAPI.FeeEstimate(ubiCtx, ubiPlatform, ubiNetwork).Execute()
	if err != nil {
		t.Fatalf("failed to get fee estimate via ubiquity: %+v", err)
	}

	mediumMap := estimate.EstimatedFees.Medium.(map[string]interface{})
	gasTipCapFloat, ok := mediumMap["max_priority_fee"]
	if !ok {
		t.Fatalf("max_priority_fee is missing")
	}
	gasTipCap, err := tran.StringToBigInt(strconv.FormatFloat(gasTipCapFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Fatalf("failed to convert gasTipCapFloat %s to big int", gasTipCapFloat)
	}
	req.GasTipCap = gasTipCap

	maxTotalFeeFloat, ok := mediumMap["max_total_fee"]
	if !ok {
		t.Fatalf("max_total_fee is missing")
	}
	maxTotalFee, err := tran.StringToBigInt(strconv.FormatFloat(maxTotalFeeFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Fatalf("failed to convert maxTotalFeeFloat %s to big int", maxTotalFeeFloat)
	}
	req.GasFeeCap = maxTotalFee

	rawSignedTx, txHash, err := tran.BuildTx(req, _sdk, chainConfig)
	if err != nil {
		t.Fatalf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", rawSignedTx, txHash)

	receipt, httpResponse, err := ubiAPIClient.TransactionsAPI.TxSend(ubiCtx, ubiPlatform, ubiNetwork).SignedTx(ubi.SignedTx{Tx: rawSignedTx}).Execute()
	if err != nil {
		t.Fatalf("failed to send tx %s \ntx hash %s\n %+v", rawSignedTx, txHash, err)
	}

	t.Logf("receipt: %+v \nhttp response: %+v", receipt, httpResponse)
}

// TestBuildTx_Cancel
// It uses Blockdaemon ubiquity API to get the nonce, and gas related parameters.
// address1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
// address2: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
// address3: 0x4357fB73aF4359D2ec2dc449B90D73495F7794DD
// send to itself with 0 amount and higher fee
// ref: https://info.etherscan.com/how-to-cancel-ethereum-pending-transactions/
// NOTE: could not cancel
func TestBuildTx_Cancel(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}
	to := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}
	req := tran.TxReq{
		From: from,
		To:   to,
		//ToAddr: "0x4A2EBB506da083caC4d61f9305dF8967E595D16b",
		Amount: big.NewInt(0),
		Gas:    210000000, // for standard tx transfering ether from one address to another
		// the following will be set be code
		Nonce:     0,
		GasTipCap: big.NewInt(32000000),
		GasFeeCap: big.NewInt(32000000),
	}

	_, fromAddr, err := tran.GetAddressPubKey(req.From, _sdk)
	if err != nil {
		t.Fatalf("failed to get address pubkey: %+v", err)
	}

	balances, _, err := ubiAPIClient.AccountsAPI.GetListOfBalancesByAddress(ubiCtx, ubiPlatform, ubiNetwork, fromAddr).Execute()
	if err != nil {
		t.Fatalf("failed to get balance via ubiquity: %+v", err)
	}
	if len(balances) == 0 {
		t.Fatalf("no balance record")
	}
	balance, err := tran.StringToBigInt(balances[0].GetConfirmedBalance())
	if err != nil {
		t.Fatalf("failed to pars big int: %+v", err)
	}

	if balance.Cmp(req.Amount) != 1 {
		t.Fatalf("balance %+v <= req amount %+v", balance, req.Amount)
	}

	req.Nonce = uint64(balances[0].GetConfirmedNonce())

	estimate, _, err := ubiAPIClient.TransactionsAPI.FeeEstimate(ubiCtx, ubiPlatform, ubiNetwork).Execute()
	if err != nil {
		t.Fatalf("failed to get fee estimate via ubiquity: %+v", err)
	}

	mediumMap := estimate.EstimatedFees.Medium.(map[string]interface{})
	gasTipCapFloat, ok := mediumMap["max_priority_fee"]
	if !ok {
		t.Fatalf("max_priority_fee is missing")
	}
	gasTipCap, err := tran.StringToBigInt(strconv.FormatFloat(gasTipCapFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Fatalf("failed to convert gasTipCapFloat %s to big int", gasTipCapFloat)
	}
	req.GasTipCap = gasTipCap

	maxTotalFeeFloat, ok := mediumMap["max_total_fee"]
	if !ok {
		t.Fatalf("max_total_fee is missing")
	}
	maxTotalFee, err := tran.StringToBigInt(strconv.FormatFloat(maxTotalFeeFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Fatalf("failed to convert maxTotalFeeFloat %s to big int", maxTotalFeeFloat)
	}
	req.GasFeeCap = maxTotalFee

	rawSignedTx, txHash, err := tran.BuildTx(req, _sdk, chainConfig)
	if err != nil {
		t.Fatalf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", rawSignedTx, txHash)

	receipt, httpResponse, err := ubiAPIClient.TransactionsAPI.TxSend(ubiCtx, ubiPlatform, ubiNetwork).SignedTx(ubi.SignedTx{Tx: rawSignedTx}).Execute()
	if err != nil {
		t.Fatalf("failed to send tx %s \ntx hash %s\n %+v", rawSignedTx, txHash, err)
	}

	t.Logf("receipt: %+v \nhttp response: %+v", receipt, httpResponse)
}

// TestBuildTxToContractAddress - fund the multisig contract
// Goerli faucet does not support contract address so use this to fund multisig contract
// It uses Blockdaemon ubiquity API to get the nonce, and gas related parameters.
// funding address1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
// funding address2: 0x4A2EBB506da083caC4d61f9305dF8967E595D16b
// funding address3: 0x4357fB73aF4359D2ec2dc449B90D73495F7794DD
// lemon contract address 5: 0x63492Cf60244e7Faf88eb0f54c16F26781fec79D
// lemon contract address 4: 0xb6e471768858c22f15fb9e1Eb19da5c4a4094861
// lemon contract address 3: 0x411aB98BD362570702A73f78a3eaEAe62FcB8e2B
// lemon contract address 2: 0x5825342Ec9880fB2bc75feb41Be62165F40cd254
// lemon contract address 1: 0x6428caCED5eB5A5be79C3bF4A295Ff75b2957f94
// paxos contract address 2: 0xffb6ef7E9D920Fac51eb8F490A03C9BC99ed5d86
// paxos contract address 1: 0xb239A44548ec3813aCbBbe4017AcFfc541505b28
// NOTE: the only material difference from TestBuildTx is the Gas limit.
func TestBuildTxToFundContract(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version: "0179a6204ed7491ea5b27a87b541d5cb",
		Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}
	toAddr := "0x63492Cf60244e7Faf88eb0f54c16F26781fec79D"
	//toAddr := "0xb239A44548ec3813aCbBbe4017AcFfc541505b28"
	//to := kmssdk.KeyLabel{
	//	Key: "secp256k1-hsm-1",
	//	//Version: "cb848fb15e3a40b49bc41cbe957ea438",
	//	Version:   "0179a6204ed7491ea5b27a87b541d5cb",
	//	Algorithm: kmssdk.Secp256k1,
	//}
	req := tran.TxReq{
		From:   from,
		ToAddr: toAddr,
		Amount: big.NewInt(5e16),
		Gas:    250000, // for standard tx transfering ether from one address to multisig contract
		// the following will be set be code
		Nonce:     0,
		GasTipCap: big.NewInt(32000),
		GasFeeCap: big.NewInt(320000),
	}

	_, fromAddr, err := tran.GetAddressPubKey(req.From, _sdk)
	if err != nil {
		t.Fatalf("failed to get address pubkey: %+v", err)
	}

	balances, _, err := ubiAPIClient.AccountsAPI.GetListOfBalancesByAddress(ubiCtx, ubiPlatform, ubiNetwork, fromAddr).Execute()
	if err != nil {
		t.Fatalf("failed to get balance via ubiquity: %+v", err)
	}
	if len(balances) == 0 {
		t.Fatalf("no balance record")
	}
	balance, err := tran.StringToBigInt(balances[0].GetConfirmedBalance())
	if err != nil {
		t.Fatalf("failed to pars big int: %+v", err)
	}

	if balance.Cmp(req.Amount) != 1 {
		t.Fatalf("balance %+v <= req amount %+v", balance, req.Amount)
	}

	req.Nonce = uint64(balances[0].GetConfirmedNonce())

	estimate, _, err := ubiAPIClient.TransactionsAPI.FeeEstimate(ubiCtx, ubiPlatform, ubiNetwork).Execute()
	if err != nil {
		t.Fatalf("failed to get fee estimate via ubiquity: %+v", err)
	}

	mediumMap := estimate.EstimatedFees.Medium.(map[string]interface{})
	gasTipCapFloat, ok := mediumMap["max_priority_fee"]
	if !ok {
		t.Fatalf("max_priority_fee is missing")
	}
	gasTipCap, err := tran.StringToBigInt(strconv.FormatFloat(gasTipCapFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Fatalf("failed to convert gasTipCapFloat %s to big int", gasTipCapFloat)
	}
	req.GasTipCap = gasTipCap

	maxTotalFeeFloat, ok := mediumMap["max_total_fee"]
	if !ok {
		t.Fatalf("max_total_fee is missing")
	}
	maxTotalFee, err := tran.StringToBigInt(strconv.FormatFloat(maxTotalFeeFloat.(float64), 'e', -1, 32))
	if err != nil {
		t.Fatalf("failed to convert maxTotalFeeFloat %s to big int", maxTotalFeeFloat)
	}
	req.GasFeeCap = maxTotalFee

	rawSignedTx, txHash, err := tran.BuildTx(req, _sdk, chainConfig)
	if err != nil {
		t.Fatalf("failed to build tx: %+v", err)
	}

	t.Logf("signed tx: %s \ntx hash: %s", rawSignedTx, txHash)

	receipt, httpResponse, err := ubiAPIClient.TransactionsAPI.TxSend(ubiCtx, ubiPlatform, ubiNetwork).SignedTx(ubi.SignedTx{Tx: rawSignedTx}).Execute()
	if err != nil {
		t.Fatalf("failed to send tx %s \ntx hash %s\n %+v", rawSignedTx, txHash, err)
	}

	t.Logf("receipt: %+v \nhttp response: %+v", receipt, httpResponse)
}

/*
---- lemon 5: full with debugging
contract address: 0x63492Cf60244e7Faf88eb0f54c16F26781fec79D
signed tx: 0xf91d890584125028408316e3608080b91d3760806040523480156200001157600080fd5b5060405162001c5738038062001c578339818101604052810190620000379190620005a8565b600a8251111580156200004b575081518311155b8015620000585750600083115b6200009a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000919062000684565b60405180910390fd5b6000805b8351811015620001e1578173ffffffffffffffffffffffffffffffffffffffff16848281518110620000d557620000d4620006a6565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff161162000136576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200012d9062000725565b60405180910390fd5b600160026000868481518110620001525762000151620006a6565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550838181518110620001c157620001c0620006a6565b5b602002602001015191508080620001d89062000776565b9150506200009e565b508260039080519060200190620001fa929190620002d1565b50836001819055507fd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac5647260001b7fb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e660001b7fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc660001b84307f251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc060001b604051602001620002ab969594939291906200085a565b6040516020818303038152906040528051906020012060048190555050505050620008c7565b8280548282559060005260206000209081019282156200034d579160200282015b828111156200034c5782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190620002f2565b5b5090506200035c919062000360565b5090565b5b808211156200037b57600081600090555060010162000361565b5090565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b620003a88162000393565b8114620003b457600080fd5b50565b600081519050620003c8816200039d565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200041e82620003d3565b810181811067ffffffffffffffff8211171562000440576200043f620003e4565b5b80604052505050565b6000620004556200037f565b905062000463828262000413565b919050565b600067ffffffffffffffff821115620004865762000485620003e4565b5b602082029050602081019050919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620004c9826200049c565b9050919050565b620004db81620004bc565b8114620004e757600080fd5b50565b600081519050620004fb81620004d0565b92915050565b600062000518620005128462000468565b62000449565b905080838252602082019050602084028301858111156200053e576200053d62000497565b5b835b818110156200056b5780620005568882620004ea565b84526020840193505060208101905062000540565b5050509392505050565b600082601f8301126200058d576200058c620003ce565b5b81516200059f84826020860162000501565b91505092915050565b600080600060608486031215620005c457620005c362000389565b5b6000620005d486828701620003b7565b935050602084015167ffffffffffffffff811115620005f857620005f76200038e565b5b620006068682870162000575565b92505060406200061986828701620003b7565b9150509250925092565b600082825260208201905092915050565b7f303c7468726573686f6c643c6f776e6572732e6c656e67746800000000000000600082015250565b60006200066c60198362000623565b9150620006798262000634565b602082019050919050565b600060208201905081810360008301526200069f816200065d565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f7265706561746564206f776e6572206f72206e6f7420736f7274656400000000600082015250565b60006200070d601c8362000623565b91506200071a82620006d5565b602082019050919050565b600060208201905081810360008301526200074081620006fe565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000620007838262000393565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203620007b857620007b762000747565b5b600182019050919050565b6000819050919050565b620007d881620007c3565b82525050565b620007e98162000393565b82525050565b6000819050919050565b60006200081a620008146200080e846200049c565b620007ef565b6200049c565b9050919050565b60006200082e82620007f9565b9050919050565b6000620008428262000821565b9050919050565b620008548162000835565b82525050565b600060c082019050620008716000830189620007cd565b620008806020830188620007cd565b6200088f6040830187620007cd565b6200089e6060830186620007de565b620008ad608083018562000849565b620008bc60a0830184620007cd565b979650505050505050565b61138080620008d76000396000f3fe6080604052600436106100745760003560e01c8063aa5df9e21161004e578063aa5df9e2146100ff578063affed0e01461013c578063ca7541ee14610167578063f87c78c7146101925761007b565b80630d8e6e2c1461008057806342cde4e8146100ab578063a0ab9653146100d65761007b565b3661007b57005b600080fd5b34801561008c57600080fd5b506100956101d1565b6040516100a2919061083c565b60405180910390f35b3480156100b757600080fd5b506100c061020e565b6040516100cd9190610877565b60405180910390f35b3480156100e257600080fd5b506100fd60048036038101906100f89190610c5f565b610214565b005b34801561010b57600080fd5b5061012660048036038101906101219190610d85565b6106ac565b6040516101339190610dc1565b60405180910390f35b34801561014857600080fd5b506101516106eb565b60405161015e9190610877565b60405180910390f35b34801561017357600080fd5b5061017c6106f1565b6040516101899190610877565b60405180910390f35b34801561019e57600080fd5b506101b960048036038101906101b49190610ddc565b6106fe565b6040516101c893929190610e82565b60405180910390f35b60606040518060400160405280600481526020017f322e333300000000000000000000000000000000000000000000000000000000815250905090565b60015481565b7f1f3748a2491ab38f2844a84540e798f06e240ec9764f4034a4b85d9de95a309760405160405180910390a1600154875114610285576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161027c90610f05565b60405180910390fd5b85518751148015610297575087518751145b6102d6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102cd90610f71565b60405180910390fd5b7f7b3f83512e4134c9157a582e9b708d7b8535a483ffdac94c37ecebb8a3b63c04336040516103059190610dc1565b60405180910390a13373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614806103735750600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b6103b2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103a990610fdd565b60405180910390fd5b60008060006103c488888888886106fe565b9250925092507f672ea124619314d71be6f340ecb56da6038c8d6b8ffca9bb1be62086d0a73d908383836040516103fd93929190610e82565b60405180910390a16000805b600154811015610637577fbb8c691c28385da4e4d29a158e660fad1c741f388f2170d3c9b67b6d71ab6d12816040516104429190610877565b60405180910390a160006001848f848151811061046257610461610ffd565b5b60200260200101518f858151811061047d5761047c610ffd565b5b60200260200101518f868151811061049857610497610ffd565b5b6020026020010151604051600081526020016040526040516104bd949392919061103b565b6020604051602081039080840390855afa1580156104df573d6000803e3d6000fd5b5050506020604051035190507f464d905a75ac90e0d07b8c2a0cb67371b7f6abd04160c6e323686e3d9beb72b9818360405161051c929190611080565b60405180910390a18273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161180156105a85750600260008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff165b6105e7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105de906110f5565b60405180910390fd5b7f4e57899e25b61543bc91679ea2a1d9edf6409fad79e539a73eab5f03c06d77cd8183604051610618929190611080565b60405180910390a180925050808061062f90611144565b915050610409565b506001600054610647919061118c565b60008190555060008080895160208b018c8e8bf190508061069d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106949061120c565b60405180910390fd5b50505050505050505050505050565b600381815481106106bc57600080fd5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60005481565b6000600380549050905090565b6000806000807f3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d760001b898989805190602001206000548a8a60405160200161074d979695949392919061122c565b60405160208183030381529060405280519060200120905060006004548260405160200161077c929190611313565b60405160208183030381529060405280519060200120905060045482829450945094505050955095509592505050565b600081519050919050565b600082825260208201905092915050565b60005b838110156107e65780820151818401526020810190506107cb565b60008484015250505050565b6000601f19601f8301169050919050565b600061080e826107ac565b61081881856107b7565b93506108288185602086016107c8565b610831816107f2565b840191505092915050565b600060208201905081810360008301526108568184610803565b905092915050565b6000819050919050565b6108718161085e565b82525050565b600060208201905061088c6000830184610868565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6108e3826107f2565b810181811067ffffffffffffffff82111715610902576109016108ab565b5b80604052505050565b6000610915610892565b905061092182826108da565b919050565b600067ffffffffffffffff821115610941576109406108ab565b5b602082029050602081019050919050565b600080fd5b600060ff82169050919050565b61096d81610957565b811461097857600080fd5b50565b60008135905061098a81610964565b92915050565b60006109a361099e84610926565b61090b565b905080838252602082019050602084028301858111156109c6576109c5610952565b5b835b818110156109ef57806109db888261097b565b8452602084019350506020810190506109c8565b5050509392505050565b600082601f830112610a0e57610a0d6108a6565b5b8135610a1e848260208601610990565b91505092915050565b600067ffffffffffffffff821115610a4257610a416108ab565b5b602082029050602081019050919050565b6000819050919050565b610a6681610a53565b8114610a7157600080fd5b50565b600081359050610a8381610a5d565b92915050565b6000610a9c610a9784610a27565b61090b565b90508083825260208201905060208402830185811115610abf57610abe610952565b5b835b81811015610ae85780610ad48882610a74565b845260208401935050602081019050610ac1565b5050509392505050565b600082601f830112610b0757610b066108a6565b5b8135610b17848260208601610a89565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610b4b82610b20565b9050919050565b610b5b81610b40565b8114610b6657600080fd5b50565b600081359050610b7881610b52565b92915050565b610b878161085e565b8114610b9257600080fd5b50565b600081359050610ba481610b7e565b92915050565b600080fd5b600067ffffffffffffffff821115610bca57610bc96108ab565b5b610bd3826107f2565b9050602081019050919050565b82818337600083830152505050565b6000610c02610bfd84610baf565b61090b565b905082815260208101848484011115610c1e57610c1d610baa565b5b610c29848285610be0565b509392505050565b600082601f830112610c4657610c456108a6565b5b8135610c56848260208601610bef565b91505092915050565b600080600080600080600080610100898b031215610c8057610c7f61089c565b5b600089013567ffffffffffffffff811115610c9e57610c9d6108a1565b5b610caa8b828c016109f9565b985050602089013567ffffffffffffffff811115610ccb57610cca6108a1565b5b610cd78b828c01610af2565b975050604089013567ffffffffffffffff811115610cf857610cf76108a1565b5b610d048b828c01610af2565b9650506060610d158b828c01610b69565b9550506080610d268b828c01610b95565b94505060a089013567ffffffffffffffff811115610d4757610d466108a1565b5b610d538b828c01610c31565b93505060c0610d648b828c01610b69565b92505060e0610d758b828c01610b95565b9150509295985092959890939650565b600060208284031215610d9b57610d9a61089c565b5b6000610da984828501610b95565b91505092915050565b610dbb81610b40565b82525050565b6000602082019050610dd66000830184610db2565b92915050565b600080600080600060a08688031215610df857610df761089c565b5b6000610e0688828901610b69565b9550506020610e1788828901610b95565b945050604086013567ffffffffffffffff811115610e3857610e376108a1565b5b610e4488828901610c31565b9350506060610e5588828901610b69565b9250506080610e6688828901610b95565b9150509295509295909350565b610e7c81610a53565b82525050565b6000606082019050610e976000830186610e73565b610ea46020830185610e73565b610eb16040830184610e73565b949350505050565b7f6e6f7420657175616c20746f207468726573686f6c6400000000000000000000600082015250565b6000610eef6016836107b7565b9150610efa82610eb9565b602082019050919050565b60006020820190508181036000830152610f1e81610ee2565b9050919050565b7f6c656e677468206e6f74206d6174636800000000000000000000000000000000600082015250565b6000610f5b6010836107b7565b9150610f6682610f25565b602082019050919050565b60006020820190508181036000830152610f8a81610f4e565b9050919050565b7f77726f6e67206578656375746f72000000000000000000000000000000000000600082015250565b6000610fc7600e836107b7565b9150610fd282610f91565b602082019050919050565b60006020820190508181036000830152610ff681610fba565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b61103581610957565b82525050565b60006080820190506110506000830187610e73565b61105d602083018661102c565b61106a6040830185610e73565b6110776060830184610e73565b95945050505050565b60006040820190506110956000830185610db2565b6110a26020830184610868565b9392505050565b7f76657269667920736967206661696c6564000000000000000000000000000000600082015250565b60006110df6011836107b7565b91506110ea826110a9565b602082019050919050565b6000602082019050818103600083015261110e816110d2565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061114f8261085e565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361118157611180611115565b5b600182019050919050565b60006111978261085e565b91506111a28361085e565b92508282019050808211156111ba576111b9611115565b5b92915050565b7f6e6f745f73756363657373000000000000000000000000000000000000000000600082015250565b60006111f6600b836107b7565b9150611201826111c0565b602082019050919050565b60006020820190508181036000830152611225816111e9565b9050919050565b600060e082019050611241600083018a610e73565b61124e6020830189610db2565b61125b6040830188610868565b6112686060830187610e73565b6112756080830186610868565b61128260a0830185610db2565b61128f60c0830184610868565b98975050505050505050565b600081905092915050565b7f1901000000000000000000000000000000000000000000000000000000000000600082015250565b60006112dc60028361129b565b91506112e7826112a6565b600282019050919050565b6000819050919050565b61130d61130882610a53565b6112f2565b82525050565b600061131e826112cf565b915061132a82856112fc565b60208201915061133a82846112fc565b602082019150819050939250505056fea2646970667358221220b3442fcaac621dad85405ad44fc878ce1e249b3dc35c283e5d3260435fa1bf9264736f6c6343000811003300000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000030000000000000000000000004357fb73af4359d2ec2dc449b90d73495f7794dd0000000000000000000000004a2ebb506da083cac4d61f9305df8967e595d16b0000000000000000000000005b85f5666c9494e69a7adb0cce95ada892ab36072da04c6672b70d239532fd2314889f7f29bd9cb95eb22132b09ae2eceb1c4dc83e28a03a53c242b9a46b3cf50ddf3db2f9eb6ef7af739808b43a512f7e7978a136378b
tx hash: 0xd8a46c129375d1f5b422ea257f79288d8042836e5c79dc6b79e2ad124a0df93c
contract: &{ContractCaller:{contract:0xc0001a2f00} ContractTransactor:{contract:0xc0001a2f00} ContractFilterer:{contract:0xc0001a2f00}}
---- lemon 4: debugging - skip sig verification
contract address: 0xb6e471768858c22f15fb9e1Eb19da5c4a4094861
tx hash: 0x2243fddd9cfd5a1e7390875c776f20ce5ed04bf0d8d73f181503168df4a5cf91
contract: &{ContractCaller:{contract:0xc0002bac80} ContractTransactor:{contract:0xc0002bac80} ContractFilterer:{contract:0xc0002bac80}}
---- lemon 3: more debugging events
contract address: 0x411aB98BD362570702A73f78a3eaEAe62FcB8e2B
tx hash: 0x399ae8d30db505d85cf7dd73651518eb5d9cfc9616e88edd329962b998540e18
contract: &{ContractCaller:{contract:0xc000392c80} ContractTransactor:{contract:0xc000392c80} ContractFilterer:{contract:0xc000392c80}}
---- lemon 2: add getHashes
contract address: 0x5825342Ec9880fB2bc75feb41Be62165F40cd254
tx hash: 0x3ca9ed8c6b8e6e2907a3d911fc2c8812b135271da45f4a8828c10e3bb44707bf
contract: &{ContractCaller:{contract:0xc0002b8280} ContractTransactor:{contract:0xc0002b8280} ContractFilterer:{contract:0xc0002b8280}}
---- lemmon 1
contract address: 0x6428caCED5eB5A5be79C3bF4A295Ff75b2957f94
tx hash: 0x8fcf05e91eccc8e96a27402c239e62542302d79fc99fc656222032e80fe1d105
contract: &{ContractCaller:{contract:0xc0002af400} ContractTransactor:{contract:0xc0002af400} ContractFilterer:{contract:0xc0002af400}}
---- paxos 2
contract address: 0xffb6ef7E9D920Fac51eb8F490A03C9BC99ed5d86
tx hash: 0x3cfd185094a2f5fcfbed3231bc7289f933ef30b10f35d660a377a9caf4b49583
contract: &{ContractCaller:{contract:0xc00012bb80} ContractTransactor:{contract:0xc00012bb80} ContractFilterer:{contract:0xc00012bb80}}
---- paxos 1: used toLower() for addresses in comparison
contract address: 0xb239A44548ec3813aCbBbe4017AcFfc541505b28
tx hash: 0x2d1a557eda7e7bd0d5090325db2f6a6d0eb9744058ee9bc9fd3af8aa4c8e44d6
contract: &{ContractCaller:{contract:0xc000210a00} ContractTransactor:{contract:0xc000210a00} ContractFilterer:{contract:0xc000210a00}}
*/
func TestBuildDeployContract(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	keyLabel1 := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	keyLabel2 := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Version: "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	keyLabel3 := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	req := tran.MultisigDeployTxReq{
		From:         from,
		MultisigFrom: []kmssdk.KeyLabel{keyLabel1, keyLabel2, keyLabel3},
		M:            2,
		GasLimit:     1500000, // TODO: how to calc it?
	}

	fromAddress, _, err := tran.GetAddressPubKey(req.From, _sdk)
	if err != nil {
		t.Fatalf("failed to get address pubkey: %+v", err)
	}

	req.FromAddress = fromAddress

	//rawUrl := fmt.Sprintf("%s/%s/%s/native", ubiNativeURL, ubiPlatform, ubiNetwork)
	client, err := ethclient.Dial(quicknodeURL)
	if err != nil {
		t.Fatalf("failed to get ethclient: %+v", err)
	}

	gasPrice, err := client.SuggestGasPrice(ubiCtx)
	if err != nil {
		t.Fatalf("failed to get gas price: %+v", err)
	}
	req.GasPrice = gasPrice

	nonce, err := client.PendingNonceAt(ubiCtx, *fromAddress)
	if err != nil {
		t.Fatalf("failed to get pending nonce: %+v", err)
	}
	req.Nonce = nonce

	contractAddress, signedTx, _contract, err := tran.BuildDeployContractTx(req, _sdk, client, chainConfig.ChainID)
	if err != nil {
		t.Fatalf("failed to build tx: %+v", err)
	}

	raw, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal signed tx: %+v", err)
	}

	t.Logf("contract address: %s\nsigned tx: %s\ntx hash: %s\ncontract: %+v", contractAddress.String(), hexutil.Encode(raw), signedTx.Hash().String(), _contract)
	//t.Logf("signed tx: %s \ntx hash: %s", rawSignedTx, rawSignedTx.Hash().String())

	/*
		receipt, httpResponse, err := ubiAPIClient.TransactionsAPI.TxSend(ubiCtx, ubiPlatform, ubiNetwork).SignedTx(ubi.SignedTx{Tx: rawSignedTx}).Execute()
		if err != nil {
			t.Fatalf("failed to send tx %s \ntx hash %s\n %+v", rawSignedTx, txHash, err)
		}

		t.Logf("receipt: %+v \nhttp response: %+v", receipt, httpResponse)
	*/
}

/*
contract address: 0x63492Cf60244e7Faf88eb0f54c16F26781fec79D
*/
func TestLoadingContract(t *testing.T) {
	address := common.HexToAddress("0x63492Cf60244e7Faf88eb0f54c16F26781fec79D")
	ctx := context.Background()
	client, err := ethclient.Dial(quicknodeURL)
	FailOnErr(t, err, "FonDial")

	instance, err := contract.NewContract(address, client)
	FailOnErr(t, err, "FonNewContract")

	m, err := instance.Threshold(nil)
	FailOnErr(t, err, "FonThreshold")
	t.Logf("m: %d", m.Int64())
	//_ = instance

	nonce, err := instance.Nonce(nil)
	FailOnErr(t, err, "FonNonce")
	t.Logf("nonce: %d", nonce.Int64())

	balance, err := client.BalanceAt(ctx, address, nil)
	FailOnErr(t, err, "FonBalanceAt")
	t.Logf("balance at %s: %d", address.String(), balance.Int64())

	bytecode, err := client.CodeAt(ctx, address, nil)
	FailOnErr(t, err, "FonCodeAt")
	t.Logf("address %s is contract: %t", address.String(), len(bytecode) > 0)
}

/*
executor address 1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
executor address 2: 0x4A2EBB506da083caC4d61f9305dF8967E595D16b
executor address 3: 0x4357fB73aF4359D2ec2dc449B90D73495F7794DD
lemon contract address 5: 0x63492Cf60244e7Faf88eb0f54c16F26781fec79D
lemon contract address 4: 0xb6e471768858c22f15fb9e1Eb19da5c4a4094861
lemon contract address 3: 0x411aB98BD362570702A73f78a3eaEAe62FcB8e2B
lemon contract address 2: 0x5825342Ec9880fB2bc75feb41Be62165F40cd254
lemon contract address 1: 0x6428caCED5eB5A5be79C3bF4A295Ff75b2957f94
paxos contract address 1: 0xb239A44548ec3813aCbBbe4017AcFfc541505b28
paxos contract address 2: 0xffb6ef7E9D920Fac51eb8F490A03C9BC99ed5d86
to address: 0x4A2EBB506da083caC4d61f9305dF8967E595D16b
tx hash:
*/
func TestBuildMultisigTx(t *testing.T) {
	executor := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Version: "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}
	executorAddress, _, err := tran.GetAddressPubKey(executor, _sdk)
	if err != nil {
		t.Fatalf("failed to get executor address pub key: %+v", err)
	}

	keyLabel1 := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	keyLabel2 := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Version: "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	keyLabel3 := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}

	address1, _, err := tran.GetAddressPubKey(keyLabel1, _sdk)
	if err != nil {
		t.Fatalf("failed to get address pub key 1: %+v", err)
	}
	addressInfo1 := tran.AddressInfo{KeyLabel: keyLabel1, Address: *address1}

	address2, _, err := tran.GetAddressPubKey(keyLabel2, _sdk)
	if err != nil {
		t.Fatalf("failed to get address pub key 2: %+v", err)
	}
	addressInfo2 := tran.AddressInfo{KeyLabel: keyLabel2, Address: *address2}

	address3, _, err := tran.GetAddressPubKey(keyLabel3, _sdk)
	if err != nil {
		t.Fatalf("failed to get address pub key 3: %+v", err)
	}
	addressInfo3 := tran.AddressInfo{KeyLabel: keyLabel3, Address: *address3}

	contractAddress := common.HexToAddress("0x63492Cf60244e7Faf88eb0f54c16F26781fec79D")

	to := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}
	toAddress, _, err := tran.GetAddressPubKey(to, _sdk)
	if err != nil {
		t.Fatalf("failed to get to address pub key: %+v", err)
	}

	req := tran.MultisigTxReq{
		Executor:            executor,
		ExecutorAddress:     executorAddress,
		MultisigAddressInfo: []*tran.AddressInfo{&addressInfo1, &addressInfo2, &addressInfo3},
		ContractAddress:     &contractAddress,
		M:                   2,
		To:                  to,
		ToAddress:           toAddress,
		Amount:              big.NewInt(1e14),
		GasLimit:            2500000, // TODO: how to calc it?
		Data:                []byte(""),
	}

	//rawUrl := fmt.Sprintf("%s/%s/%s/native", ubiNativeURL, ubiPlatform, ubiNetwork)
	client, err := ethclient.Dial(quicknodeURL)
	if err != nil {
		t.Fatalf("failed to get ethclient: %+v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatalf("failed to get gas price: %+v", err)
	}
	req.GasPrice = gasPrice

	nonce, err := client.PendingNonceAt(context.Background(), contractAddress)
	if err != nil {
		t.Fatalf("failed to get contract pending nonce: %+v", err)
	}
	req.ContractNonce = nonce

	nonce, err = client.PendingNonceAt(context.Background(), *executorAddress)
	if err != nil {
		t.Fatalf("failed to get executor pending nonce: %+v", err)
	}
	req.ExecutorNonce = nonce

	instance, err := contract.NewContract(contractAddress, client)
	if err != nil {
		t.Fatalf("failed to have new contract instance: %+v", err)
	}

	variableNonce, err := instance.Nonce(nil)
	if err != nil {
		t.Fatalf("failed to get variable nonce: %+v", err)
	}
	req.ContractVariableNonce = variableNonce.Uint64()

	signedTx, err := tran.BuildMultisigTx(req, _sdk, client, chainConfig.ChainID)
	if err != nil {
		t.Fatalf("failed to build tx: %+v", err)
	}

	raw, err := signedTx.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal signed tx: %+v", err)
	}

	t.Logf("contract address: %s\nsigned tx: %s\ntx hash: %s", contractAddress.String(), hexutil.Encode(raw), signedTx.Hash().String())
	/*
		{ // 完了检查确实转账成功
			time.Sleep(time.Minute)
			client, err := ethclient.Dial(quicknodeURL)
			if err != nil {
				t.Fatalf("failed to dial: %+v", err)
			}

			bal, err := client.BalanceAt(context.Background(), contractAddress, nil)
			if err != nil {
				t.Fatalf("failed to get balance: %+v", err)
			}

			t.Logf("new balance: %d", bal.Int64())
		}
		time.Sleep(time.Minute)
	*/
	/*
		receipt, httpResponse, err := ubiAPIClient.TransactionsAPI.TxSend(ubiCtx, ubiPlatform, ubiNetwork).SignedTx(ubi.SignedTx{Tx: rawSignedTx}).Execute()
		if err != nil {
			t.Fatalf("failed to send tx %s \ntx hash %s\n %+v", rawSignedTx, txHash, err)
		}

		t.Logf("receipt: %+v \nhttp response: %+v", receipt, httpResponse)
	*/
}
