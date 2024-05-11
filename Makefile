ifndef DOCKER_IMAGE_NAME
DOCKER_IMAGE_NAME=simple-telemetry-publisher
endif

ifndef DOCKER_IMAGE_TAG
DOCKER_IMAGE_TAG=0.0.1
endif

docker-image:
	@echo "Building docker image..."
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

docker-save:
	@echo "Saving docker image..."
	docker save $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -o $(DOCKER_IMAGE_NAME)_$(DOCKER_IMAGE_TAG).tar

compile:
	@echo "Compile binary..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o simple-telemetry-publisher ./cmd/simple-telemetry-publisher/main.go

run:
	@echo "Running binary..."
	./simple-telemetry-publisher --config ./fixtures/simple-publisher-config.yaml

test:
	@echo "Running tests..."
	go test -v ./...
