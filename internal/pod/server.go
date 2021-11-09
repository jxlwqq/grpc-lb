package pod

import (
	"context"
	"github.com/jxlwqq/grpc-lb/api/protobuf/pod"
)

type server struct {
	pod.UnimplementedPodServer
	podName string
}

func NewServer(podName string) pod.PodServer {
	return &server{
		podName: podName,
	}
}

func (s *server) GetInfo(ctx context.Context, req *pod.Request) (*pod.Response, error) {
	return &pod.Response{
		Name: s.podName,
	}, nil
}
