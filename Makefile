DSRCS=$(wildcard *.go) go.mod
SRCS=$(wildcard cmd/incdb/main.go) go.mod

all: incdb incdbd data

incdb: $(SRCS)
	go build -o incdb cmd/incdb/*.go

incdbd: $(DSRCS)
	go build -o incdbd *.go

data:
	mkdir -p data

test: all
	go test ./...

testv: all
	go test -v ./...

clean:
	rm -f incdb incdbd data/*

.PHONY: test testv clean
