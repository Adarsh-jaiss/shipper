build:
	@go build -o bin/shipper

run: build
	@./bin/shipper
