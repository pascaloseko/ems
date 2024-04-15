.PHONY: generate server

generate:
	go get -d github.com/99designs/gqlgen && go run github.com/99designs/gqlgen generate

server:
	go run server.go