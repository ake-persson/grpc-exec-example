all: proto

proto:
	protoc -I info/ info/info.proto --go_out=plugins=grpc:info/
	protoc -I exec/ exec/exec.proto --go_out=plugins=grpc:exec/
	protoc -I auth/ auth/auth.proto --go_out=plugins=grpc:auth/

proto-python:
	python -m grpc_tools.protoc -I info --python_out=info --grpc_python_out=info info/info.proto
	python -m grpc_tools.protoc -I exec --python_out=exec --grpc_python_out=exec exec/exec.proto
	python -m grpc_tools.protoc -I auth --python_out=auth --grpc_python_out=auth auth/auth.proto

build:
	cd auth-server && go build
	cd exec-server && go build
	cd info-server && go build
	cd catalog-server && go build
	cd client && go build

darwin:
	cd auth-server && GOOS=darwin GOARCH=amd64 go build
	cd exec-server && GOOS=darwin GOARCH=amd64 go build
	cd info-server && GOOS=darwin GOARCH=amd64 go build
	cd catalog-server && GOOS=darwin GOARCH=amd64 go build
	cd client && GOOS=darwin GOARCH=amd64 go build

linux:
	cd auth-server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	cd exec-server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	cd info-server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	cd catalog-server && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
	cd client && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

.PHONY: proto build darwin linux
