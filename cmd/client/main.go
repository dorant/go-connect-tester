package main

import (
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/dorant/go-grpc-tester/proto/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	address     = "localhost:50000"
	defaultName = "world"
)

func waitForGrpcReady(cc grpc.ClientConn, timeout int) error {
	fmt.Printf("Wait %d sec for ready state...: (%v)", timeout, cc.GetState())
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	for {
		s := cc.GetState()
		if s == connectivity.Ready {
			break
		}
		if !cc.WaitForStateChange(ctx, s) {
			// ctx got timeout or canceled.
			return ctx.Err()
		}
	}
	return nil
}

func main() {
	// Set up a connection to the server.:
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	err = waitForGrpcReady(*conn, 10)
	if err != nil {
		log.Printf("did not connect: %v\n", err)
	}
	fmt.Printf("Done (%v)\n", conn.GetState())

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	for i := 0; i < 40; i++ {
		log.Printf("State: %v\n", conn.GetState())
		fmt.Printf("Wait no. %d...: ", i)
		time.Sleep(2 * time.Second)

		r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			log.Printf("could not greet: %v\n", err)
		} else {
			fmt.Printf("Done: %s\n", r)
		}
	}
}
