SOURCE_COMMIT_SHA := $(shell git rev-parse HEAD)

ENVS := SOURCE_COMMIT=${SOURCE_COMMIT_SHA} COMPOSE_BAKE=true


.PHONY: run dev build-dev prod fprod logs-prod go-to-server-container templ db-status db-up db-reset build-prod-image update-deps

run: dev

dev:
	${ENVS} docker compose -f compose.yaml up

build-dev:
	${ENVS} docker compose -f compose.yaml up --build

fdev:
	${ENVS} docker compose -f compose.yaml down

prod:
	docker compose -f compose.test-prod.yaml up --build -d

fprod:
	docker compose -f compose.test-prod.yaml down

logs-prod:
	docker compose -f compose.test-prod.yaml logs -f -n 100

go-to-server-container:
	docker exec -it --tty agents-scales-server /bin/sh

db-status:
	docker exec -it --tty agents-scales-server /bin/manage -c migrate-status

db-up:
	docker exec -it --tty agents-scales-server /bin/manage -c migrate-up

db-reset:
	docker exec -it --tty agents-scales-server /bin/manage -c migrate-reset

templ:
	docker exec -it --tty agents-scales-server templ generate

build-prod-image:
	docker buildx build --build-arg SOURCE_COMMIT="${SOURCE_COMMIT_SHA}" --target prod -t docker.telepat.online/agents-scales-image:latest .

update-deps:
	docker exec -it --tty agents-scales-server go get -u ./...
