build_protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/btchdwallet/wallet.proto
old:
	#protoc -I=. --go_out=proto/btchdwallet proto/btchdwallet/wallet.proto
	#protoc -I=. --go-grpc_out=proto/btchdwallet proto/btchdwallet/wallet.proto
	#protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative proto/btchdwallet/wallet.proto

