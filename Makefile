
default:
	@echo "Please use one of the following commands:"
	@echo "make test"
	@echo "make build"
	@echo "make run"

tests:
	go test  ./... 

build:
	go build -o ./build/ ./..

run:
	go run ./..

generate-test-csv:
	go run . --generate-csv 3 100
