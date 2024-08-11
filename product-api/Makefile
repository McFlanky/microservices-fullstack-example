build:
	@go build -o bin/api

run: build
	@./bin/api

check-install:
	@which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger 

swagger: check-install
	@swagger generate spec -o ./swagger.yaml --scan-models

generate_client:
	@cd sdk && swagger generate client -f ../swagger.yaml -A product-api
