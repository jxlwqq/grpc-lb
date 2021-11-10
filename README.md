# Kubernetes 中的 gRPC 负载均衡

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
| `make init`  | 安装 protoc-gen-go 和 protoc-gen-grpc |
| `make protoc`  | 基于 proto 文件，生成 *_pb.go 和 *_grpc.pb.go |
| `make docker-build`  | 构建 docker 镜像 |
| `make kube-deploy` | 在集群中部署服务 |
| `make kube-delete` | 删除服务 |
| `make istio-inject` | 注入 Istio 边车 |

具体逻辑，请查看 Makefile 文件。


### L4 vs L7 负载均衡

所谓的四层就是基于 IP + 端口的负载均衡，而七层就是基于 URL 等应用层信息的负载均衡； Kubernetes 内置的 Service 负载均衡基于 iptables/ipvs 实现，仅支持 L4。换句话说， Service 支持 HTTP/1.1 协议，不支持 HTTP/2 协议。 

而 Envoy(Istio) 则更为全能，支持被 gRPC 请求和响应的作为路由和负载均衡底层的所有 HTTP/2 功能。

### 项目架构

本项目分别测试 Service 和 Envoy(Istio) 对 HTTP/RPC 负载均衡的支持情况。

* cmd/server/main.go: 服务端，同时提供 HTTP 和 RPC 服务。响应的数据为服务端容器所在的 Pod 名称，（基于 [Downward API](https://kubernetes.io/zh/docs/tasks/inject-data-application/environment-variable-expose-pod-information/)）。
* cmd/client-http/main.go: HTTP 客户端，通过 HTTP 方式，循环调用服务端接口，并打印返回值。
* cmd/client-grpc/main.go: gRPC 客户端，通过 RPC 方式，循环远程调用服务端方法，并打印返回值。

### 测试原理

服务端 server 在 Kubernetes 集群中以 Deployment 的方式部署 3 个副本，3 个副本的 Pod 名称各不相同，而 client-http 和 client-grpc 则会每秒调用一次服务端，并打印返回值。如果返回值中，三个 Pod 的名称都存在，则表明正在进行有效的负载均衡，否则，则表明未进行有效的负载均衡。

### 测试 Service

构建镜像并部署在集群中：

```shell
make docker-build # 构建镜像
make kube-deploy  # 在集群中部署服务
```

详细逻辑请查看 Makefile 文件。

查看 Pod：

```shell
kubectl get pods
```

返回：

```shell
NAME                           READY   STATUS    RESTARTS   AGE
client-grpc-6c565594f4-tdf75   1/1     Running   0          2m48s
client-http-55d95c744d-f7nx4   1/1     Running   0          2m49s
server-7c4bfd74d-29c69         1/1     Running   0          2m51s
server-7c4bfd74d-4btvw         1/1     Running   0          2m51s
server-7c4bfd74d-fk8zf         1/1     Running   0          2m51s
```

查看 client-http Pod 的日志：

> Pod 名称的后缀部分是随机生成的，请替换为你自己的。

```shell
kubectl logs client-http-55d95c744d-f7nx4
```

返回：
```shell
#1: server-7c4bfd74d-4btvw
#2: server-7c4bfd74d-4btvw
#3: server-7c4bfd74d-29c69
#4: server-7c4bfd74d-fk8zf
#5: server-7c4bfd74d-fk8zf
#6: server-7c4bfd74d-29c69
#7: server-7c4bfd74d-fk8zf
#8: server-7c4bfd74d-4btvw
#9: server-7c4bfd74d-fk8zf
```

查看 client-grpc Pod 的日志：

```shell
kubectl logs client-grpc-6c565594f4-tdf75
```

返回：

```shell
#1: server-7c4bfd74d-fk8zf
#2: server-7c4bfd74d-fk8zf
#3: server-7c4bfd74d-fk8zf
#4: server-7c4bfd74d-fk8zf
#5: server-7c4bfd74d-fk8zf
#6: server-7c4bfd74d-fk8zf
#7: server-7c4bfd74d-fk8zf
#8: server-7c4bfd74d-fk8zf
#9: server-7c4bfd74d-fk8zf
```


可以看出，HTTP 请求在进行有效负载，而 RPC 请求在进行无效负载。

### 测试 Envoy(Istio)

我们在集群中已经部署了一个 Istio，但是没有设置自动注入的命令空间，所以我们在这里进行手动注入。

手动注入：

```shell
make istio-inject # 注入 Istio 边车
```

详细逻辑请查看 Makefile 文件。

查看 Pod：

```shell
kubectl get pods
```

返回：

```shell
NAME                           READY   STATUS    RESTARTS   AGE
client-grpc-7864f57779-f6blx   2/2     Running   0          17s
client-http-f8964854c-jclkd    2/2     Running   0          21s
server-7846bd6bb4-bcfws        2/2     Running   0          27s
server-7846bd6bb4-fv29s        2/2     Running   0          40s
server-7846bd6bb4-hzqj6        2/2     Running   0          34s
```

查看 client-http Pod 的日志：

```shell
kubectl logs client-http-f8964854c-jclkd
```

返回：

```shell
#1: server-7846bd6bb4-hzqj6
#2: server-7846bd6bb4-fv29s
#3: server-7846bd6bb4-hzqj6
#4: server-7846bd6bb4-hzqj6
#5: server-7846bd6bb4-hzqj6
#6: server-7846bd6bb4-hzqj6
#7: server-7846bd6bb4-hzqj6
#8: server-7846bd6bb4-bcfws
#9: server-7846bd6bb4-fv29s
```

查看 client-grpc Pod 的日志：
```shell
kubectl logs client-grpc-7864f57779-f6blx
```

返回：

```shell
#1: server-7846bd6bb4-fv29s
#2: server-7846bd6bb4-hzqj6
#3: server-7846bd6bb4-fv29s
#4: server-7846bd6bb4-bcfws
#5: server-7846bd6bb4-fv29s
#6: server-7846bd6bb4-hzqj6
#7: server-7846bd6bb4-fv29s
#8: server-7846bd6bb4-bcfws
#9: server-7846bd6bb4-fv29s
```

可以看出，HTTP 请求 和 RPC 请求均在进行有效负载。


### 清理

```shell
make kube-delete
istioctl x uninstall --purge
```

