all:
	go build -o client.bin ./client/main.go
	go build -o node.bin ./node/main.go
