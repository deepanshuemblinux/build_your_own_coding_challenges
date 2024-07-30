build:
	go build -o bin/lbapp .

lb: build
	bin/lbapp lb $(addr)

backend: build
	bin/lbapp backend $(addr)