## Example of BTC Transaction Build
It demonstrates the core utility functions in btcsuite for transaction building and signing.  
For simplicity, privKey, destination address, utxo, etc, are hard coded.  
It is a standalone app which does not have much relationship with the rest of the project.

### Run
```shell
go run main.go
```
It prints the signed transaction request, and its hash. The request can be broadcast via Blockcypher to the testnet.  
NOTE: It has already been submitted successfully and rebroadcast will result in double-spend error.