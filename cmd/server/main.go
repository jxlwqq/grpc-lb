package main

import (
	"fmt"
	pb "github.com/jxlwqq/grpc-lb/api/protobuf/pod"
	"github.com/jxlwqq/grpc-lb/internal/pod"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

const (
	httpPort = ":80"
	grpcPort = ":30051"
)

var podName = flag.String("pod-name", "", "pod name")

func main() {

	flag.Parse()

	go func() {
		log.Printf("http server listen on %s", httpPort)
		http.HandleFunc("/pod", func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, *podName)
		})

		if err := http.ListenAndServe(httpPort, nil); err != nil {
			log.Fatalf("http server failed to listen: %v", err)
		}
	}()

	log.Printf("grpc server listen on %s", grpcPort)
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("grpc server failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := pod.NewServer(*podName)

	pb.RegisterPodServer(s, server)

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("grpc server failed to listen: %v", err)
	}
}
