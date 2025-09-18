.PHONY: build clean deploy

generate:
	go get github.com/99designs/gqlgen
	cd graphql/self-service && go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen
	cd graphql/self-queue && go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen
	cd graphql/backoffice && go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen

client:
	go get github.com/Khan/genqlient/generate && cd service/shopify && go run -v github.com/Khan/genqlient && cd ../custom && go run -v github.com/Khan/genqlient
	go mod tidy

all: generate client

build:
	echo "building backend ...."
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/bootstrap main.go
	zip -j bin/backoffice.zip bin/bootstrap

clean:
	rm -rf ./bin ./vendor Gopkg.lock
