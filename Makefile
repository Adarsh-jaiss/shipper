.PHONY: build run infra push

# Build the main application
build:
	@go build -o bin/shipper

# Run the built application
run: build
	@./bin/shipper

# Build and run infrastructure-related code


# Push changes to git
push:
	@git init
	@git add .
	@git commit -s -m "$(msg)"
	@git push origin main
