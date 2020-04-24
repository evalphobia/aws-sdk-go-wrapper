.PHONY: init lint test-coverage send-coverage __setup_test

GO111MODULE=on

init:
	go mod download

lint:
	@type golangci-lint > /dev/null || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

test-coverage: __setup_test
	go test -race -covermode atomic -coverprofile=gotest.cov ./...

send-coverage:
	@type goveralls > /dev/null || go get -u github.com/mattn/goveralls
	goveralls -coverprofile=gotest.cov -service=github

__setup_test:
	@mkdir -p ./.local/s3_data
	@chmod 777 ./.local/s3_data
	docker-compose up -d
