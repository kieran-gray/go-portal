install: build

build:
	docker compose build --no-cache app dev

run-local:
	docker compose --profile local up --build

run-live:
	docker compose --profile live up --build

seed-live:
	docker compose run --no-deps app sh -c '\
	/uploadMockFiles --env-file ./app/config/.env.live'

seed-local:
	docker compose run dev sh -c '\
	/uploadMockFiles --env-file ./config/.env.local'