include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up -d todoapp-postgres

env-down:
	docker compose down todoapp-postgres

env-cleanup:
	@read -p "Delete all volume files? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  	docker compose down todoapp-postgres port-forwarder && \
	  	rm -rf ${PROJECT_ROOT}/out/pgdata
	  	echo "Volume files has been deleted"; \
	else \
	  echo "Removing volume files has been canceled"; \
	  fi

env-cleanup-windows:
	@read -p "Delete all volume files? [y/N]: " ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		docker compose down todoapp-postgres && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Volume files have been deleted"; \
	else \
		echo "Removing volume files has been canceled"; \
	fi





migrate-create:
	@if [ -z "$(seq)" ]; then \
  	echo "Variable seq is missing" \
  	exit 1; \
  	fi
	docker compose run --rm todoapp-postgres-migrate \
		-create
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"


migrate-create-windows:
	@if [ -z "$(seq)" ]; then \
  	echo "Variable seq is missing" \
  	exit 1; \
  	fi
	docker compose run --rm todoapp-postgres-migrate create -ext sql -dir /migrations -seq "$(seq)"

migrate-up:
		make migrate-action action=up

migrate-down:
		make migrate-action action=down

migrate-action:
		@if [ -z "$(action)" ]; then \
		echo "Variable action is missing" \
		exit 1; \
		fi;

		docker compose run --rm todoapp-postgres-migrate \
    	--path /migrations \
    	--database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
    	"$(action)"

migrate-up-windows:
	make migrate-action action=up SHELL="$(SHELL)"

migrate-down-windows:
	make migrate-action action=down SHELL="$(SHELL)"

migrate-action-windows:
	@if [ -z "$(action)" ]; then \
		echo "Variable action is missing"; \
		exit 1; \
	fi
	docker compose run --rm todoapp-postgres-migrate -path /migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable" $(action)

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

app-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go