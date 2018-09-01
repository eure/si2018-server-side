.DEFAULT_GOAL := help

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+\%*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# go-swagger

init: ## init ( "dep ensure" etc... )
	dep ensure

generate: generate-server generate-client ## generate everything

generate-server: ## generate server by swagger
	swagger generate server -f ./si2018.yml -a summerintern

generate-client: ## generate client by swagger
	swagger generate client -f ./si2018.yml -a summerintern

build: generate ## build server
	go build ./cmd/summer-intern2018-server/main.go

run: build ## run server
	./main --port=8888

# migrations

init-db: ## drop & create summer_intern_2018 database
	mysql -u ${LOCAL_DB_USERNAME} -h ${LOCAL_DB_HOSTNAME} -e "drop database if exists ${LOCAL_DB_DBNAME};"
	mysql -u ${LOCAL_DB_USERNAME} -h ${LOCAL_DB_HOSTNAME} -e "create database if not exists ${LOCAL_DB_DBNAME};"

migrate: ## just wrapping goose up
	goose up

dummy: ## insert dummy data
	go run misc/dummy/*.go

setup-db: init-db migrate dummy ## recreate db and insert dummy data
