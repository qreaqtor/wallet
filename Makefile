SHELL := /bin/bash

codegen:
	oapi-codegen -generate types,client,spec -o internal/generated/api/api.go -package genrated_api ./api/api.yaml
	go generate ./...

run-dev:
	set -a; source $(CURDIR)/config/dev.env; set +a; \
	docker-compose up

run-postgres:
	@set -e; \
	set -a; source $(CURDIR)/config/local.env; set +a; \
	cp $(CURDIR)/config/local.env $(CURDIR)/.env; \
	docker-compose up postgres
