SOURCE_COMMIT_SHA := $(shell git rev-parse HEAD)

ENVS := SOURCE_COMMIT=${SOURCE_COMMIT_SHA} COMPOSE_BAKE=true


.PHONY: run dev build-dev prod fprod logs-prod go-to-server-container templ db-status db-up db-reset build-prod-image

run: dev

dev:
	${ENVS} docker compose -f compose.yaml up

build-dev:
	${ENVS} docker compose -f compose.yaml up --build

fdev:
	${ENVS} docker compose -f compose.yaml down

prod:
	${ENVS} docker compose -f compose.prod.yaml up --build -d

fprod:
	${ENVS} docker compose -f compose.prod.yaml down

logs-prod:
	${ENVS} docker compose -f compose.prod.yaml logs -f -n 100

go-to-server-container:
	docker exec -it --tty agents-scales-server /bin/sh

db-status:
	docker exec -it --tty agents-scales-server goose postgres "$(shell docker exec -it agents-scales-server /bin/get_db_string)" -dir=migrations status

db-up:
	docker exec -it --tty agents-scales-server goose postgres "$(shell docker exec -it agents-scales-server /bin/get_db_string)" -dir=migrations up

db-reset:
	docker exec -it --tty agents-scales-server goose postgres "$(shell docker exec -it agents-scales-server /bin/get_db_string)" -dir=migrations reset

templ:
	docker exec -it --tty agents-scales-server templ generate

build-prod-image:
	${ENVS} docker build -t docker.telepat.online/agents-scales-image:latest --target prod .
