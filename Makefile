.PHONY: build clean deploy

build:
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o cmd/lambda/lambda bin/lambda
	env go build -ldflags="-s -w" -o cmd/server/server bin/server

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
