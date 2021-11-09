# Kubernetes 中的 gRPC 负载平衡

### 安装环境依赖

* docker-desktop >= 4.1.1
* kubernetes >= 1.21.5
* go >= 1.17
* protobuf >= 3.17.3
* istioctl >= 1.11.4

下载安装 Docker Desktop ，并启动内置的 Kubernetes 集群。

```shell
# 安装 Go
brew install go
# 安装 Protobuf
brew install protobuf
# 安装 Istio
brew install istioctl
kubectl config use-context docker-desktop
istioctl install -y
```

### Makefile 介绍

|  命令   | 说明  |
|  ----  | ----  |
| `make init`  | 安装 protoc-gen-go 和 protoc-gen-grpc 二进制文件 |
| `make protoc`  | 基于 proto 文件，生成 pb.go 和 grpc.pb.go 文件 |
| `make docker-build`  | 构建镜像 |
| `make kube-deploy` | 在集群中部署服务 |
| `make kube-delete` | 删除服务 |
| `make istio-inject` | 注入 Istio 边车 |


### L4 vs L7 负载均衡

所谓四层就是基于 IP + 端口的负载均衡，七层就是基于 URL 等应用层信息的负载均衡；

### 项目架构

本项目分别测试 Service 和 Envoy(Istio) 对 HTTP/RPC 负载均衡的支持情况。

* cmd/server/main.go: 服务端，同时提供 HTTP 和 RPC 服务。响应的数据为服务端容器所在的 Pod 名称，基于 [Downward API](https://kubernetes.io/zh/docs/tasks/inject-data-application/environment-variable-expose-pod-information/)。
* cmd/client-http/main.go: HTTP 客户端，通过 HTTP 方式，循环调用服务端接口，并打印返回值。
* cmd/client-rpc/main.go: gRPC 客户端，通过 RPC 方式，循环远程调用服务端方法，并打印返回值。


