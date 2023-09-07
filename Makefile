BINARY_NAME=gokeeper
BINARY_PATH=bin

proto:
    protoc --go_out=api/v1 --go_opt=paths=source_relative --proto_path=api/v1 --go-grpc_out=api/v1 --go-grpc_opt=paths=source_relative api/v1/*.proto


server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080


#build-client:
#    GOARCH=amd64 GOOS=darwin go build -o ${BINARY_PATH}/${BINARY_NAME}-darwin main.go
#    GOARCH=amd64 GOOS=linux go build -o ${BINARY_PATH}/${BINARY_NAME}-linux main.go
#    GOARCH=amd64 GOOS=windows go build -o ${BINARY_PATH}/${BINARY_NAME}-windows main.go