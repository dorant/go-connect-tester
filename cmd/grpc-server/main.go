//go:generate protoc -I ../../proto/hello --go_out=plugins=grpc:../../proto/hello ../../proto/hello/hello.proto

package main

import (
	"log"
	"net"

	pb "github.com/dorant/go-connect-tester/proto/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50000"
)

type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Call from client '%s' received\n", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	log.Printf("Start listen to: %s\n", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
