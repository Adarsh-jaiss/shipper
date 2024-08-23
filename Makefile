build:
	@go build -o bin/shipper

run: build
	@./bin/shipper

infra:
	@go build -o bin/shipper-infra /infra
	@./bin/shipper-infra

push:
	@git init
	@git add .
	@git commit -s -m "$(msg)"
	@git push origin main

