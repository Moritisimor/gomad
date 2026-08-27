BINARY := gomad

.PHONY: all build test vet clean

all: build

build:
	go build -o $(BINARY) ./cmd/gomad

test:
	go test ./...

vet:
	go vet ./...

clean:
	$(RM) $(BINARY) interpreter.test coverage.out
