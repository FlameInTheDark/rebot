API_IMAGE_TAG ?= latest
GO_MIGRATE_VERSION ?= v4.14.1

build: build-api

build-api:
	docker build \
		-f zarf/docker/api/Dockerfile \
		-t rebot-api-amd64:$(API_IMAGE_TAG) \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

install: build up wait-db ## Build containers and up all services

run: up wait-db ## Up all services

PROJECT_NAME="rebot"

up:
	docker-compose -p ${PROJECT_NAME} -f zarf/compose/docker-compose.yaml -f zarf/compose/config-compose.yaml up --detach --remove-orphans

down:
	docker-compose -p ${PROJECT_NAME} -f zarf/compose/docker-compose.yaml down --remove-orphans

stop:
	docker-compose -p ${PROJECT_NAME} -f zarf/compose/docker-compose.yaml stop

logs:
	docker-compose -p ${PROJECT_NAME} -f zarf/compose/docker-compose.yaml logs -f

wait-db:
	docker-compose -p ${PROJECT_NAME} --f zarf/compose/docker-compose.yaml -f zarf/compose/config-compose.yaml run wait -c rebot-db:5432

create-migration: ## Create migration file in db/migrations directory. Migration should be named by "name" argument. Example: create-migration name=create_foos
	docker run -v "${PWD}/migrations:/migrations" --network host migrate/migrate:$(GO_MIGRATE_VERSION) -path=/migrations \
 		create -ext sql -dir /migrations $(name)