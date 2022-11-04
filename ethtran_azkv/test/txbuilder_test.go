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
// lemon contract address 3: 0x411aB98BD362570702A73f78a3eaEAe62FcB8e2B
// lemon contract address 2: 0x5825342Ec9880fB2bc75feb41Be62165F40cd254
// lemon contract address 1: 0x6428caCED5eB5A5be79C3bF4A295Ff75b2957f94
// paxos contract address 2: 0xffb6ef7E9D920Fac51eb8F490A03C9BC99ed5d86
// paxos contract address 1: 0xb239A44548ec3813aCbBbe4017AcFfc541505b28
// NOTE: the only material difference from TestBuildTx is the Gas limit.
func TestBuildTxToFundContract(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version: "0179a6204ed7491ea5b27a87b541d5cb",
		//Version:   "b6aec266b6a147f7a1c40fe842504650",
		Algorithm: kmssdk.Secp256k1,
	}
	toAddr := "0x411aB98BD362570702A73f78a3eaEAe62FcB8e2B"
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
---- lemon 4: debugging - skip sig verification
contract address: 0xb6e471768858c22f15fb9e1Eb19da5c4a4094861
signed tx: 0xf91c5b038508d86e94ff8316e3608080b91c0860806040523480156200001157600080fd5b5060405162001b2838038062001b288339818101604052810190620000379190620005a8565b600a8251111580156200004b575081518311155b8015620000585750600083115b6200009a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000919062000684565b60405180910390fd5b6000805b8351811015620001e1578173ffffffffffffffffffffffffffffffffffffffff16848281518110620000d557620000d4620006a6565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff161162000136576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200012d9062000725565b60405180910390fd5b600160026000868481518110620001525762000151620006a6565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550838181518110620001c157620001c0620006a6565b5b602002602001015191508080620001d89062000776565b9150506200009e565b508260039080519060200190620001fa929190620002d1565b50836001819055507fd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac5647260001b7fb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e660001b7fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc660001b84307f251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc060001b604051602001620002ab969594939291906200085a565b6040516020818303038152906040528051906020012060048190555050505050620008c7565b8280548282559060005260206000209081019282156200034d579160200282015b828111156200034c5782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190620002f2565b5b5090506200035c919062000360565b5090565b5b808211156200037b57600081600090555060010162000361565b5090565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b620003a88162000393565b8114620003b457600080fd5b50565b600081519050620003c8816200039d565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200041e82620003d3565b810181811067ffffffffffffffff8211171562000440576200043f620003e4565b5b80604052505050565b6000620004556200037f565b905062000463828262000413565b919050565b600067ffffffffffffffff821115620004865762000485620003e4565b5b602082029050602081019050919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620004c9826200049c565b9050919050565b620004db81620004bc565b8114620004e757600080fd5b50565b600081519050620004fb81620004d0565b92915050565b600062000518620005128462000468565b62000449565b905080838252602082019050602084028301858111156200053e576200053d62000497565b5b835b818110156200056b5780620005568882620004ea565b84526020840193505060208101905062000540565b5050509392505050565b600082601f8301126200058d576200058c620003ce565b5b81516200059f84826020860162000501565b91505092915050565b600080600060608486031215620005c457620005c362000389565b5b6000620005d486828701620003b7565b935050602084015167ffffffffffffffff811115620005f857620005f76200038e565b5b620006068682870162000575565b92505060406200061986828701620003b7565b9150509250925092565b600082825260208201905092915050565b7f303c7468726573686f6c643c6f776e6572732e6c656e67746800000000000000600082015250565b60006200066c60198362000623565b9150620006798262000634565b602082019050919050565b600060208201905081810360008301526200069f816200065d565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f7265706561746564206f776e6572206f72206e6f7420736f7274656400000000600082015250565b60006200070d601c8362000623565b91506200071a82620006d5565b602082019050919050565b600060208201905081810360008301526200074081620006fe565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000620007838262000393565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203620007b857620007b762000747565b5b600182019050919050565b6000819050919050565b620007d881620007c3565b82525050565b620007e98162000393565b82525050565b6000819050919050565b60006200081a620008146200080e846200049c565b620007ef565b6200049c565b9050919050565b60006200082e82620007f9565b9050919050565b6000620008428262000821565b9050919050565b620008548162000835565b82525050565b600060c082019050620008716000830189620007cd565b620008806020830188620007cd565b6200088f6040830187620007cd565b6200089e6060830186620007de565b620008ad608083018562000849565b620008bc60a0830184620007cd565b979650505050505050565b61125180620008d76000396000f3fe6080604052600436106100745760003560e01c8063aa5df9e21161004e578063aa5df9e2146100ff578063affed0e01461013c578063ca7541ee14610167578063f87c78c7146101925761007b565b80630d8e6e2c1461008057806342cde4e8146100ab578063a0ab9653146100d65761007b565b3661007b57005b600080fd5b34801561008c57600080fd5b506100956101d1565b6040516100a29190610779565b60405180910390f35b3480156100b757600080fd5b506100c061020e565b6040516100cd91906107b4565b60405180910390f35b3480156100e257600080fd5b506100fd60048036038101906100f89190610b9c565b610214565b005b34801561010b57600080fd5b5061012660048036038101906101219190610cc2565b6105e9565b6040516101339190610cfe565b60405180910390f35b34801561014857600080fd5b50610151610628565b60405161015e91906107b4565b60405180910390f35b34801561017357600080fd5b5061017c61062e565b60405161018991906107b4565b60405180910390f35b34801561019e57600080fd5b506101b960048036038101906101b49190610d19565b61063b565b6040516101c893929190610dbf565b60405180910390f35b60606040518060400160405280600481526020017f322e333300000000000000000000000000000000000000000000000000000000815250905090565b60015481565b7f1f3748a2491ab38f2844a84540e798f06e240ec9764f4034a4b85d9de95a309760405160405180910390a1600154875114610285576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161027c90610e42565b60405180910390fd5b85518751148015610297575087518751145b6102d6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102cd90610eae565b60405180910390fd5b7f7b3f83512e4134c9157a582e9b708d7b8535a483ffdac94c37ecebb8a3b63c04336040516103059190610cfe565b60405180910390a13373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614806103735750600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b6103b2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103a990610f1a565b60405180910390fd5b60008060006103c4888888888861063b565b9250925092507f672ea124619314d71be6f340ecb56da6038c8d6b8ffca9bb1be62086d0a73d908383836040516103fd93929190610dbf565b60405180910390a16000805b600154811015610574577fbb8c691c28385da4e4d29a158e660fad1c741f388f2170d3c9b67b6d71ab6d128160405161044291906107b4565b60405180910390a160006001848f848151811061046257610461610f3a565b5b60200260200101518f858151811061047d5761047c610f3a565b5b60200260200101518f868151811061049857610497610f3a565b5b6020026020010151604051600081526020016040526040516104bd9493929190610f78565b6020604051602081039080840390855afa1580156104df573d6000803e3d6000fd5b5050506020604051035190507f464d905a75ac90e0d07b8c2a0cb67371b7f6abd04160c6e323686e3d9beb72b9818360405161051c929190610fbd565b60405180910390a17f4e57899e25b61543bc91679ea2a1d9edf6409fad79e539a73eab5f03c06d77cd8183604051610555929190610fbd565b60405180910390a180925050808061056c90611015565b915050610409565b506001600054610584919061105d565b60008190555060008080895160208b018c8e8bf19050806105da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105d1906110dd565b60405180910390fd5b50505050505050505050505050565b600381815481106105f957600080fd5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60005481565b6000600380549050905090565b6000806000807f3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d760001b898989805190602001206000548a8a60405160200161068a97969594939291906110fd565b6040516020818303038152906040528051906020012090506000600454826040516020016106b99291906111e4565b60405160208183030381529060405280519060200120905060045482829450945094505050955095509592505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610723578082015181840152602081019050610708565b60008484015250505050565b6000601f19601f8301169050919050565b600061074b826106e9565b61075581856106f4565b9350610765818560208601610705565b61076e8161072f565b840191505092915050565b600060208201905081810360008301526107938184610740565b905092915050565b6000819050919050565b6107ae8161079b565b82525050565b60006020820190506107c960008301846107a5565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6108208261072f565b810181811067ffffffffffffffff8211171561083f5761083e6107e8565b5b80604052505050565b60006108526107cf565b905061085e8282610817565b919050565b600067ffffffffffffffff82111561087e5761087d6107e8565b5b602082029050602081019050919050565b600080fd5b600060ff82169050919050565b6108aa81610894565b81146108b557600080fd5b50565b6000813590506108c7816108a1565b92915050565b60006108e06108db84610863565b610848565b905080838252602082019050602084028301858111156109035761090261088f565b5b835b8181101561092c578061091888826108b8565b845260208401935050602081019050610905565b5050509392505050565b600082601f83011261094b5761094a6107e3565b5b813561095b8482602086016108cd565b91505092915050565b600067ffffffffffffffff82111561097f5761097e6107e8565b5b602082029050602081019050919050565b6000819050919050565b6109a381610990565b81146109ae57600080fd5b50565b6000813590506109c08161099a565b92915050565b60006109d96109d484610964565b610848565b905080838252602082019050602084028301858111156109fc576109fb61088f565b5b835b81811015610a255780610a1188826109b1565b8452602084019350506020810190506109fe565b5050509392505050565b600082601f830112610a4457610a436107e3565b5b8135610a548482602086016109c6565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610a8882610a5d565b9050919050565b610a9881610a7d565b8114610aa357600080fd5b50565b600081359050610ab581610a8f565b92915050565b610ac48161079b565b8114610acf57600080fd5b50565b600081359050610ae181610abb565b92915050565b600080fd5b600067ffffffffffffffff821115610b0757610b066107e8565b5b610b108261072f565b9050602081019050919050565b82818337600083830152505050565b6000610b3f610b3a84610aec565b610848565b905082815260208101848484011115610b5b57610b5a610ae7565b5b610b66848285610b1d565b509392505050565b600082601f830112610b8357610b826107e3565b5b8135610b93848260208601610b2c565b91505092915050565b600080600080600080600080610100898b031215610bbd57610bbc6107d9565b5b600089013567ffffffffffffffff811115610bdb57610bda6107de565b5b610be78b828c01610936565b985050602089013567ffffffffffffffff811115610c0857610c076107de565b5b610c148b828c01610a2f565b975050604089013567ffffffffffffffff811115610c3557610c346107de565b5b610c418b828c01610a2f565b9650506060610c528b828c01610aa6565b9550506080610c638b828c01610ad2565b94505060a089013567ffffffffffffffff811115610c8457610c836107de565b5b610c908b828c01610b6e565b93505060c0610ca18b828c01610aa6565b92505060e0610cb28b828c01610ad2565b9150509295985092959890939650565b600060208284031215610cd857610cd76107d9565b5b6000610ce684828501610ad2565b91505092915050565b610cf881610a7d565b82525050565b6000602082019050610d136000830184610cef565b92915050565b600080600080600060a08688031215610d3557610d346107d9565b5b6000610d4388828901610aa6565b9550506020610d5488828901610ad2565b945050604086013567ffffffffffffffff811115610d7557610d746107de565b5b610d8188828901610b6e565b9350506060610d9288828901610aa6565b9250506080610da388828901610ad2565b9150509295509295909350565b610db981610990565b82525050565b6000606082019050610dd46000830186610db0565b610de16020830185610db0565b610dee6040830184610db0565b949350505050565b7f6e6f7420657175616c20746f207468726573686f6c6400000000000000000000600082015250565b6000610e2c6016836106f4565b9150610e3782610df6565b602082019050919050565b60006020820190508181036000830152610e5b81610e1f565b9050919050565b7f6c656e677468206e6f74206d6174636800000000000000000000000000000000600082015250565b6000610e986010836106f4565b9150610ea382610e62565b602082019050919050565b60006020820190508181036000830152610ec781610e8b565b9050919050565b7f77726f6e67206578656375746f72000000000000000000000000000000000000600082015250565b6000610f04600e836106f4565b9150610f0f82610ece565b602082019050919050565b60006020820190508181036000830152610f3381610ef7565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b610f7281610894565b82525050565b6000608082019050610f8d6000830187610db0565b610f9a6020830186610f69565b610fa76040830185610db0565b610fb46060830184610db0565b95945050505050565b6000604082019050610fd26000830185610cef565b610fdf60208301846107a5565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006110208261079b565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361105257611051610fe6565b5b600182019050919050565b60006110688261079b565b91506110738361079b565b925082820190508082111561108b5761108a610fe6565b5b92915050565b7f6e6f745f73756363657373000000000000000000000000000000000000000000600082015250565b60006110c7600b836106f4565b91506110d282611091565b602082019050919050565b600060208201905081810360008301526110f6816110ba565b9050919050565b600060e082019050611112600083018a610db0565b61111f6020830189610cef565b61112c60408301886107a5565b6111396060830187610db0565b61114660808301866107a5565b61115360a0830185610cef565b61116060c08301846107a5565b98975050505050505050565b600081905092915050565b7f1901000000000000000000000000000000000000000000000000000000000000600082015250565b60006111ad60028361116c565b91506111b882611177565b600282019050919050565b6000819050919050565b6111de6111d982610990565b6111c3565b82525050565b60006111ef826111a0565b91506111fb82856111cd565b60208201915061120b82846111cd565b602082019150819050939250505056fea26469706673582212204887eb57ea7f040878de6d126780e3bbc08471aded194d4ef5414fb1d72847de64736f6c6343000811003300000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000030000000000000000000000004357fb73af4359d2ec2dc449b90d73495f7794dd0000000000000000000000004a2ebb506da083cac4d61f9305df8967e595d16b0000000000000000000000005b85f5666c9494e69a7adb0cce95ada892ab36072ea0d1b5a667ec0018be20c8c37da8a22e0a9fcf79022c6edd2d08fe4684f62d6afea05d2bdb4e5a8d63cda199614393608233d85c828e0f25e836368b5ff5a80c6a9d
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

	contractAddress, signedTx, _contract, err := tran.BuildDeployContractTx(req, _sdk, client, chainConfig)
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
contract address: 0xffb6ef7E9D920Fac51eb8F490A03C9BC99ed5d86
previous contract address: 0xb239A44548ec3813aCbBbe4017AcFfc541505b28
*/
func TestLoadingContract(t *testing.T) {
	address := common.HexToAddress("0x5825342Ec9880fB2bc75feb41Be62165F40cd254")
	client, err := ethclient.Dial(quicknodeURL)
	if err != nil {
		t.Fatalf("failed to get ethclient: %+v", err)
	}

	instance, err := contract.NewContract(address, client)
	if err != nil {
		t.Fatalf("failed to have new contract instance: %+v", err)
	}

	m, err := instance.Threshold(nil)
	if err != nil {
		t.Fatalf("failed to get threshold: %+v", err)
	}
	t.Logf("m: %d", m.Int64())
	//_ = instance

	nonce, err := instance.Nonce(nil)
	if err != nil {
		t.Fatalf("failed to get nonce: %+v", err)
	}
	t.Logf("nonce: %d", nonce.Int64())
}

/*
executor address 1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
executor address 2: 0x4A2EBB506da083caC4d61f9305dF8967E595D16b
executor address 3: 0x4357fB73aF4359D2ec2dc449B90D73495F7794DD
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

	contractAddress := common.HexToAddress("0x411aB98BD362570702A73f78a3eaEAe62FcB8e2B")

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

	signedTx, err := tran.BuildMultisigTx(req, _sdk, client, chainConfig)
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
