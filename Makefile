
build:
	go build -o bin/hcl-generator cmd/hcl-generator/*.go
	go build -o bin/shutdown cmd/shutdown/*.go
	go build -o bin/delete-server cmd/delete-server/*.go
