.PHONY: build clean deploy

generate:
	cd graphql/self-service && go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen

client:
	cd service/shopify && go get github.com/Khan/genqlient/generate && go run -v github.com/Khan/genqlient
	go mod tidy

all: generate client

build:
	echo "building backend ...."
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/bootstrap main.go
	zip -j bin/backoffice.zip bin/bootstrap

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build_backoffice
	npx serverless deploy --stage production --verbose
