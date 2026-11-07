
build:
	go build ./...
run:
	for d in examples/*/*/main.go; do  echo "$$d"; go run ./$$d; done

test:
	go test ./...


# The behaviour of the examples is recorded and checked with ttyframes
# (github.com/mandelsoft/ttyframes, pinned by tools.go).
TTYFRAMES := go run github.com/mandelsoft/ttyframes/cmd/ttyframes

# behaviour records the behaviour of the examples (files expected.ttyframes).
behaviour:
	$(TTYFRAMES) record ./examples/...

# check-behaviour compares the behaviour of the examples with the recorded one.
check-behaviour:
	$(TTYFRAMES) check ./examples/...

.PHONY: build run test behaviour check-behaviour
