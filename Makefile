include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up -d todoapp-postgres

env-down:
	docker compose down todoapp-postgres

env-cleanup:
	@read -p "Delete all volume files? [y/N]: " ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		docker compose down todoapp-postgres && \
		rm -rf out/pgdata && \
		echo "Volume files have been deleted"; \
	else \
		echo "Removing volume files has been canceled"; \
	fi


#env-cleanup:
#	@read -p "Delete all volume files? [y/N]: " ans; \
#	if [ "$$ans" = "y" ]; then \
#	  	docker compose down todoapp-postgres && \
#	  	rm -rf out/pgdata
#	  	echo "Volume files has been deleted"; \
#	else \
#	  echo "Removing volume files has been canceled"; \
#	  fi

#docker compose run отрабатывает 1 раз , up - поднимается и живет
# -create команда чтобы создать файлы миграции
# -ext sql название расширений файла
# -dir папка где искать файлы для миграций
# -seq название для миграций migrate-create seq=some_value
migrate-create:
	@if [ -z "$(seq)" ]; then \
  	echo "Variable seq is missing" \
  	exit 1; \
  	fi
	docker compose run --rm todoapp-postgres-migrate create -ext sql -dir /migrations -seq "$(seq)"
#	docker compose run --rm todoapp-postgres-migrate \
#		-create
#		-ext sql \
#		-dir /migrations \
#		-seq "$(seq)"



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


#TODO надо поставить make и протестировать все команды, + что то осознанное вписать в переменные связанные с БД