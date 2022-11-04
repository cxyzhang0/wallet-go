## Simple Multig Signature Contract

refs:  
contract from this with minor change: eth-multisig or github.com/lemon-sunxiansong/walletsdk
contract orig: https://github.com/christianlundkvist/simple-multisig  
go: https://goethereumbook.org/smart-contract-deploy/

```script
solc --bin SimpleMultiSig.sol -o ./ --overwrite
solc --abi SimpleMultiSig.sol -o ./ --overwrite
abigen --abi=SimpleMultiSig.abi --bin=SimpleMultiSig.bin --pkg=contract --out simplemultisig.go
```
after simplemultisig.go is generated, edit it by changing package name to contract1 from contract