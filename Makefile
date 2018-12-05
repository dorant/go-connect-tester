
HTTP_CLIENT_TAG = 'bjornsv/http-client:1.1'
HTTP_SERVER_TAG = 'bjornsv/http-server:1.1'
GRPC_CLIENT_TAG = 'bjornsv/grpc-client:1.2'
GRPC_SERVER_TAG = 'bjornsv/grpc-server:1.2'

.PHONY: http-client http-server grpc-client grpc-server gen deps

all: http-client http-server grpc-client grpc-server

http-client:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o http-client cmd/http-client/main.go
	docker build . -t $(HTTP_CLIENT_TAG) -f Dockerfile.http-client

http-server:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o http-server cmd/http-server/main.go
	docker build . -t $(HTTP_SERVER_TAG) -f Dockerfile.http-server

grpc-client:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o grpc-client cmd/grpc-client/main.go
	docker build . -t $(GRPC_CLIENT_TAG) -f Dockerfile.grpc-client

grpc-server:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o grpc-server cmd/grpc-server/main.go
	docker build . -t $(GRPC_SERVER_TAG) -f Dockerfile.grpc-server

push:
	@@echo
	@@echo "Dont forget: docker login -u <dockerhub-username>"
	@@echo
	docker push $(HTTP_CLIENT_TAG)
	docker push $(HTTP_SERVER_TAG)
	docker push $(GRPC_CLIENT_TAG)
	docker push $(GRPC_SERVER_TAG)


clean:
	rm -rf http-server http-client grpc-server grpc-client
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

