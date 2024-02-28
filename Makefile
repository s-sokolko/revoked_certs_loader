SHELL := /bin/bash
container_name=revoked_certs_loader


.PHONY: build

all: build


build:
	go build -o load_revoked cmd/revoked_certs_loader/main.go

run:
	go run cmd/revoked_certs_loader/main.go

container:
	sudo docker build -t $(container_name) --no-cache -f build/Dockerfile .

container_push:
	sudo docker login
	sudo docker push $(container_name)

