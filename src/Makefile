# Project variables
PROJECT_NAME := chaos-kube
REGISTRY := my-docker-registry
VERSION := $(shell git describe --tags --always)

# Services
SERVICES := api chaos-controller chaos-configurator chaos-executor log-aggregator

# Kubernetes namespace for deployment
NAMESPACE := chaos-kube

# Docker build and push
.PHONY: build push deploy clean run

# Default target
all: build

# Build all services
build: $(patsubst %,build-%, $(SERVICES))

# Build each service
build-%:
	@echo "Building service $*..."
	@cd cmd/$* && GOOS=linux GOARCH=amd64 go build -o ../../bin/$* main.go
	@echo "Building Docker image for $*..."
	docker build -t $(REGISTRY)/$(PROJECT_NAME)-$*:$(VERSION) -f cmd/$*/Dockerfile .

# Push all services to registry
push: $(patsubst %,push-%, $(SERVICES))

# Push each service to Docker registry
push-%:
	@echo "Pushing $* to Docker registry..."
	docker push $(REGISTRY)/$(PROJECT_NAME)-$*:$(VERSION)

# Deploy all services to Kubernetes
deploy: $(patsubst %,deploy-%, $(SERVICES))

# Deploy each service using kubectl apply or Helm
deploy-%:
	@echo "Deploying $* to Kubernetes..."
	kubectl apply -f deployments/$*.yaml -n $(NAMESPACE)

# Test all services
test: 
	go test ./...

# Run tests for each service
test-%:
	@echo "Running tests for $*..."
	cd $* && go test ./...

# Clean build artifacts and images
clean:
	@echo "Cleaning up..."
	rm -rf bin/*
	@for service in $(SERVICES); do \
		docker rmi $(REGISTRY)/$(PROJECT_NAME)-$$service:$(VERSION) || true; \
	done

# Individual service targets for convenience
# Build, push, and deploy a single service
%: test-% build-% push-% deploy-%
	@echo "Service $* has been tested, built, pushed, and deployed."

# Run target to test, build, and deploy all services
run: test build deploy
	@echo "All services have been tested, built, and deployed."


# Example usage:
#   make build             - Build all services
#   make push              - Push all services to Docker registry
#   make deploy            - Deploy all services to Kubernetes
#   make build-api         - Build only the api service
#   make api               - Build, push, and deploy the api service
#   make clean             - Clean up all binaries and Docker images
