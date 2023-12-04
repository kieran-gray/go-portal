ifeq ("$(env)", "")
	ENVIRONMENT = "local"
else
	ENVIRONMENT = $(env)
endif

install: build

build:
	docker compose build --no-cache app dev

run:
	docker compose --profile $(ENVIRONMENT) up --build

seed:
	docker compose run --no-deps app sh -c '\
	/uploadMockFiles --env-file ./app/config/.env.$(ENVIRONMENT)'
