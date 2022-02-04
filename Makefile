BINARY_NAME=price-server

build:
	go build -o ${BINARY_NAME} main.go

run:
	./${BINARY_NAME}

build_and_run: build run

test: 
	go test -v ./...

clean:
	go clean
	rm ${BINARY_NAME}