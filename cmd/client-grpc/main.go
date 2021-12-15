package main

import (
	"context"
	"fmt"
	"github.com/jxlwqq/grpc-lb/api/protobuf"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"time"
)

const (
	waitDuration = 1 * time.Second
)

var counter int

var host = flag.String("host", "", "grpc host port")

func callAndShowResponse(client protobuf.PodClient) {
	resp, err := client.GetInfo(context.Background(), &protobuf.Request{})
	if err != nil {
		panic(err)
	}

	counter++
	fmt.Printf("#%d: %s\n", counter, resp.Name)
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := protobuf.NewPodClient(conn)
	for {
		callAndShowResponse(client)
		time.Sleep(waitDuration)
	}
}
