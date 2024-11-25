# Project variables
PROJECT_NAME := chaos-kube
REGISTRY := pandects
VERSION := $(shell git describe --tags --always)
TIMESTAMP := $(shell date +%s)
# Services
SERVICES := chaos-executor chaos-controller #api chaos-configurator log-aggregator

# Docker build and push
.PHONY: build push deploy clean run

# Default target
all: build

# Build all services
build: $(patsubst %,build-%, $(SERVICES))

# Build each service
build-%:
	@echo "Building Docker image for $*..."
	docker build -t $(REGISTRY)/$*:$(TIMESTAMP) -f src/cmd/$*/Dockerfile src
#docker build -t pandects/chaos-executor:latest -f src/cmd/chaos-executor/Dockerfile src

# Push all services to registry
push: $(patsubst %,push-%, $(SERVICES))

# Push each service to Docker registry
push-%:
	@echo "Pushing $* to Docker registry..."
	docker push $(REGISTRY)/$*:latest

# Deploy all services to Kubernetes
deploy:
	helm install chaos-kube chaos-kube-chart -f chaos-kube-chart/values/values-dev.yaml --create-namespace || helm upgrade chaos-kube chaos-kube-chart -f chaos-kube-chart/values/values-dev.yaml


# Test all services
test: 
	@cd src && go test ./...

# Clean build artifacts and images
clean:
	@echo "Cleaning up..."
	rm -rf bin/*
	@for service in $(SERVICES); do \
		docker rmi $(REGISTRY)/$$service:$(VERSION) || true; \
	done

# Run target to test, build, and deploy all services
run: test build push deploy
	@echo "All services have been tested, built, and deployed."


# Example usage:
#   make build             - Build all services
#   make push              - Push all services to Docker registry
#   make deploy            - Deploy all services to Kubernetes
#   make build-api         - Build only the api service
#   make api               - Build, push, and deploy the api service
#   make clean             - Clean up all binaries and Docker images
