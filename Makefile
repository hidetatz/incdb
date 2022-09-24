SRCS=$(wildcard *.go)

incdb: $(SRCS)
	go build -o incdb *.go

test: incdb
	go test ./...

clean:
	rm -f incdb data/incdb.data data/incdb.catalog

.PHONY: test clean
