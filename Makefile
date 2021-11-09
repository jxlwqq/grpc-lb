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
	docker build -t jxlwqq/server -f ./cmd/server/Dockerfile .
	docker build -t jxlwqq/client-http -f ./cmd/client-http/Dockerfile .
	docker build -t jxlwqq/client-grpc -f ./cmd/client-grpc/Dockerfile .

.PHONY: kube-deploy
kube-deploy:
	kubectl apply -f manifests

.PHONY: kube-delete
kube-delete:
	kubectl delete -f manifests

.PHONY: kube-inject
kube-inject:
	istioctl kube-inject -f manifests/server.yaml | kubectl apply -f -
	istioctl kube-inject -f manifests/client-http.yaml | kubectl apply -f -
	istioctl kube-inject -f manifests/client-grpc.yaml | kubectl apply -f -
