SHELL := /bin/bash

codegen:
	oapi-codegen -generate types,client,spec -o internal/generated/api/api.go -package genrated_api ./api/api.yaml
	go generate ./...

run-dev:
	sed -i 's/^DATABASE_ADDRESS=.*/DATABASE_ADDRESS=postgres:5432/' $(CURDIR)/config.env
	set -a; source $(CURDIR)/config.env; set +a; \
	docker-compose up

run-postgres:
	sed -i 's/^DATABASE_ADDRESS=.*/DATABASE_ADDRESS=localhost:5432/' $(CURDIR)/config.env
	set -a; source $(CURDIR)/config.env; set +a; \
	docker-compose up postgres

