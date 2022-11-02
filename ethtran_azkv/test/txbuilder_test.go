package test

import (
	"context"
	kmssdk "github.com/cxyzhang0/wallet-go/azkv/sdk"
	tran "github.com/cxyzhang0/wallet-go/ethtran_azkv"
	"github.com/cxyzhang0/wallet-go/ethtran_azkv/contract"
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

// TestBuildTxToContractAddress - fund the multisig contract
// Goerli faucet does not support contract address so use this to fund multisig contract
// It uses Blockdaemon ubiquity API to get the nonce, and gas related parameters.
// address1: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
// address2: 0xb239A44548ec3813aCbBbe4017AcFfc541505b28
// NOTE: the only material difference from TestBuildTx is the Gas limit.
func TestBuildTxToFundContract(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
		Algorithm: kmssdk.Secp256k1,
	}
	toAddr := "0xb239A44548ec3813aCbBbe4017AcFfc541505b28"
	//to := kmssdk.KeyLabel{
	//	Key: "secp256k1-hsm-1",
	//	//Version: "cb848fb15e3a40b49bc41cbe957ea438",
	//	Version:   "0179a6204ed7491ea5b27a87b541d5cb",
	//	Algorithm: kmssdk.Secp256k1,
	//}
	req := tran.TxReq{
		From:   from,
		ToAddr: toAddr,
		Amount: big.NewInt(1e15),
		Gas:    25000, // for standard tx transfering ether from one address to multisig contract
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

/*
contract address: 0xb239A44548ec3813aCbBbe4017AcFfc541505b28
signed tx: 0xf91cad0a83047835832dc6c08080b91c5c60a06040523480156200001157600080fd5b5060405162001b7c38038062001b7c833981810160405281019062000037919062000610565b6200004983836200011860201b60201c565b7fd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac5647260001b7fb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e660001b7fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc660001b83307f251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc060001b604051602001620000f29695949392919062000722565b60405160208183030381529060405280519060200120608081815250505050506200083a565b60148151111580156200012c575080518211155b8015620001395750600082115b6200014357600080fd5b60005b6003805490508110156200020357600060026000600384815481106200017157620001706200078f565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508080620001fa90620007ed565b91505062000146565b506000805b825181101562000313578173ffffffffffffffffffffffffffffffffffffffff168382815181106200023f576200023e6200078f565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff16116200026857600080fd5b6001600260008584815181106200028457620002836200078f565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550828181518110620002f357620002f26200078f565b5b6020026020010151915080806200030a90620007ed565b91505062000208565b5081600390805190602001906200032c92919062000339565b5082600181905550505050565b828054828255906000526020600020908101928215620003b5579160200282015b82811115620003b45782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550916020019190600101906200035a565b5b509050620003c49190620003c8565b5090565b5b80821115620003e3576000816000905550600101620003c9565b5090565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b6200041081620003fb565b81146200041c57600080fd5b50565b600081519050620004308162000405565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b62000486826200043b565b810181811067ffffffffffffffff82111715620004a857620004a76200044c565b5b80604052505050565b6000620004bd620003e7565b9050620004cb82826200047b565b919050565b600067ffffffffffffffff821115620004ee57620004ed6200044c565b5b602082029050602081019050919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620005318262000504565b9050919050565b620005438162000524565b81146200054f57600080fd5b50565b600081519050620005638162000538565b92915050565b6000620005806200057a84620004d0565b620004b1565b90508083825260208201905060208402830185811115620005a657620005a5620004ff565b5b835b81811015620005d35780620005be888262000552565b845260208401935050602081019050620005a8565b5050509392505050565b600082601f830112620005f557620005f462000436565b5b81516200060784826020860162000569565b91505092915050565b6000806000606084860312156200062c576200062b620003f1565b5b60006200063c868287016200041f565b935050602084015167ffffffffffffffff81111562000660576200065f620003f6565b5b6200066e86828701620005dd565b925050604062000681868287016200041f565b9150509250925092565b6000819050919050565b620006a0816200068b565b82525050565b620006b181620003fb565b82525050565b6000819050919050565b6000620006e2620006dc620006d68462000504565b620006b7565b62000504565b9050919050565b6000620006f682620006c1565b9050919050565b60006200070a82620006e9565b9050919050565b6200071c81620006fd565b82525050565b600060c08201905062000739600083018962000695565b62000748602083018862000695565b62000757604083018762000695565b620007666060830186620006a6565b62000775608083018562000711565b6200078460a083018462000695565b979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000620007fa82620003fb565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036200082f576200082e620007be565b5b600182019050919050565b60805161132662000856600039600061027b01526113266000f3fe6080604052600436106100595760003560e01c806342cde4e814610065578063a0ab965314610090578063aa5df9e2146100b9578063affe39c1146100f6578063affed0e014610121578063f3182e851461014c57610060565b3661006057005b600080fd5b34801561007157600080fd5b5061007a610175565b6040516100879190610898565b60405180910390f35b34801561009c57600080fd5b506100b760048036038101906100b29190610c91565b61017b565b005b3480156100c557600080fd5b506100e060048036038101906100db9190610db7565b6104b8565b6040516100ed9190610df3565b60405180910390f35b34801561010257600080fd5b5061010b6104f7565b6040516101189190610ecc565b60405180910390f35b34801561012d57600080fd5b50610136610585565b6040516101439190610898565b60405180910390f35b34801561015857600080fd5b50610173600480360381019061016e9190610fb1565b61058b565b005b60015481565b60015487511461018a57600080fd5b8551875114801561019c575087518751145b6101a557600080fd5b3373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16148061020b5750600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b61021457600080fd5b60007f3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d760001b86868680519060200120600054878760405160200161025f979695949392919061101c565b60405160208183030381529060405280519060200120905060007f0000000000000000000000000000000000000000000000000000000000000000826040516020016102ac929190611103565b6040516020818303038152906040528051906020012090506000805b6001548110156104175760006001848e84815181106102ea576102e961113a565b5b60200260200101518e85815181106103055761030461113a565b5b60200260200101518e86815181106103205761031f61113a565b5b6020026020010151604051600081526020016040526040516103459493929190611178565b6020604051602081039080840390855afa158015610367573d6000803e3d6000fd5b5050506020604051035190508273ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161180156103f75750600260008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff165b61040057600080fd5b80925050808061040f906111ec565b9150506102c8565b5060016000546104279190611234565b60008190555060008873ffffffffffffffffffffffffffffffffffffffff168886908960405161045791906112d9565b600060405180830381858888f193505050503d8060008114610495576040519150601f19603f3d011682016040523d82523d6000602084013e61049a565b606091505b505080915050806104aa57600080fd5b505050505050505050505050565b600381815481106104c857600080fd5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6060600380548060200260200160405190810160405280929190818152602001828054801561057b57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610531575b5050505050905090565b60005481565b3073ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146105c357600080fd5b6105cd82826105d1565b5050565b60148151111580156105e4575080518211155b80156105f05750600082115b6105f957600080fd5b60005b6003805490508110156106b257600060026000600384815481106106235761062261113a565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555080806106aa906111ec565b9150506105fc565b506000805b82518110156107b4578173ffffffffffffffffffffffffffffffffffffffff168382815181106106ea576106e961113a565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff161161071257600080fd5b60016002600085848151811061072b5761072a61113a565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508281815181106107975761079661113a565b5b6020026020010151915080806107ac906111ec565b9150506106b7565b5081600390805190602001906107cb9291906107d8565b5082600181905550505050565b828054828255906000526020600020908101928215610851579160200282015b828111156108505782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550916020019190600101906107f8565b5b50905061085e9190610862565b5090565b5b8082111561087b576000816000905550600101610863565b5090565b6000819050919050565b6108928161087f565b82525050565b60006020820190506108ad6000830184610889565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610915826108cc565b810181811067ffffffffffffffff82111715610934576109336108dd565b5b80604052505050565b60006109476108b3565b9050610953828261090c565b919050565b600067ffffffffffffffff821115610973576109726108dd565b5b602082029050602081019050919050565b600080fd5b600060ff82169050919050565b61099f81610989565b81146109aa57600080fd5b50565b6000813590506109bc81610996565b92915050565b60006109d56109d084610958565b61093d565b905080838252602082019050602084028301858111156109f8576109f7610984565b5b835b81811015610a215780610a0d88826109ad565b8452602084019350506020810190506109fa565b5050509392505050565b600082601f830112610a4057610a3f6108c7565b5b8135610a508482602086016109c2565b91505092915050565b600067ffffffffffffffff821115610a7457610a736108dd565b5b602082029050602081019050919050565b6000819050919050565b610a9881610a85565b8114610aa357600080fd5b50565b600081359050610ab581610a8f565b92915050565b6000610ace610ac984610a59565b61093d565b90508083825260208201905060208402830185811115610af157610af0610984565b5b835b81811015610b1a5780610b068882610aa6565b845260208401935050602081019050610af3565b5050509392505050565b600082601f830112610b3957610b386108c7565b5b8135610b49848260208601610abb565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610b7d82610b52565b9050919050565b610b8d81610b72565b8114610b9857600080fd5b50565b600081359050610baa81610b84565b92915050565b610bb98161087f565b8114610bc457600080fd5b50565b600081359050610bd681610bb0565b92915050565b600080fd5b600067ffffffffffffffff821115610bfc57610bfb6108dd565b5b610c05826108cc565b9050602081019050919050565b82818337600083830152505050565b6000610c34610c2f84610be1565b61093d565b905082815260208101848484011115610c5057610c4f610bdc565b5b610c5b848285610c12565b509392505050565b600082601f830112610c7857610c776108c7565b5b8135610c88848260208601610c21565b91505092915050565b600080600080600080600080610100898b031215610cb257610cb16108bd565b5b600089013567ffffffffffffffff811115610cd057610ccf6108c2565b5b610cdc8b828c01610a2b565b985050602089013567ffffffffffffffff811115610cfd57610cfc6108c2565b5b610d098b828c01610b24565b975050604089013567ffffffffffffffff811115610d2a57610d296108c2565b5b610d368b828c01610b24565b9650506060610d478b828c01610b9b565b9550506080610d588b828c01610bc7565b94505060a089013567ffffffffffffffff811115610d7957610d786108c2565b5b610d858b828c01610c63565b93505060c0610d968b828c01610b9b565b92505060e0610da78b828c01610bc7565b9150509295985092959890939650565b600060208284031215610dcd57610dcc6108bd565b5b6000610ddb84828501610bc7565b91505092915050565b610ded81610b72565b82525050565b6000602082019050610e086000830184610de4565b92915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b610e4381610b72565b82525050565b6000610e558383610e3a565b60208301905092915050565b6000602082019050919050565b6000610e7982610e0e565b610e838185610e19565b9350610e8e83610e2a565b8060005b83811015610ebf578151610ea68882610e49565b9750610eb183610e61565b925050600181019050610e92565b5085935050505092915050565b60006020820190508181036000830152610ee68184610e6e565b905092915050565b600067ffffffffffffffff821115610f0957610f086108dd565b5b602082029050602081019050919050565b6000610f2d610f2884610eee565b61093d565b90508083825260208201905060208402830185811115610f5057610f4f610984565b5b835b81811015610f795780610f658882610b9b565b845260208401935050602081019050610f52565b5050509392505050565b600082601f830112610f9857610f976108c7565b5b8135610fa8848260208601610f1a565b91505092915050565b60008060408385031215610fc857610fc76108bd565b5b6000610fd685828601610bc7565b925050602083013567ffffffffffffffff811115610ff757610ff66108c2565b5b61100385828601610f83565b9150509250929050565b61101681610a85565b82525050565b600060e082019050611031600083018a61100d565b61103e6020830189610de4565b61104b6040830188610889565b611058606083018761100d565b6110656080830186610889565b61107260a0830185610de4565b61107f60c0830184610889565b98975050505050505050565b600081905092915050565b7f1901000000000000000000000000000000000000000000000000000000000000600082015250565b60006110cc60028361108b565b91506110d782611096565b600282019050919050565b6000819050919050565b6110fd6110f882610a85565b6110e2565b82525050565b600061110e826110bf565b915061111a82856110ec565b60208201915061112a82846110ec565b6020820191508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b61117281610989565b82525050565b600060808201905061118d600083018761100d565b61119a6020830186611169565b6111a7604083018561100d565b6111b4606083018461100d565b95945050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006111f78261087f565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611229576112286111bd565b5b600182019050919050565b600061123f8261087f565b915061124a8361087f565b9250828201905080821115611262576112616111bd565b5b92915050565b600081519050919050565b600081905092915050565b60005b8381101561129c578082015181840152602081019050611281565b60008484015250505050565b60006112b382611268565b6112bd8185611273565b93506112cd81856020860161127e565b80840191505092915050565b60006112e582846112a8565b91508190509291505056fea26469706673582212200015b84f8f3ecac80515c3a2ebfc997a4743ade0de072f015adfcf0835dd9ae164736f6c6343000811003300000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000030000000000000000000000004357fb73af4359d2ec2dc449b90d73495f7794dd0000000000000000000000004a2ebb506da083cac4d61f9305df8967e595d16b0000000000000000000000005b85f5666c9494e69a7adb0cce95ada892ab36072da02512644c0788479baf05efbec5171bcf6a156be82b3baf5bf2dffd3ec8a3dc0da04309aeb3209f15fc4bb2299626d986501e1920cd460238c6f8c5f9854d56202a
tx hash: 0x2d1a557eda7e7bd0d5090325db2f6a6d0eb9744058ee9bc9fd3af8aa4c8e44d6
contract: &{ContractCaller:{contract:0xc000210a00} ContractTransactor:{contract:0xc000210a00} ContractFilterer:{contract:0xc000210a00}}
*/
func TestBuildDeployContract(t *testing.T) {
	from := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
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
contract address: 0xb239A44548ec3813aCbBbe4017AcFfc541505b28
*/
func TestLoadingContract(t *testing.T) {
	address := common.HexToAddress("0xb239A44548ec3813aCbBbe4017AcFfc541505b28")
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
executor address: 0x5b85f5666C9494e69A7ADB0CCe95ada892aB3607
contract address: 0xb239A44548ec3813aCbBbe4017AcFfc541505b28
to address: 0x4A2EBB506da083caC4d61f9305dF8967E595D16b
tx hash:
*/
func TestBuildMultisigTx(t *testing.T) {
	executor := kmssdk.KeyLabel{
		Key:     "secp256k1-hsm-1",
		Version: "cb848fb15e3a40b49bc41cbe957ea438",
		//Version:   "0179a6204ed7491ea5b27a87b541d5cb",
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

	contractAddress := common.HexToAddress("0xb239A44548ec3813aCbBbe4017AcFfc541505b28")

	to := kmssdk.KeyLabel{
		Key: "secp256k1-hsm-1",
		//Version: "cb848fb15e3a40b49bc41cbe957ea438",
		Version:   "0179a6204ed7491ea5b27a87b541d5cb",
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
		Amount:              big.NewInt(1e15),
		GasLimit:            1500000, // TODO: how to calc it?
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
		receipt, httpResponse, err := ubiAPIClient.TransactionsAPI.TxSend(ubiCtx, ubiPlatform, ubiNetwork).SignedTx(ubi.SignedTx{Tx: rawSignedTx}).Execute()
		if err != nil {
			t.Fatalf("failed to send tx %s \ntx hash %s\n %+v", rawSignedTx, txHash, err)
		}

		t.Logf("receipt: %+v \nhttp response: %+v", receipt, httpResponse)
	*/
}
