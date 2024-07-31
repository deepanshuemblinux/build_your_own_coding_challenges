build:
	go build -o bin/go-redis .

server: build
	./bin/go-redis server

client: build
	./bin/go-redis client