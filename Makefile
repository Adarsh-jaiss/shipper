.PHONY: server 

server:
	@echo "Building server..."
	@go build -o bin/server ./server
	@echo "Setting executable permissions..."
	@chmod +x bin/server
	@echo "Running server..."
	@./bin/server