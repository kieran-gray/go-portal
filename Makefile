ifeq ("$(env)", "")
	ENVIRONMENT = "local"
else
	ENVIRONMENT = $(env)
endif

install: build

build:
	docker compose build

run:
	docker compose --profile $(ENVIRONMENT) up --force-recreate

seed:
	docker compose run --no-deps dev /uploadMockFiles --env-file ./app/config/.env.local
