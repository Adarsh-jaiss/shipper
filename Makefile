build:
	@go build -o bin/shipper

run: build
	@./bin/shipper

push:
	@git init
	@git add .
	@git commit -s -m "$(msg)"
	@git push origin main
