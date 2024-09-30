docker_dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
docker_dev:
	@docker compose up --build

status:
	@goose postgres "$(shell /bin/get_db_string)" -dir=migrations status

up:
	@goose postgres "$(shell /bin/get_db_string)" -dir=migrations up

reset:
	@goose postgres "$(shell /bin/get_db_string)" -dir=migrations reset
