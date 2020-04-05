.PHONY: build clean deploy

build:
	env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o bin/lambda cmd/lambda/main.go 	
	env go build -ldflags="-s -w" -o bin/server cmd/server/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
