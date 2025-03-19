
build:
	go build ./...
run:
	for d in examples/*/*/main.go; do  echo "$$d"; go run ./$$d; done

test:
	go test ./...

