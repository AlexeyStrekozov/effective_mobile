include .env
export 

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d effective_mobile-postgres

env-down:
	@docker compose down effective_mobile-postgres

env-cleanup:
	@read -p "Clear all environment files? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down effective_mobile-postgres port-forwarder && \
		rm -rf {PROJECT_ROOT}/out/pgdata && \
		echo "Environment files have been cleared"; \
	else \
		echo "Cleaning up the environment of the canceled"; \
	fi;

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing parameter 'seq'. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm effective_mobile-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing parameter 'action'. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm effective_mobile-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@effective_mobile-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"${action}"

logs-cleanup:
	@read -p "Clear all logs files? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Logs files have been cleared"; \
	else \
		echo "Cleaning up the logs of the canceled"; \
	fi;

dev-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/app/main.go

app-deploy:
	@docker compose up -d --build app

app-undeploy:
	@docker compose down app

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g /cmd/app/main.go
		-o docs \
		--parseInternal \
		--parseDependency

ps:
	@docker compose ps