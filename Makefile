API_IMAGE_TAG ?= latest
GO_MIGRATE_VERSION ?= v4.14.1


build-api:
	docker build \
		-f zarf/docker/api/Dockerfile \
		-t rebot-api-amd64:$(API_IMAGE_TAG) \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

build-commander:
	docker build \
		-f zarf/docker/commander/Dockerfile \
		-t rebot-commander-amd64:$(API_IMAGE_TAG) \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

build-weather:
	docker build \
		-f zarf/docker/weather/Dockerfile \
		-t rebot-weather-amd64:$(API_IMAGE_TAG) \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

build: build-api build-commander build-weather

install: build up wait-db ## Build containers and up all services

run: up wait-db ## Up all services

lint:
	golangci-lint run ./... --out-format code-climate

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
	docker-compose -p ${PROJECT_NAME} -f zarf/compose/docker-compose.yaml -f zarf/compose/config-compose.yaml run wait -c rebot-db:5432

create-migration: ## Create migration file in db/migrations directory. Migration should be named by "name" argument. Example: create-migration name=create_foos
	docker run -v "${PWD}/migration:/migration" --network host migrate/migrate:$(GO_MIGRATE_VERSION) -path=/migration \
 		create -ext sql -dir /migration $(name)

models:
	sqlc generate

models-wsl:
	wsl.exe sqlc generate

tidy:
	go mod tidy
	go mod vendor