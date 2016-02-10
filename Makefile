SHORT_NAME := workflow-manager

# Enable vendor/ directory support.
export GO15VENDOREXPERIMENT=1

# SemVer with build information is defined in the SemVer 2 spec, but Docker
# doesn't allow +, so we use -.
VERSION := git-$(shell git rev-parse --short HEAD)

DEV_ENV_IMAGE := quay.io/deis/go-dev:0.9.0
DEV_ENV_WORK_DIR := /go/src/github.com/deis/${SHORT_NAME}
DEV_ENV_CMD := docker run --rm -e CGO_ENABLED=0 -v ${PWD}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${DEV_ENV_IMAGE}

# Common flags passed into Go's linker.
LDFLAGS := "-s -X main.version=${VERSION}"

# Docker Root FS
BINDIR := ./rootfs/bin

# Legacy support for DEV_REGISTRY, plus new support for DEIS_REGISTRY.
ifdef ${DEV_REGISTRY}
  DEIS_REGISTRY = ${DEV_REGISTRY}/
endif

IMAGE_PREFIX ?= deis/

# Docker image name
IMAGE := ${DEIS_REGISTRY}${IMAGE_PREFIX}${SHORT_NAME}:${VERSION}

all: build docker-build docker-push

# Containerized dependency resolution / initial workspace setup
bootstrap:
	${DEV_ENV_CMD} glide install

# This illustrates a two-stage Docker build. docker-compile runs inside of
# the Docker environment. Other alternatives are cross-compiling, doing
# the build as a `docker build`.
build:
	mkdir -p ${BINDIR}
	${DEV_ENV_CMD} go build -o ${BINDIR}/boot -a -installsuffix cgo -ldflags ${LDFLAGS} boot.go

test:
	${DEV_ENV_CMD} sh -c 'go test -v $$(glide nv)'

# For cases where we're building from local
# We also alter the RC file to set the image name.
docker-build:
	docker build --rm -t ${IMAGE} rootfs

# Push to a registry that Kubernetes can access.
docker-push:
	docker push ${IMAGE}
