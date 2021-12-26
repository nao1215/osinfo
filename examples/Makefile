build: deps 
	go build -o example example.go

run: build
	./example

clean: 
	rm example
	rm -rf vendor

deps: 
	go mod vendor