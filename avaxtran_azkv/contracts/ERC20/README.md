## ERC20 Contract

refs:  
contract from this: https://github.com/OpenZeppelin/openzeppelin-contracts/blob/release-v4.8/contracts/token/ERC20/ERC20.sol 
go: https://goethereumbook.org/smart-contract-deploy/

```script
solc --bin WFUSD.sol -o ./ --overwrite
solc --abi WFUSD.sol -o ./ --overwrite
abigen --abi=WFUSD.abi --bin=WFUSD.bin --pkg=erc20 --out wfusd.go
```