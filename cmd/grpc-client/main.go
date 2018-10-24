package main

import (
	"log"
	"os"
	"time"

	pb "github.com/dorant/go-connect-tester/proto/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "50000"
	}
	address := host + ":" + port

	log.Printf("Client connects to server at: %s\n", address)

	hostname, _ := os.Hostname()

	// Set up a connection to the server:
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	for {
		r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: hostname})
		if err != nil {
			log.Printf("Failed calling grpc::SayHello(): %v\n", err)
		} else {
			log.Printf("Call to server returned: %s\n", r)
		}

		time.Sleep(10 * time.Second)
	}
}
