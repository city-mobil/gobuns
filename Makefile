.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: env_up
env_up:
	docker-compose up -d

.PHONY: env_down
env_down:
	docker-compose down -v --rmi local --remove-orphans

.PHONY: lint
lint:
	golangci-lint run -v ./...

.PHONY: test
test: env_up
	go mod tidy
	GOPATH=`go env GOPATH` docker-compose -f docker-compose.yml -f docker-compose.ci.yml run tests

.PHONY: gen
gen:
	go generate ./...