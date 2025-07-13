.PHONY: run build test-docs

run: build
	./target/main -p https://demo.unified-streaming.com/k8s/features/stable/video/tears-of-steel/tears-of-steel.ism/.mpd

build:
	go build -o target/main .

test-docs:
	go doc -all

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=tmp/coverage.out ./...
	go tool cover -html=tmp/coverage.out