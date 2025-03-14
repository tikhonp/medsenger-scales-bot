run: dev

dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
dev:
	docker compose -f compose.yaml up

build-dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
build-dev:
	docker compose -f compose.yaml up --build

prod: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
prod:
	docker compose -f compose.prod.yaml up --build -d

fprod:
	docker compose -f compose.prod.yaml down

logs-prod:
	docker compose -f compose.prod.yaml logs -f -n 100

go-to-server-container:
	docker exec -it agents-scales-server /bin/bash

db-status:
	@goose postgres "$(shell /bin/get_db_string)" -dir=migrations status

db-up:
	@goose postgres "$(shell /bin/get_db_string)" -dir=migrations up

db-reset:
	@goose postgres "$(shell /bin/get_db_string)" -dir=migrations reset

templ:
	@templ generate
