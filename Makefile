gen-calculator-protos:
	protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import calculator/protos/calculator/calculator