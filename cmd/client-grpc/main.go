package main

import (
	"context"
	"fmt"
	pb "github.com/jxlwqq/grpc-lb/api/protobuf/pod"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"time"
)

const (
	waitDuration = 1 * time.Second
)

var counter int

var host = flag.String("host", "", "grpc host port")

func callAndShowResponse(client pb.PodClient) {
	resp, err := client.GetInfo(context.Background(), &pb.Request{})
	if err != nil {
		panic(err)
	}

	counter++
	fmt.Printf("#%d: %s\n", counter, resp.Name)
}

func main() {
	flag.Parse()
	fmt.Printf("Connecting to %s...\n", *host)
	conn, err := grpc.Dial(*host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("Making rpc...\n")
	client := pb.NewPodClient(conn)
	for {
		callAndShowResponse(client)
		time.Sleep(waitDuration)
	}
}
