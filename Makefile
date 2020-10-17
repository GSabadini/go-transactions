PWD = $(shell pwd -L)
IMAGE_NAME = gsabadini/go-transactions
DOCKER_RUN = docker run --rm -it -v ${PWD}:/app -w /app ${IMAGE_NAME}

init:
	cp .env.example .env

fmt:
	go fmt ./...

test:
	go test ./...

up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans

logs:
	docker-compose logs -f app

build-prod:
	docker build -t ${IMAGE_NAME} -f Dockerfile .

ci:
	docker run --rm -it -v ${PWD}:/app -w /app golangci/golangci-lint:v1.24.0 golangci-lint run ./...
#	--exclude-use-default=false \
#	--enable=golint \
#	--enable=gocyclo \
#	--enable=goconst \
#	--enable=unconvert \
#	./...