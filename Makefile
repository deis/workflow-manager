SHORT_NAME := workflow-manager

# Enable vendor/ directory support.
export GO15VENDOREXPERIMENT=1

include versioning.mk

DEV_ENV_IMAGE := quay.io/deis/go-dev:0.9.0
DEV_ENV_WORK_DIR := /go/src/github.com/deis/${SHORT_NAME}
DEV_ENV_CMD := docker run --rm -e CGO_ENABLED=0 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}

# Common flags passed into Go's linker.
LDFLAGS := "-s -X main.version=${VERSION}"

# Docker Root FS
BINDIR := ${CURDIR}/rootfs/bin

# Legacy support for DEV_REGISTRY, plus new support for DEIS_REGISTRY.
ifdef ${DEV_REGISTRY}
  DEIS_REGISTRY = ${DEV_REGISTRY}/
endif

all: build docker-build docker-push

# Containerized dependency resolution / initial workspace setup
bootstrap:
	${DEV_ENV_CMD} glide install

# This illustrates a two-stage Docker build. docker-compile runs inside of
# the Docker environment. Other alternatives are cross-compiling, doing
# the build as a `docker build`.
build:
	mkdir -p ${BINDIR}
	${DEV_ENV_CMD} go build -o rootfs/bin/boot -a -installsuffix cgo -ldflags ${LDFLAGS} boot.go

test:
	${DEV_ENV_CMD} sh -c 'go test -v $$(glide nv)'

# For cases where we're building from local
# We also alter the RC file to set the image name.
docker-build:
	docker build --rm -t ${IMAGE} rootfs
	docker tag -f ${IMAGE} ${MUTABLE_IMAGE}
