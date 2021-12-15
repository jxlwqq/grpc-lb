.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: protoc
protoc:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	@ if ! which protoc-gen-go > /dev/null; then \
		echo "error: protoc-gen-go not installed" >&2; \
		exit 1; \
	fi
	@ if ! which protoc-gen-go-grpc > /dev/null; then \
		echo "error: protoc-gen-go-grpc not installed" >&2; \
		exit 1; \
	fi
	for file in $$(git ls-files '*.proto'); do \
		protoc -I $$(dirname $$file) \
		--go_out=:$$(dirname $$file) --go_opt=paths=source_relative \
		--go-grpc_out=:$$(dirname $$file) --go-grpc_opt=paths=source_relative \
		$$file; \
	done

.PHONY: docker-build
docker-build:
	docker build -t grpc-lb/server -f ./build/docker/server/Dockerfile .
	docker build -t grpc-lb/client-http -f ./build/docker/client-http/Dockerfile .
	docker build -t grpc-lb/client-grpc -f ./build/docker/client-grpc/Dockerfile .

.PHONY: kube-deploy
kube-deploy:
	kubectl apply -f deployments/server.yaml
	kubectl rollout status deployment/server
	kubectl apply -f deployments/client-http.yaml
	kubectl rollout status deployment/client-http
	kubectl apply -f deployments/client-grpc.yaml
	kubectl rollout status deployment/client-grpc
	kubectl get pods

.PHONY: kube-delete
kube-delete:
	kubectl delete -f deployments

.PHONY: istio-inject
istio-inject:
	istioctl kube-inject -f deployments/server.yaml | kubectl apply -f -
	kubectl rollout status deployment/server
	istioctl kube-inject -f deployments/client-http.yaml | kubectl apply -f -
	kubectl rollout status deployment/client-http
	istioctl kube-inject -f deployments/client-grpc.yaml | kubectl apply -f -
	kubectl rollout status deployment/client-grpc
	kubectl get pods