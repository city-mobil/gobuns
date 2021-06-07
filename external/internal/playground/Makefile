build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ping/bin/ping ping/ping.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o pong/bin/pong pong/pong.go
	docker build -t go-buns/ping ./ping
	docker build -t go-buns/pong ./pong

up: build
	docker-compose up

down:
	docker-compose down -v