all: build

test:
	go test -v -race -cover .

build:
	go build -o target/sync-mvn-deps ./*.go

clean:
	@rm -rf target