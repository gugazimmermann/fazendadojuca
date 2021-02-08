.PHONY: build clean deploy remove

build:
	go get -u ./...
	env GOOS=linux GOARCH=amd64 go build -o bin/animals animals/main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/breed breed/main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/gender gender/main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/purity_level purity_level/main.go

clean:
	rm -rf ./bin ./vendor ./.serverless Gopkg.lock

deploy: clean build
	sls deploy --verbose

remove: clean
	sls remove