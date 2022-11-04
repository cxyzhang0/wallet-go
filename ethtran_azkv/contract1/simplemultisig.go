// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract1

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold_\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"owners_\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"sperator\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"txInputHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"totalHash\",\"type\":\"bytes32\"}],\"name\":\"ExecuteLog\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"ExecuteStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ExecuteVerifySender\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"RecoverStart\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"RecoverVerify\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"RecoverdAddr\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint8[]\",\"name\":\"sigV\",\"type\":\"uint8[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"sigR\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"sigS\",\"type\":\"bytes32[]\"},{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"destination\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"getHashes\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwersLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ownersArr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001b2838038062001b288339818101604052810190620000379190620005a8565b600a8251111580156200004b575081518311155b8015620000585750600083115b6200009a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000919062000684565b60405180910390fd5b6000805b8351811015620001e1578173ffffffffffffffffffffffffffffffffffffffff16848281518110620000d557620000d4620006a6565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff161162000136576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016200012d9062000725565b60405180910390fd5b600160026000868481518110620001525762000151620006a6565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550838181518110620001c157620001c0620006a6565b5b602002602001015191508080620001d89062000776565b9150506200009e565b508260039080519060200190620001fa929190620002d1565b50836001819055507fd87cd6ef79d4e2b95e15ce8abf732db51ec771f1ca2edccf22a46c729ac5647260001b7fb7a0bfa1b79f2443f4d73ebb9259cddbcd510b18be6fc4da7d1aa7b1786e73e660001b7fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc660001b84307f251543af6a222378665a76fe38dbceae4871a070b7fdaf5c6c30cf758dc33cc060001b604051602001620002ab969594939291906200085a565b6040516020818303038152906040528051906020012060048190555050505050620008c7565b8280548282559060005260206000209081019282156200034d579160200282015b828111156200034c5782518260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555091602001919060010190620002f2565b5b5090506200035c919062000360565b5090565b5b808211156200037b57600081600090555060010162000361565b5090565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b620003a88162000393565b8114620003b457600080fd5b50565b600081519050620003c8816200039d565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200041e82620003d3565b810181811067ffffffffffffffff8211171562000440576200043f620003e4565b5b80604052505050565b6000620004556200037f565b905062000463828262000413565b919050565b600067ffffffffffffffff821115620004865762000485620003e4565b5b602082029050602081019050919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620004c9826200049c565b9050919050565b620004db81620004bc565b8114620004e757600080fd5b50565b600081519050620004fb81620004d0565b92915050565b600062000518620005128462000468565b62000449565b905080838252602082019050602084028301858111156200053e576200053d62000497565b5b835b818110156200056b5780620005568882620004ea565b84526020840193505060208101905062000540565b5050509392505050565b600082601f8301126200058d576200058c620003ce565b5b81516200059f84826020860162000501565b91505092915050565b600080600060608486031215620005c457620005c362000389565b5b6000620005d486828701620003b7565b935050602084015167ffffffffffffffff811115620005f857620005f76200038e565b5b620006068682870162000575565b92505060406200061986828701620003b7565b9150509250925092565b600082825260208201905092915050565b7f303c7468726573686f6c643c6f776e6572732e6c656e67746800000000000000600082015250565b60006200066c60198362000623565b9150620006798262000634565b602082019050919050565b600060208201905081810360008301526200069f816200065d565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f7265706561746564206f776e6572206f72206e6f7420736f7274656400000000600082015250565b60006200070d601c8362000623565b91506200071a82620006d5565b602082019050919050565b600060208201905081810360008301526200074081620006fe565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000620007838262000393565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203620007b857620007b762000747565b5b600182019050919050565b6000819050919050565b620007d881620007c3565b82525050565b620007e98162000393565b82525050565b6000819050919050565b60006200081a620008146200080e846200049c565b620007ef565b6200049c565b9050919050565b60006200082e82620007f9565b9050919050565b6000620008428262000821565b9050919050565b620008548162000835565b82525050565b600060c082019050620008716000830189620007cd565b620008806020830188620007cd565b6200088f6040830187620007cd565b6200089e6060830186620007de565b620008ad608083018562000849565b620008bc60a0830184620007cd565b979650505050505050565b61125180620008d76000396000f3fe6080604052600436106100745760003560e01c8063aa5df9e21161004e578063aa5df9e2146100ff578063affed0e01461013c578063ca7541ee14610167578063f87c78c7146101925761007b565b80630d8e6e2c1461008057806342cde4e8146100ab578063a0ab9653146100d65761007b565b3661007b57005b600080fd5b34801561008c57600080fd5b506100956101d1565b6040516100a29190610779565b60405180910390f35b3480156100b757600080fd5b506100c061020e565b6040516100cd91906107b4565b60405180910390f35b3480156100e257600080fd5b506100fd60048036038101906100f89190610b9c565b610214565b005b34801561010b57600080fd5b5061012660048036038101906101219190610cc2565b6105e9565b6040516101339190610cfe565b60405180910390f35b34801561014857600080fd5b50610151610628565b60405161015e91906107b4565b60405180910390f35b34801561017357600080fd5b5061017c61062e565b60405161018991906107b4565b60405180910390f35b34801561019e57600080fd5b506101b960048036038101906101b49190610d19565b61063b565b6040516101c893929190610dbf565b60405180910390f35b60606040518060400160405280600481526020017f322e333300000000000000000000000000000000000000000000000000000000815250905090565b60015481565b7f1f3748a2491ab38f2844a84540e798f06e240ec9764f4034a4b85d9de95a309760405160405180910390a1600154875114610285576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161027c90610e42565b60405180910390fd5b85518751148015610297575087518751145b6102d6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102cd90610eae565b60405180910390fd5b7f7b3f83512e4134c9157a582e9b708d7b8535a483ffdac94c37ecebb8a3b63c04336040516103059190610cfe565b60405180910390a13373ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614806103735750600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b6103b2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103a990610f1a565b60405180910390fd5b60008060006103c4888888888861063b565b9250925092507f672ea124619314d71be6f340ecb56da6038c8d6b8ffca9bb1be62086d0a73d908383836040516103fd93929190610dbf565b60405180910390a16000805b600154811015610574577fbb8c691c28385da4e4d29a158e660fad1c741f388f2170d3c9b67b6d71ab6d128160405161044291906107b4565b60405180910390a160006001848f848151811061046257610461610f3a565b5b60200260200101518f858151811061047d5761047c610f3a565b5b60200260200101518f868151811061049857610497610f3a565b5b6020026020010151604051600081526020016040526040516104bd9493929190610f78565b6020604051602081039080840390855afa1580156104df573d6000803e3d6000fd5b5050506020604051035190507f464d905a75ac90e0d07b8c2a0cb67371b7f6abd04160c6e323686e3d9beb72b9818360405161051c929190610fbd565b60405180910390a17f4e57899e25b61543bc91679ea2a1d9edf6409fad79e539a73eab5f03c06d77cd8183604051610555929190610fbd565b60405180910390a180925050808061056c90611015565b915050610409565b506001600054610584919061105d565b60008190555060008080895160208b018c8e8bf19050806105da576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105d1906110dd565b60405180910390fd5b50505050505050505050505050565b600381815481106105f957600080fd5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60005481565b6000600380549050905090565b6000806000807f3ee892349ae4bbe61dce18f95115b5dc02daf49204cc602458cd4c1f540d56d760001b898989805190602001206000548a8a60405160200161068a97969594939291906110fd565b6040516020818303038152906040528051906020012090506000600454826040516020016106b99291906111e4565b60405160208183030381529060405280519060200120905060045482829450945094505050955095509592505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610723578082015181840152602081019050610708565b60008484015250505050565b6000601f19601f8301169050919050565b600061074b826106e9565b61075581856106f4565b9350610765818560208601610705565b61076e8161072f565b840191505092915050565b600060208201905081810360008301526107938184610740565b905092915050565b6000819050919050565b6107ae8161079b565b82525050565b60006020820190506107c960008301846107a5565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6108208261072f565b810181811067ffffffffffffffff8211171561083f5761083e6107e8565b5b80604052505050565b60006108526107cf565b905061085e8282610817565b919050565b600067ffffffffffffffff82111561087e5761087d6107e8565b5b602082029050602081019050919050565b600080fd5b600060ff82169050919050565b6108aa81610894565b81146108b557600080fd5b50565b6000813590506108c7816108a1565b92915050565b60006108e06108db84610863565b610848565b905080838252602082019050602084028301858111156109035761090261088f565b5b835b8181101561092c578061091888826108b8565b845260208401935050602081019050610905565b5050509392505050565b600082601f83011261094b5761094a6107e3565b5b813561095b8482602086016108cd565b91505092915050565b600067ffffffffffffffff82111561097f5761097e6107e8565b5b602082029050602081019050919050565b6000819050919050565b6109a381610990565b81146109ae57600080fd5b50565b6000813590506109c08161099a565b92915050565b60006109d96109d484610964565b610848565b905080838252602082019050602084028301858111156109fc576109fb61088f565b5b835b81811015610a255780610a1188826109b1565b8452602084019350506020810190506109fe565b5050509392505050565b600082601f830112610a4457610a436107e3565b5b8135610a548482602086016109c6565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610a8882610a5d565b9050919050565b610a9881610a7d565b8114610aa357600080fd5b50565b600081359050610ab581610a8f565b92915050565b610ac48161079b565b8114610acf57600080fd5b50565b600081359050610ae181610abb565b92915050565b600080fd5b600067ffffffffffffffff821115610b0757610b066107e8565b5b610b108261072f565b9050602081019050919050565b82818337600083830152505050565b6000610b3f610b3a84610aec565b610848565b905082815260208101848484011115610b5b57610b5a610ae7565b5b610b66848285610b1d565b509392505050565b600082601f830112610b8357610b826107e3565b5b8135610b93848260208601610b2c565b91505092915050565b600080600080600080600080610100898b031215610bbd57610bbc6107d9565b5b600089013567ffffffffffffffff811115610bdb57610bda6107de565b5b610be78b828c01610936565b985050602089013567ffffffffffffffff811115610c0857610c076107de565b5b610c148b828c01610a2f565b975050604089013567ffffffffffffffff811115610c3557610c346107de565b5b610c418b828c01610a2f565b9650506060610c528b828c01610aa6565b9550506080610c638b828c01610ad2565b94505060a089013567ffffffffffffffff811115610c8457610c836107de565b5b610c908b828c01610b6e565b93505060c0610ca18b828c01610aa6565b92505060e0610cb28b828c01610ad2565b9150509295985092959890939650565b600060208284031215610cd857610cd76107d9565b5b6000610ce684828501610ad2565b91505092915050565b610cf881610a7d565b82525050565b6000602082019050610d136000830184610cef565b92915050565b600080600080600060a08688031215610d3557610d346107d9565b5b6000610d4388828901610aa6565b9550506020610d5488828901610ad2565b945050604086013567ffffffffffffffff811115610d7557610d746107de565b5b610d8188828901610b6e565b9350506060610d9288828901610aa6565b9250506080610da388828901610ad2565b9150509295509295909350565b610db981610990565b82525050565b6000606082019050610dd46000830186610db0565b610de16020830185610db0565b610dee6040830184610db0565b949350505050565b7f6e6f7420657175616c20746f207468726573686f6c6400000000000000000000600082015250565b6000610e2c6016836106f4565b9150610e3782610df6565b602082019050919050565b60006020820190508181036000830152610e5b81610e1f565b9050919050565b7f6c656e677468206e6f74206d6174636800000000000000000000000000000000600082015250565b6000610e986010836106f4565b9150610ea382610e62565b602082019050919050565b60006020820190508181036000830152610ec781610e8b565b9050919050565b7f77726f6e67206578656375746f72000000000000000000000000000000000000600082015250565b6000610f04600e836106f4565b9150610f0f82610ece565b602082019050919050565b60006020820190508181036000830152610f3381610ef7565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b610f7281610894565b82525050565b6000608082019050610f8d6000830187610db0565b610f9a6020830186610f69565b610fa76040830185610db0565b610fb46060830184610db0565b95945050505050565b6000604082019050610fd26000830185610cef565b610fdf60208301846107a5565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006110208261079b565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361105257611051610fe6565b5b600182019050919050565b60006110688261079b565b91506110738361079b565b925082820190508082111561108b5761108a610fe6565b5b92915050565b7f6e6f745f73756363657373000000000000000000000000000000000000000000600082015250565b60006110c7600b836106f4565b91506110d282611091565b602082019050919050565b600060208201905081810360008301526110f6816110ba565b9050919050565b600060e082019050611112600083018a610db0565b61111f6020830189610cef565b61112c60408301886107a5565b6111396060830187610db0565b61114660808301866107a5565b61115360a0830185610cef565b61116060c08301846107a5565b98975050505050505050565b600081905092915050565b7f1901000000000000000000000000000000000000000000000000000000000000600082015250565b60006111ad60028361116c565b91506111b882611177565b600282019050919050565b6000819050919050565b6111de6111d982610990565b6111c3565b82525050565b60006111ef826111a0565b91506111fb82856111cd565b60208201915061120b82846111cd565b602082019150819050939250505056fea26469706673582212204887eb57ea7f040878de6d126780e3bbc08471aded194d4ef5414fb1d72847de64736f6c63430008110033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend, threshold_ *big.Int, owners_ []common.Address, chainId *big.Int) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend, threshold_, owners_, chainId)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// GetHashes is a free data retrieval call binding the contract method 0xf87c78c7.
//
// Solidity: function getHashes(address destination, uint256 value, bytes data, address executor, uint256 gasLimit) view returns(bytes32, bytes32, bytes32)
func (_Contract *ContractCaller) GetHashes(opts *bind.CallOpts, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) ([32]byte, [32]byte, [32]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getHashes", destination, value, data, executor, gasLimit)

	if err != nil {
		return *new([32]byte), *new([32]byte), *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	out2 := *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return out0, out1, out2, err

}

// GetHashes is a free data retrieval call binding the contract method 0xf87c78c7.
//
// Solidity: function getHashes(address destination, uint256 value, bytes data, address executor, uint256 gasLimit) view returns(bytes32, bytes32, bytes32)
func (_Contract *ContractSession) GetHashes(destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) ([32]byte, [32]byte, [32]byte, error) {
	return _Contract.Contract.GetHashes(&_Contract.CallOpts, destination, value, data, executor, gasLimit)
}

// GetHashes is a free data retrieval call binding the contract method 0xf87c78c7.
//
// Solidity: function getHashes(address destination, uint256 value, bytes data, address executor, uint256 gasLimit) view returns(bytes32, bytes32, bytes32)
func (_Contract *ContractCallerSession) GetHashes(destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) ([32]byte, [32]byte, [32]byte, error) {
	return _Contract.Contract.GetHashes(&_Contract.CallOpts, destination, value, data, executor, gasLimit)
}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
//
// Solidity: function getOwersLength() view returns(uint256)
func (_Contract *ContractCaller) GetOwersLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getOwersLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
//
// Solidity: function getOwersLength() view returns(uint256)
func (_Contract *ContractSession) GetOwersLength() (*big.Int, error) {
	return _Contract.Contract.GetOwersLength(&_Contract.CallOpts)
}

// GetOwersLength is a free data retrieval call binding the contract method 0xca7541ee.
//
// Solidity: function getOwersLength() view returns(uint256)
func (_Contract *ContractCallerSession) GetOwersLength() (*big.Int, error) {
	return _Contract.Contract.GetOwersLength(&_Contract.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() pure returns(string)
func (_Contract *ContractCaller) GetVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() pure returns(string)
func (_Contract *ContractSession) GetVersion() (string, error) {
	return _Contract.Contract.GetVersion(&_Contract.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() pure returns(string)
func (_Contract *ContractCallerSession) GetVersion() (string, error) {
	return _Contract.Contract.GetVersion(&_Contract.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_Contract *ContractCaller) Nonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "nonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_Contract *ContractSession) Nonce() (*big.Int, error) {
	return _Contract.Contract.Nonce(&_Contract.CallOpts)
}

// Nonce is a free data retrieval call binding the contract method 0xaffed0e0.
//
// Solidity: function nonce() view returns(uint256)
func (_Contract *ContractCallerSession) Nonce() (*big.Int, error) {
	return _Contract.Contract.Nonce(&_Contract.CallOpts)
}

// OwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
//
// Solidity: function ownersArr(uint256 ) view returns(address)
func (_Contract *ContractCaller) OwnersArr(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "ownersArr", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
//
// Solidity: function ownersArr(uint256 ) view returns(address)
func (_Contract *ContractSession) OwnersArr(arg0 *big.Int) (common.Address, error) {
	return _Contract.Contract.OwnersArr(&_Contract.CallOpts, arg0)
}

// OwnersArr is a free data retrieval call binding the contract method 0xaa5df9e2.
//
// Solidity: function ownersArr(uint256 ) view returns(address)
func (_Contract *ContractCallerSession) OwnersArr(arg0 *big.Int) (common.Address, error) {
	return _Contract.Contract.OwnersArr(&_Contract.CallOpts, arg0)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint256)
func (_Contract *ContractCaller) Threshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "threshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint256)
func (_Contract *ContractSession) Threshold() (*big.Int, error) {
	return _Contract.Contract.Threshold(&_Contract.CallOpts)
}

// Threshold is a free data retrieval call binding the contract method 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint256)
func (_Contract *ContractCallerSession) Threshold() (*big.Int, error) {
	return _Contract.Contract.Threshold(&_Contract.CallOpts)
}

// Execute is a paid mutator transaction binding the contract method 0xa0ab9653.
//
// Solidity: function execute(uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_Contract *ContractTransactor) Execute(opts *bind.TransactOpts, sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "execute", sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}

// Execute is a paid mutator transaction binding the contract method 0xa0ab9653.
//
// Solidity: function execute(uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_Contract *ContractSession) Execute(sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Execute(&_Contract.TransactOpts, sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}

// Execute is a paid mutator transaction binding the contract method 0xa0ab9653.
//
// Solidity: function execute(uint8[] sigV, bytes32[] sigR, bytes32[] sigS, address destination, uint256 value, bytes data, address executor, uint256 gasLimit) returns()
func (_Contract *ContractTransactorSession) Execute(sigV []uint8, sigR [][32]byte, sigS [][32]byte, destination common.Address, value *big.Int, data []byte, executor common.Address, gasLimit *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Execute(&_Contract.TransactOpts, sigV, sigR, sigS, destination, value, data, executor, gasLimit)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Contract *ContractTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Contract *ContractSession) Receive() (*types.Transaction, error) {
	return _Contract.Contract.Receive(&_Contract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Contract *ContractTransactorSession) Receive() (*types.Transaction, error) {
	return _Contract.Contract.Receive(&_Contract.TransactOpts)
}

// ContractExecuteLogIterator is returned from FilterExecuteLog and is used to iterate over the raw logs and unpacked data for ExecuteLog events raised by the Contract contract.
type ContractExecuteLogIterator struct {
	Event *ContractExecuteLog // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractExecuteLogIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractExecuteLog)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractExecuteLog)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractExecuteLogIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractExecuteLogIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractExecuteLog represents a ExecuteLog event raised by the Contract contract.
type ContractExecuteLog struct {
	Sperator    [32]byte
	TxInputHash [32]byte
	TotalHash   [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterExecuteLog is a free log retrieval operation binding the contract event 0x672ea124619314d71be6f340ecb56da6038c8d6b8ffca9bb1be62086d0a73d90.
//
// Solidity: event ExecuteLog(bytes32 sperator, bytes32 txInputHash, bytes32 totalHash)
func (_Contract *ContractFilterer) FilterExecuteLog(opts *bind.FilterOpts) (*ContractExecuteLogIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ExecuteLog")
	if err != nil {
		return nil, err
	}
	return &ContractExecuteLogIterator{contract: _Contract.contract, event: "ExecuteLog", logs: logs, sub: sub}, nil
}

// WatchExecuteLog is a free log subscription operation binding the contract event 0x672ea124619314d71be6f340ecb56da6038c8d6b8ffca9bb1be62086d0a73d90.
//
// Solidity: event ExecuteLog(bytes32 sperator, bytes32 txInputHash, bytes32 totalHash)
func (_Contract *ContractFilterer) WatchExecuteLog(opts *bind.WatchOpts, sink chan<- *ContractExecuteLog) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ExecuteLog")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractExecuteLog)
				if err := _Contract.contract.UnpackLog(event, "ExecuteLog", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseExecuteLog is a log parse operation binding the contract event 0x672ea124619314d71be6f340ecb56da6038c8d6b8ffca9bb1be62086d0a73d90.
//
// Solidity: event ExecuteLog(bytes32 sperator, bytes32 txInputHash, bytes32 totalHash)
func (_Contract *ContractFilterer) ParseExecuteLog(log types.Log) (*ContractExecuteLog, error) {
	event := new(ContractExecuteLog)
	if err := _Contract.contract.UnpackLog(event, "ExecuteLog", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractExecuteStartIterator is returned from FilterExecuteStart and is used to iterate over the raw logs and unpacked data for ExecuteStart events raised by the Contract contract.
type ContractExecuteStartIterator struct {
	Event *ContractExecuteStart // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractExecuteStartIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractExecuteStart)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractExecuteStart)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractExecuteStartIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractExecuteStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractExecuteStart represents a ExecuteStart event raised by the Contract contract.
type ContractExecuteStart struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterExecuteStart is a free log retrieval operation binding the contract event 0x1f3748a2491ab38f2844a84540e798f06e240ec9764f4034a4b85d9de95a3097.
//
// Solidity: event ExecuteStart()
func (_Contract *ContractFilterer) FilterExecuteStart(opts *bind.FilterOpts) (*ContractExecuteStartIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ExecuteStart")
	if err != nil {
		return nil, err
	}
	return &ContractExecuteStartIterator{contract: _Contract.contract, event: "ExecuteStart", logs: logs, sub: sub}, nil
}

// WatchExecuteStart is a free log subscription operation binding the contract event 0x1f3748a2491ab38f2844a84540e798f06e240ec9764f4034a4b85d9de95a3097.
//
// Solidity: event ExecuteStart()
func (_Contract *ContractFilterer) WatchExecuteStart(opts *bind.WatchOpts, sink chan<- *ContractExecuteStart) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ExecuteStart")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractExecuteStart)
				if err := _Contract.contract.UnpackLog(event, "ExecuteStart", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseExecuteStart is a log parse operation binding the contract event 0x1f3748a2491ab38f2844a84540e798f06e240ec9764f4034a4b85d9de95a3097.
//
// Solidity: event ExecuteStart()
func (_Contract *ContractFilterer) ParseExecuteStart(log types.Log) (*ContractExecuteStart, error) {
	event := new(ContractExecuteStart)
	if err := _Contract.contract.UnpackLog(event, "ExecuteStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractExecuteVerifySenderIterator is returned from FilterExecuteVerifySender and is used to iterate over the raw logs and unpacked data for ExecuteVerifySender events raised by the Contract contract.
type ContractExecuteVerifySenderIterator struct {
	Event *ContractExecuteVerifySender // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractExecuteVerifySenderIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractExecuteVerifySender)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractExecuteVerifySender)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractExecuteVerifySenderIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractExecuteVerifySenderIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractExecuteVerifySender represents a ExecuteVerifySender event raised by the Contract contract.
type ContractExecuteVerifySender struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExecuteVerifySender is a free log retrieval operation binding the contract event 0x7b3f83512e4134c9157a582e9b708d7b8535a483ffdac94c37ecebb8a3b63c04.
//
// Solidity: event ExecuteVerifySender(address sender)
func (_Contract *ContractFilterer) FilterExecuteVerifySender(opts *bind.FilterOpts) (*ContractExecuteVerifySenderIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ExecuteVerifySender")
	if err != nil {
		return nil, err
	}
	return &ContractExecuteVerifySenderIterator{contract: _Contract.contract, event: "ExecuteVerifySender", logs: logs, sub: sub}, nil
}

// WatchExecuteVerifySender is a free log subscription operation binding the contract event 0x7b3f83512e4134c9157a582e9b708d7b8535a483ffdac94c37ecebb8a3b63c04.
//
// Solidity: event ExecuteVerifySender(address sender)
func (_Contract *ContractFilterer) WatchExecuteVerifySender(opts *bind.WatchOpts, sink chan<- *ContractExecuteVerifySender) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ExecuteVerifySender")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractExecuteVerifySender)
				if err := _Contract.contract.UnpackLog(event, "ExecuteVerifySender", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseExecuteVerifySender is a log parse operation binding the contract event 0x7b3f83512e4134c9157a582e9b708d7b8535a483ffdac94c37ecebb8a3b63c04.
//
// Solidity: event ExecuteVerifySender(address sender)
func (_Contract *ContractFilterer) ParseExecuteVerifySender(log types.Log) (*ContractExecuteVerifySender, error) {
	event := new(ContractExecuteVerifySender)
	if err := _Contract.contract.UnpackLog(event, "ExecuteVerifySender", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRecoverStartIterator is returned from FilterRecoverStart and is used to iterate over the raw logs and unpacked data for RecoverStart events raised by the Contract contract.
type ContractRecoverStartIterator struct {
	Event *ContractRecoverStart // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractRecoverStartIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRecoverStart)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractRecoverStart)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractRecoverStartIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRecoverStartIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRecoverStart represents a RecoverStart event raised by the Contract contract.
type ContractRecoverStart struct {
	I   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRecoverStart is a free log retrieval operation binding the contract event 0xbb8c691c28385da4e4d29a158e660fad1c741f388f2170d3c9b67b6d71ab6d12.
//
// Solidity: event RecoverStart(uint256 i)
func (_Contract *ContractFilterer) FilterRecoverStart(opts *bind.FilterOpts) (*ContractRecoverStartIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "RecoverStart")
	if err != nil {
		return nil, err
	}
	return &ContractRecoverStartIterator{contract: _Contract.contract, event: "RecoverStart", logs: logs, sub: sub}, nil
}

// WatchRecoverStart is a free log subscription operation binding the contract event 0xbb8c691c28385da4e4d29a158e660fad1c741f388f2170d3c9b67b6d71ab6d12.
//
// Solidity: event RecoverStart(uint256 i)
func (_Contract *ContractFilterer) WatchRecoverStart(opts *bind.WatchOpts, sink chan<- *ContractRecoverStart) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "RecoverStart")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRecoverStart)
				if err := _Contract.contract.UnpackLog(event, "RecoverStart", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRecoverStart is a log parse operation binding the contract event 0xbb8c691c28385da4e4d29a158e660fad1c741f388f2170d3c9b67b6d71ab6d12.
//
// Solidity: event RecoverStart(uint256 i)
func (_Contract *ContractFilterer) ParseRecoverStart(log types.Log) (*ContractRecoverStart, error) {
	event := new(ContractRecoverStart)
	if err := _Contract.contract.UnpackLog(event, "RecoverStart", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRecoverVerifyIterator is returned from FilterRecoverVerify and is used to iterate over the raw logs and unpacked data for RecoverVerify events raised by the Contract contract.
type ContractRecoverVerifyIterator struct {
	Event *ContractRecoverVerify // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractRecoverVerifyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRecoverVerify)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractRecoverVerify)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractRecoverVerifyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRecoverVerifyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRecoverVerify represents a RecoverVerify event raised by the Contract contract.
type ContractRecoverVerify struct {
	Addr common.Address
	I    *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRecoverVerify is a free log retrieval operation binding the contract event 0x464d905a75ac90e0d07b8c2a0cb67371b7f6abd04160c6e323686e3d9beb72b9.
//
// Solidity: event RecoverVerify(address addr, uint256 i)
func (_Contract *ContractFilterer) FilterRecoverVerify(opts *bind.FilterOpts) (*ContractRecoverVerifyIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "RecoverVerify")
	if err != nil {
		return nil, err
	}
	return &ContractRecoverVerifyIterator{contract: _Contract.contract, event: "RecoverVerify", logs: logs, sub: sub}, nil
}

// WatchRecoverVerify is a free log subscription operation binding the contract event 0x464d905a75ac90e0d07b8c2a0cb67371b7f6abd04160c6e323686e3d9beb72b9.
//
// Solidity: event RecoverVerify(address addr, uint256 i)
func (_Contract *ContractFilterer) WatchRecoverVerify(opts *bind.WatchOpts, sink chan<- *ContractRecoverVerify) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "RecoverVerify")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRecoverVerify)
				if err := _Contract.contract.UnpackLog(event, "RecoverVerify", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRecoverVerify is a log parse operation binding the contract event 0x464d905a75ac90e0d07b8c2a0cb67371b7f6abd04160c6e323686e3d9beb72b9.
//
// Solidity: event RecoverVerify(address addr, uint256 i)
func (_Contract *ContractFilterer) ParseRecoverVerify(log types.Log) (*ContractRecoverVerify, error) {
	event := new(ContractRecoverVerify)
	if err := _Contract.contract.UnpackLog(event, "RecoverVerify", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractRecoverdAddrIterator is returned from FilterRecoverdAddr and is used to iterate over the raw logs and unpacked data for RecoverdAddr events raised by the Contract contract.
type ContractRecoverdAddrIterator struct {
	Event *ContractRecoverdAddr // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractRecoverdAddrIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRecoverdAddr)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractRecoverdAddr)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractRecoverdAddrIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRecoverdAddrIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRecoverdAddr represents a RecoverdAddr event raised by the Contract contract.
type ContractRecoverdAddr struct {
	Addr common.Address
	I    *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRecoverdAddr is a free log retrieval operation binding the contract event 0x4e57899e25b61543bc91679ea2a1d9edf6409fad79e539a73eab5f03c06d77cd.
//
// Solidity: event RecoverdAddr(address addr, uint256 i)
func (_Contract *ContractFilterer) FilterRecoverdAddr(opts *bind.FilterOpts) (*ContractRecoverdAddrIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "RecoverdAddr")
	if err != nil {
		return nil, err
	}
	return &ContractRecoverdAddrIterator{contract: _Contract.contract, event: "RecoverdAddr", logs: logs, sub: sub}, nil
}

// WatchRecoverdAddr is a free log subscription operation binding the contract event 0x4e57899e25b61543bc91679ea2a1d9edf6409fad79e539a73eab5f03c06d77cd.
//
// Solidity: event RecoverdAddr(address addr, uint256 i)
func (_Contract *ContractFilterer) WatchRecoverdAddr(opts *bind.WatchOpts, sink chan<- *ContractRecoverdAddr) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "RecoverdAddr")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRecoverdAddr)
				if err := _Contract.contract.UnpackLog(event, "RecoverdAddr", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRecoverdAddr is a log parse operation binding the contract event 0x4e57899e25b61543bc91679ea2a1d9edf6409fad79e539a73eab5f03c06d77cd.
//
// Solidity: event RecoverdAddr(address addr, uint256 i)
func (_Contract *ContractFilterer) ParseRecoverdAddr(log types.Log) (*ContractRecoverdAddr, error) {
	event := new(ContractRecoverdAddr)
	if err := _Contract.contract.UnpackLog(event, "RecoverdAddr", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
