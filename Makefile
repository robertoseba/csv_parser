
default:
	@echo "Please use one of the following commands:"
	@echo "make test"
	@echo "make build"
	@echo "make run"

tests:
	go test  ./... 

compile:
	go build -o ./build/

run:
	go run . 

view-profile:
	go tool pprof cpu.pprof

bench: 
	go test -bench=. -benchmem ./...