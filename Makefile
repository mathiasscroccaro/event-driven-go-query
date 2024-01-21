format:
	go fmt ./...

test:
	go test -cover ./...

build_worker:
	go build -o worker cmd/worker/main.go

build_api:
	go build -o api cmd/api/main.go

build: build_worker build_api

clean:
	rm -rf worker
	rm -rf api