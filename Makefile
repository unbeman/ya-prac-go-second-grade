BINARY_NAME=gokeeper
CERT_PATH=cert

proto:
	protoc --go_out=api/v1 --go_opt=paths=source_relative --proto_path=api/v1 --go-grpc_out=api/v1 --go-grpc_opt=paths=source_relative api/v1/*.proto

mocks:
	mockgen -destination=internal/server/database/mock/database_mock.go "github.com/unbeman/ya-prac-go-second-grade/internal/server/database" Database
	mockgen -destination=internal/server/utils/mock/jwt_mock.go "github.com/unbeman/ya-prac-go-second-grade/internal/server/utils" IJWT

gen-certs:
	openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout ${CERT_PATH}/RootCA.key -out ${CERT_PATH}/RootCA.pem -subj "/C=RU/CN=Example-Root-CA"
	openssl x509 -outform pem -in ${CERT_PATH}/RootCA.pem -out ${CERT_PATH}/RootCA.crt
	openssl req -new -nodes -newkey rsa:2048 -keyout ${CERT_PATH}/server.key -out ${CERT_PATH}/server.csr -subj "/C=RU/ST=Russia/L=Moscow/O=Example-Certificates/CN=localhost"
	openssl x509 -req -sha256 -days 1024 -in server.csr -CA ${CERT_PATH}/RootCA.pem -CAkey ${CERT_PATH}/RootCA.key -CAcreateserial -extfile ${CERT_PATH}/domains.ext -out ${CERT_PATH}/server.crt

gen-jwt-key:
	openssl genrsa -out cert/jwt_key.pem 1024

local-server:
	go run cmd/server/main.go

local-client:
	go run cmd/client/main.go


build-client:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin cmd/client/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux cmd/client/main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows cmd/client/main.go

dc-server-up:
	docker-compose up -d