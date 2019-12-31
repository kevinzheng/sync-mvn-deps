all: clean test build install

test:
	go test -v -race -cover .

build:
	go build -o target/sync-mvn-deps ./*.go

clean:
	rm -rf target

install:
	cp target/sync-mvn-deps /usr/local/bin/