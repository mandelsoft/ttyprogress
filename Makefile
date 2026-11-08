
build:
	go build ./...
run:
	for d in examples/*/*/main.go; do  echo "$$d"; go run ./$$d; done

test:
	go test ./...


# The behaviour of the examples is recorded and checked with ttyframes
# (github.com/mandelsoft/ttyframes, pinned by tools.go).
TTYFRAMES := go run github.com/mandelsoft/ttyframes/cmd/ttyframes

# behaviour records the behaviour of the examples (files expected.ttyframes.json).
behaviour:
	$(TTYFRAMES) record ./examples/...

# check-behaviour compares the behaviour of the examples with the recorded one.
check-behaviour:
	$(TTYFRAMES) check ./examples/...

# The README.md is generated with mdref (github.com/mandelsoft/mdref,
# pinned by tools.go) from doc/src/README.md, which includes the code
# fragments of the examples marked by "--- begin <key> ---" and
# "--- end <key> ---" comments.
MDREF := go run github.com/mandelsoft/mdref

# doc generates the README.md.
doc:
	$(MDREF) --headings --unix doc/src .
	chmod 644 README.md

# check-doc checks, that the README.md is up to date with its source
# and the examples. The generation header (the first three lines)
# depends on the target directory and is not compared.
check-doc:
	@tmp=$$(mktemp -d); trap 'rm -rf "$$tmp"' EXIT; \
	$(MDREF) --headings --unix doc/src "$$tmp" && \
	tail -n +4 README.md > "$$tmp/actual" && \
	tail -n +4 "$$tmp/README.md" > "$$tmp/expected" && \
	diff "$$tmp/actual" "$$tmp/expected" || \
	{ echo "README.md is not up to date, run make doc"; exit 1; }

.PHONY: build run test behaviour check-behaviour doc check-doc
