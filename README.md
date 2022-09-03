## BITCOIN HD WALLET WITH GOLANG AND GRPC  
Sample codes for grpc with protobuf to create client/server to
create mnemonic-based BTC wallet, retrieve wallet from mnemonic and get balance for a given address via Blockcypher API.

### Prerequisites
brew install protobuf  
protoc --version  

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest  
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  

### References:
https://hackernoon.com/how-to-create-a-bitcoin-hd-wallet-with-golang-and-grpc-part-l-u51d3wwm <br />
https://levelup.gitconnected.com/bitcoin-hd-wallet-with-golang-and-grpc-part-l-56d8df05c602

### build protobuf
```shell
make build_protoc
```
### server
#### build
```shell
go build main.go
````
#### run
```shell
./main  
or  
go run main.go
```
### client
```shell
cd client
```
#### build
```shell
go build client.go
```
#### run
```shell
// create a wallet
./client -m=create-wallet  
or
go run client.go -m=create-wallet  

// get wallet from mnemonic
./client -m=get-wallet -mne="crucial clinic obscure good creek brand sunset grit coral mention off hint" 

// get balance for adddress
./client -m=get-balance -addr=1Go23sv8vR81YuV1hHGsUrdyjLcGVUpCDy
```
