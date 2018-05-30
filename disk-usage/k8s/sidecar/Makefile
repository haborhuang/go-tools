PKG="github.com/haborhuang/go-tools/disk-usage/k8s/sidecar"
IMAGE?="du-mon-sidecar"

docker-build:
	docker run -i --rm -v `pwd`:/go/src/$(PKG) -w /go/src/$(PKG) golang:1.8.3-alpine3.6 go build -o sidecar

build-image: docker-build
	docker build -t $(IMAGE) .

.PHONY: docker-build build-image
