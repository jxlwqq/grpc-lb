package pod

import (
	"context"
	"github.com/jxlwqq/grpc-lb/api/protobuf"
)

type server struct {
	protobuf.UnimplementedPodServer
	podName string
}

func NewServer(podName string) protobuf.PodServer {
	return &server{
		podName: podName,
	}
}

func (s *server) GetInfo(ctx context.Context, req *protobuf.Request) (*protobuf.Response, error) {
	return &protobuf.Response{
		Name: s.podName,
	}, nil
}
