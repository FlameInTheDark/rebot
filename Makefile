API_IMAGE_TAG ?= latest
GO_MIGRATE_VERSION ?= v4.14.1


build-api:
	docker build \
		-f zarf/docker/api/Dockerfile \
		-t flameinthedark/rebot-api:$(API_IMAGE_TAG) \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

build-commander:
	docker build \
		-f zarf/docker/commander/Dockerfile \
		-t flameinthedark/rebot-commander:$(API_IMAGE_TAG) \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

build-weather:
	docker build \
		-f zarf/docker/weather/Dockerfile \
		-t flameinthedark/rebot-weather:$(API_IMAGE_TAG) \
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

vendor:
	go mod vendor

## Kubernetes deployment
kube-up: k8s-app-secrets-up k8s-database-deploy-up k8s-redis-deploy-up k8s-rabbitmq-deploy-up k8s-consul-deploy-up k8s-influx-deploy-up k8s-api-deploy-up k8s-commander-deploy-up k8s-weather-deploy-up
kube-down: k8s-app-secrets-down k8s-database-deploy-down k8s-redis-deploy-down k8s-rabbitmq-deploy-down k8s-consul-deploy-down k8s-influx-deploy-down k8s-api-deploy-down k8s-commander-deploy-down k8s-weather-deploy-down

## Up
k8s-api-deploy-up:
	kubectl apply -f zarf/k8s/api/api-app-deployment.yaml

k8s-commander-deploy-up:
	kubectl apply -f zarf/k8s/commander/commander-app-deployment.yaml

k8s-weather-deploy-up:
	kubectl apply -f zarf/k8s/commander/commander-app-deployment.yaml

k8s-database-deploy-up:
	kubectl apply -f zarf/k8s/database/database-secret.yaml
	kubectl apply -f zarf/k8s/database/database-configmap.yaml
	kubectl apply -f zarf/k8s/database/database-pv.yaml
	kubectl apply -f zarf/k8s/database/database-app-deployment.yaml

k8s-redis-deploy-up:
	kubectl apply -f zarf/k8s/redis/redis-secret.yaml
	kubectl apply -f zarf/k8s/redis/redis-configmap.yaml
	kubectl apply -f zarf/k8s/redis/redis-app-deployment.yaml

k8s-rabbitmq-deploy-up:
	kubectl apply -f zarf/k8s/rabbitmq/rabbitmq-secret.yaml
	kubectl apply -f zarf/k8s/rabbitmq/rabbitmq-configmap.yaml
	kubectl apply -f zarf/k8s/rabbitmq/rabbitmq-app-deployment.yaml

k8s-consul-deploy-up:
	kubectl apply -f zarf/k8s/consul/consul-configmap.yaml
	kubectl apply -f zarf/k8s/consul/consul-app-deployment.yaml

k8s-influx-deploy-up:
	kubectl apply -f zarf/k8s/influx/influx-configmap.yaml
	kubectl apply -f zarf/k8s/influx/influx-secret.yaml

k8s-app-secrets-up:
	kubectl apply -f zarf/k8s/secrets/rebot-secrets.yaml

## Down
k8s-api-deploy-down:
	kubectl delete -f zarf/k8s/api/api-app-deployment.yaml

k8s-commander-deploy-down:
	kubectl delete -f zarf/k8s/commander/commander-app-deployment.yaml

k8s-weather-deploy-down:
	kubectl delete -f zarf/k8s/weather/weather-app-deployment.yaml

k8s-database-deploy-down:
	kubectl delete -f zarf/k8s/database/database-secret.yaml
	kubectl delete -f zarf/k8s/database/database-configmap.yaml
	kubectl delete -f zarf/k8s/database/database-pv.yaml
	kubectl delete -f zarf/k8s/database/database-app-deployment.yaml

k8s-redis-deploy-down:
	kubectl delete -f zarf/k8s/redis/redis-secret.yaml
	kubectl delete -f zarf/k8s/redis/redis-configmap.yaml
	kubectl delete -f zarf/k8s/redis/redis-app-deployment.yaml

k8s-rabbitmq-deploy-down:
	kubectl delete -f zarf/k8s/rabbitmq/rabbitmq-secret.yaml
	kubectl delete -f zarf/k8s/rabbitmq/rabbitmq-configmap.yaml
	kubectl delete -f zarf/k8s/rabbitmq/rabbitmq-app-deployment.yaml

k8s-influx-deploy-down:
	kubectl delete -f zarf/k8s/influx/influx-configmap.yaml
	kubectl delete -f zarf/k8s/influx/influx-secret.yaml

k8s-consul-deploy-down:
	kubectl delete -f zarf/k8s/consul/consul-configmap.yaml
	kubectl delete -f zarf/k8s/consul/consul-app-deployment.yaml

k8s-app-secrets-down:
	kubectl delete -f zarf/k8s/secrets/rebot-secrets.yaml