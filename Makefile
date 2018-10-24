all: client server

CLIENT_TAG = 'bjornsv/grpc-client:1.2'
SERVER_TAG = 'bjornsv/grpc-server:1.2'

.PHONY: client server gen deps
client:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o client cmd/client/main.go
	docker build . -t $(CLIENT_TAG) -f Dockerfile.client

server:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o server cmd/server/main.go
	docker build . -t $(SERVER_TAG) -f Dockerfile.server

push:
	@@echo
	@@echo "Dont forget: docker login -u <dockerhub-username>"
	@@echo
	docker push $(CLIENT_TAG)
	docker push $(SERVER_TAG)


clean:
	rm -rf server client
gen:
	go generate ./...

deps:
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u google.golang.org/grpc
	@@echo
	@@echo "* Manually install protoc from https://github.com/protocolbuffers/protobuf/releases"
	@@echo "> curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip"
	@@echo "> unzip -o protoc-3.6.1-linux-x86_64.zip -d ~/ bin/protoc"

